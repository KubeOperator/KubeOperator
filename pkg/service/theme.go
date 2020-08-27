package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type ThemeService interface{
	GetConsumerTheme() (*dto.Theme, error)
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
