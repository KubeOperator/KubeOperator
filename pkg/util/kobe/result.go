package kobe

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type Stat struct {
	Host        string `json:"host"`
	Change      int    `json:"change"`
	Failures    int    `json:"failures"`
	Ignored     int    `json:"ignored"`
	Ok          int    `json:"ok"`
	Rescued     int    `json:"rescued"`
	Skipped     int    `json:"skipped"`
	Unreachable int    `json:"unreachable"`
}

type Duration struct {
	End   time.Time `json:"end"`
	Start time.Time `json:"start"`
}

type Task struct {
	Duration Duration                          `json:"duration"`
	Name     string                            `json:"name"`
	Hosts    map[string]map[string]interface{} `json:"hosts"`
}

type Play struct {
	Name     string   `json:"name"`
	Duration Duration `json:"duration"`
	Tasks    []Task   `json:"tasks"`
}
type HostFailedInfo map[string]string

type Result struct {
	Stats          map[string]Stat `json:"stats"`
	Plays          []Play          `json:"plays"`
	HostFailedInfo HostFailedInfo  `json:"-"`
}

func (r *Result) GatherFailedInfo() {
	hostFailed := make(map[string]string)
	for _, play := range r.Plays {
		for _, task := range play.Tasks {
			for name := range task.Hosts {
				hostResult := task.Hosts[name]
				val, ok := hostResult["failed"]
				if !ok {
					val, ok = hostResult["unreachable"]
				}
				hostFailed[name] = ""
				if ok && val.(bool) {
					b, err := json.Marshal(hostResult)
					if err != nil {
						hostFailed[name] = err.Error()
					}
					hostFailed[name] = string(b)
				}
			}
		}
	}
	r.HostFailedInfo = hostFailed
}

func ParseResult(content string) (result Result, err error) {
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return result, errors.Wrap(err, fmt.Sprintf("unmarshal contant %s failed: %v", content, err))
	}
	return
}
