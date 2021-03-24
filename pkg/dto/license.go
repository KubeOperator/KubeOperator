package dto

type License struct {
	Message     string      `json:"message"`
	Status      string      `json:"status"`
	LicenseInfo LicenseInfo `json:"license"`
}

type LicenseInfo struct {
	Corporation string `json:"corporation"`
	Product     string `json:"product"`
	Edition     string `json:"edition"`
	Count       int    `json:"count"`
	Expired     string `json:"expired"`
}
