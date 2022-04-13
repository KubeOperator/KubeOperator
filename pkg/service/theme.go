package service

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/vincent-petithory/dataurl"
)

type ThemeService interface {
	GetConsumerTheme() (*dto.Theme, error)
	SaveConsumerTheme(theme dto.Theme) (*dto.Theme, error)
}

func NewThemeService() ThemeService {
	return &themeService{}
}

type themeService struct {
}

func (*themeService) GetConsumerTheme() (*dto.Theme, error) {
	var theme model.Theme
	if err := db.DB.First(&theme).Error; err != nil {
		return nil, err
	}
	return &dto.Theme{Theme: theme}, nil
}

func (*themeService) SaveConsumerTheme(theme dto.Theme) (*dto.Theme, error) {
	var t dto.Theme

	if theme.Logo != "" {
		types := [...]string{"image/gif", "image/jpeg", "image/jpg", "image/png"}
		dataURL, err := dataurl.DecodeString(theme.Logo)
		if err != nil {
			return nil, err
		}
		result := false
		contentType := dataURL.MediaType.ContentType()
		for _, value := range types {
			if value == contentType {
				result = true
				break
			}
		}
		if !result {
			return nil, errors.New("Unsupported file type ")
		}
	}

	if err := db.DB.First(&t).Error; err != nil {
		return nil, err
	}
	t.Logo = theme.Logo
	t.SystemName = theme.SystemName
	if err := db.DB.Save(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}
