package prometheus

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/prometheus/client_golang/api"
	apiv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

const MeteringDefaultTimeout = 20 * time.Second

type PrometheusService struct {
	client apiv1.API
}

func NewPrometheusService(Endpoint string) (PrometheusService, error) {
	cfg := api.Config{
		Address: Endpoint,
	}

	client, err := api.NewClient(cfg)
	return PrometheusService{client: apiv1.NewAPI(client)}, err
}

func (p PrometheusService) GetNamedMetrics(metrics []string, ts time.Time, opts *dto.QueryOptions) []dto.Metric {
	var res []dto.Metric
	var mtx sync.Mutex
	var wg sync.WaitGroup

	for _, metric := range metrics {
		wg.Add(1)
		go func(metric string) {
			parsedResp := dto.Metric{MetricName: metric}

			value, _, err := p.client.Query(context.Background(), makeExpr(metric, opts), ts)
			if err != nil {
				parsedResp.Error = err.Error()
			} else {
				parsedResp.MetricData = parseQueryResp(value)
			}

			mtx.Lock()
			res = append(res, parsedResp)
			mtx.Unlock()

			wg.Done()
		}(metric)
	}

	wg.Wait()

	return res
}

func (p PrometheusService) GetNamedMetricsOverTime(opts *dto.QueryOptions) []dto.Metric {
	var res []dto.Metric
	var mtx sync.Mutex
	var wg sync.WaitGroup

	timeRange := apiv1.Range{
		Start: time.Unix(int64(opts.Start), 0),
		End:   time.Unix(int64(opts.End), 0),
		Step:  time.Duration(opts.Step * int(time.Minute)),
	}

	for _, metric := range opts.MetricsFilter {
		wg.Add(1)
		go func(metric string) {
			parsedResp := dto.Metric{MetricName: metric}
			value, _, err := p.client.QueryRange(context.Background(), makeExpr(metric, opts), timeRange)
			if err != nil {
				parsedResp.Error = err.Error()
			} else {
				parsedResp.MetricData = parseQueryRangeResp(value)
			}

			mtx.Lock()
			res = append(res, parsedResp)
			mtx.Unlock()

			wg.Done()
		}(metric)
	}

	wg.Wait()

	return res
}

func (p PrometheusService) GetMetadata(namespace string) []dto.Metadata {
	var meta []dto.Metadata
	var matchTarget string

	if namespace != "" {
		// Filter metrics available to members of this namespace
		matchTarget = fmt.Sprintf("{namespace=\"%s\"}", namespace)
	}
	items, err := p.client.TargetsMetadata(context.Background(), matchTarget, "", "")
	if err != nil {
		logger.Log.Error(err)
		return meta
	}

	// Deduplication
	set := make(map[string]bool)
	for _, item := range items {
		_, ok := set[item.Metric]
		if !ok {
			set[item.Metric] = true
			meta = append(meta, dto.Metadata{
				Metric: item.Metric,
				Type:   string(item.Type),
				Help:   item.Help,
			})
		}
	}

	return meta
}

func (p PrometheusService) GetMetricLabelSet(expr string, start, end time.Time) []map[string]string {
	var res []map[string]string

	labelSet, _, err := p.client.Series(context.Background(), []string{expr}, start, end)
	if err != nil {
		logger.Log.Error(err)
		return []map[string]string{}
	}

	for _, item := range labelSet {
		var tmp = map[string]string{}
		for key, val := range item {
			if key == "__name__" {
				continue
			}
			tmp[string(key)] = string(val)
		}

		res = append(res, tmp)
	}

	return res
}

func parseQueryRangeResp(value model.Value) dto.MetricData {
	res := dto.MetricData{MetricType: dto.MetricTypeMatrix}

	data, _ := value.(model.Matrix)

	for _, v := range data {
		mv := dto.MetricValue{
			Metadata: make(map[string]string),
		}

		for k, v := range v.Metric {
			mv.Metadata[string(k)] = string(v)
		}

		for _, k := range v.Values {
			mv.Series = append(mv.Series, dto.Point{float64(k.Timestamp) / 1000, float64(k.Value)})
		}

		res.MetricValues = append(res.MetricValues, mv)
	}

	return res
}

func parseQueryResp(value model.Value) dto.MetricData {
	res := dto.MetricData{MetricType: dto.MetricTypeVector}

	data, _ := value.(model.Vector)

	for _, v := range data {
		mv := dto.MetricValue{
			Metadata: make(map[string]string),
		}

		for k, v := range v.Metric {
			mv.Metadata[string(k)] = string(v)
		}

		mv.Sample = &dto.Point{float64(v.Timestamp) / 1000, float64(v.Value)}

		res.MetricValues = append(res.MetricValues, mv)
	}

	return res
}
