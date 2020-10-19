package dto

import "github.com/fairwindsops/polaris/pkg/validator"

type ClusterGrade struct {
	Score    int                                      `json:"score"`
	TotalSum Summary                                  `json:"totalSum"`
	ListSum  map[string]Summary                       `json:"listSum"`
	Results  map[string][]*validator.ControllerResult `json:"results"`
}

type Summary struct {
	Danger  int `json:"danger"`
	Warning int `json:"warning"`
	Success int `json:"success"`
}
