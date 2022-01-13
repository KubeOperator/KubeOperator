package service

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/util/license"
)

type LicenseService interface {
	Save(content string) (*dto.License, error)
	Get() (*dto.License, error)
	GetHw() (*dto.License, error)
}

type licenseService struct {
	licenseRepo repository.LicenseRepository
}

func NewLicenseService() LicenseService {
	return &licenseService{
		licenseRepo: repository.NewLicenseRepository(),
	}
}

var (
	formatLicenseError = errors.New("parse license error")
	verificationError  = errors.New("license is invalid")
	licenseNotFound    = errors.New("license not found")
)

func (l *licenseService) Save(content string) (*dto.License, error) {
	resp, err := license.Parse(content)
	if err != nil {
		return nil, formatLicenseError
	}
	if resp.Status != "valid" {
		return nil, verificationError
	}
	err = l.licenseRepo.Save(content)
	if err != nil {
		return nil, err
	}
	return &resp.License, nil
}

func (l *licenseService) Get() (*dto.License, error) {
	var ls dto.License
	lc, err := l.licenseRepo.Get()
	if err != nil {
		return &ls, licenseNotFound
	}
	if lc.ID == "" {
		return &ls, nil
	}
	resp, err := license.Parse(lc.Content)
	if err != nil {
		return nil, formatLicenseError
	}
	ls = resp.License
	ls.Status = resp.Status
	ls.Message = resp.Message
	return &ls, err
}

func (l *licenseService) GetHw() (*dto.License, error) {
	var ls dto.License
	lc, err := l.licenseRepo.GetHw()
	if err != nil {
		return &ls, licenseNotFound
	}
	if lc.ID == "" {
		return &ls, nil
	}
	resp, err := license.Parse(lc.Content)
	if err != nil {
		return &ls, formatLicenseError
	}
	ls = resp.License
	ls.Status = resp.Status
	ls.Message = resp.Message
	return &ls, err
}
