package host

type Host struct {
	Id   string `gorm:"primary_key size:64"`
	Name string `gorm:"size:128"`
	Ip   string `gorm:"size:128"`
	Port int
}
