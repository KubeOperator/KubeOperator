package credential

import (
	uuid "github.com/satori/go.uuid"
	"ko3-gin/internal/db"
	"ko3-gin/pkg/model"
)

var Service = &serviceManger{}

type serviceManger struct{}

func (s *serviceManger) Create(credential model.Credential) error {
	credential.Id = uuid.NewV4().String()
	if err := db.DB.Create(credential).Error; err != nil {
		return err
	}
	return nil
}
