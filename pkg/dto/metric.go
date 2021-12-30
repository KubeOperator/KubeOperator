package dto

const (
	MetricTypeMatrix = "matrix"
	MetricTypeVector = "vector"
)

type QueryOptions struct {
	Level         string   `json:"level"`
	NodeName      string   `json:"nodeName"`
	Start         int      `json:"start"`
	End           int      `json:"end"`
	Step          int      `json:"step"`
	MetricsFilter []string `json:"metricsFilter"`
}

type Metadata struct {
	Metric string `json:"metric,omitempty" description:"metric name"`
	Type   string `json:"type,omitempty" description:"metric type"`
	Help   string `json:"help,omitempty" description:"metric description"`
}

type Metric struct {
	MetricName string `json:"metric_name,omitempty" description:"metric name, eg. scheduler_up_sum" csv:"metric_name"`
	MetricData `json:"data,omitempty" description:"actual metric result"`
	Error      string `json:"error,omitempty" csv:"-"`
}

type MetricValues []MetricValue

type MetricData struct {
	MetricType   string `json:"resultType,omitempty" description:"result type, one of matrix, vector" csv:"metric_type"`
	MetricValues `json:"result,omitempty" description:"metric data including labels, time series and values" csv:"metric_values"`
}

type Point [2]float64
type ExportPoint [2]float64

type MetricValue struct {
	Metadata       map[string]string `json:"metric,omitempty" description:"time series labels"`
	Sample         *Point            `json:"value,omitempty" description:"time series, values of vector type"`
	Series         []Point           `json:"values,omitempty" description:"time series, values of matrix type"`
	ExportSample   *ExportPoint      `json:"exported_value,omitempty" description:"exported time series, values of vector type"`
	ExportedSeries []ExportPoint     `json:"exported_values,omitempty" description:"exported time series, values of matrix type"`

	MinValue     string `json:"min_value" description:"minimum value from monitor points"`
	MaxValue     string `json:"max_value" description:"maximum value from monitor points"`
	AvgValue     string `json:"avg_value" description:"average value from monitor points"`
	SumValue     string `json:"sum_value" description:"sum value from monitor points"`
	Fee          string `json:"fee" description:"resource fee"`
	ResourceUnit string `json:"resource_unit"`
	CurrencyUnit string `json:"currency_unit"`
}
