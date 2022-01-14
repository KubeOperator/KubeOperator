package dto

type License struct {
	Corporation string `json:"corporation" validate:"required"`
	Product     string `json:"product" validate:"required"`
	Edition     string `json:"edition" validate:"required"`
	Count       int    `json:"count" validate:"required"`
	Expired     string `json:"expired" validate:"required"`
	Message     string `json:"message" validate:"-"`
	Status      string `json:"status" validate:"-"`
}
