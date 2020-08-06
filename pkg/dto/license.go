package dto

type License struct {
	Corporation string `json:"corporation"`
	Product     string `json:"product"`
	Edition     string `json:"edition"`
	Count       int    `json:"count"`
	Expired     string `json:"expired"`
	Message     string `json:"message"`
	Status      string `json:"status"`
}
