package model

type Credential struct {
	Id         string `gorm:"primary_key size:64"`
	Name       string `gorm:"size:128"`
	User       string `gorm:"size:128"`
	Password   string `gorm:"size:256"`
	PrivateKey string
}
