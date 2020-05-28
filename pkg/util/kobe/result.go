package kobe

import (
	"encoding/json"
	"time"
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

type Result struct {
	Stats map[string]Stat `json:"stats"`
	Plays []Play          `json:"plays"`
}

func ParseResult(content string) (result Result, err error) {
	err = json.Unmarshal([]byte(content), &result)
	return
}
