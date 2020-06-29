package dto

type StorageClass struct {
	Name        string
	Provisioner string
	Vars        map[string]interface{}
}
