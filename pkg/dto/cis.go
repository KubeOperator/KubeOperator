package dto

import (
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type CisTask struct {
	model.CisTask
}

type CisBatch struct {
	Items     []CisTask
	Operation string
}

type CisTaskDetail struct {
	ClusterName    string `json:"clusterName"`
	ClusterVersion string `json:"clusterVersion"`
	model.CisTaskWithResult
	NodeList CisNodeList `json:"nodeList"`
}


type CisNodeList []CisNode
type CisTestList []CisTest
type CisResultList []CisResult

type CisNode struct {
	Id        string      `json:"id"`
	Version   string      `json:"version"`
	Text      string      `json:"text"`
	NodeType  string      `json:"node_type"`
	Tests     CisTestList `json:"tests"`
	TotalPass int         `json:"total_pass"`
	TotalFail int         `json:"total_fail"`
	TotalWarn int         `json:"total_warn"`
	TotalInfo int         `json:"total_info"`
}

type CisTest struct {
	Section string        `json:"section"`
	Pass    int           `json:"pass"`
	Fail    int           `json:"fail"`
	Warn    int           `json:"warn"`
	Info    int           `json:"info"`
	Desc    string        `json:"desc"`
	Results CisResultList `json:"results"`
}

type CisResult struct {
	TestNumber  string `json:"test_number"`
	TestDesc    string `json:"test_desc"`
	Remediation string `json:"remediation"`
	Status      string `json:"status"`
	Scored      bool   `json:"scored"`
	Reason      string `json:"reason"`
}
