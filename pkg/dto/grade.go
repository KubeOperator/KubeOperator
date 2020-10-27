package dto

type ClusterGrade struct {
	Score    int                `json:"score"`
	TotalSum Summary            `json:"totalSum"`
	ListSum  map[string]Summary `json:"listSum"`
	Results  []NamespaceResult  `json:"results"`
}

type Summary struct {
	Danger  int `json:"danger"`
	Warning int `json:"warning"`
	Success int `json:"success"`
}

type NamespaceResult struct {
	Namespace string                  `json:"namespace"`
	Results   []NamespaceResultDetail `json:"results"`
}

type NamespaceResultDetail struct {
	Name       string      `json:"name"`
	Kind       string      `json:"kind"`
	PodResults []PodResult `json:"podResults"`
}

type PodResult struct {
	ID       string `json:"id"`
	Category string `json:"category"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
	Success  bool   `json:"success"`
}
