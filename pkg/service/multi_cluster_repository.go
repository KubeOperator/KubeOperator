package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/git"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"path"
)

type MultiClusterRepositoryService interface {
	Page(num, size int) (*page.Page, error)
	List() ([]dto.MultiClusterRepository, error)
	Get(name string) (*dto.MultiClusterRepository, error)
	Create(request dto.MultiClusterRepositoryCreateRequest) (*dto.MultiClusterRepository, error)
	Delete(name string) error
	Update(name string, request dto.MultiClusterRepositoryUpdateRequest) (*dto.MultiClusterRepository, error)
	UpdateClusterRelations(req dto.UpdateRelationRequest) error
}
type multiClusterRepositoryService struct {
}

func (m *multiClusterRepositoryService) Update(name string, request dto.MultiClusterRepositoryUpdateRequest) (*dto.MultiClusterRepository, error) {
	var item model.MultiClusterRepository
	if err := db.DB.Where(model.MultiClusterRepository{Name: name}).First(&item).Error; err != nil {
		return nil, err
	}
	item.SyncEnable = request.SyncEnable
	item.GitTimeout = request.GitTimeout
	item.SyncInterval = request.SyncInterval
	if err := db.DB.Save(&item).Error; err != nil {
		return nil, err
	}

	return &dto.MultiClusterRepository{MultiClusterRepository: item}, nil
}

func (m *multiClusterRepositoryService) Delete(name string) error {
	var item model.MultiClusterRepository
	if err := db.DB.Where(model.MultiClusterRepository{Name: name}).First(&item).Error; err != nil {
		return err
	}
	if err := db.DB.Delete(&item).Error; err != nil {
		return err
	}
	return nil
}

func NewMultiClusterRepositoryService() MultiClusterRepositoryService {
	return &multiClusterRepositoryService{
	}
}
func (m *multiClusterRepositoryService) UpdateClusterRelations(req dto.UpdateRelationRequest) error {
	panic("implement me")
}

func (m *multiClusterRepositoryService) Page(num, size int) (*page.Page, error) {
	var p page.Page
	var mos []model.MultiClusterRepository
	if err := db.DB.Model(model.MultiClusterRepository{}).
		Count(&p.Total).
		Offset((num - 1) * size).
		Limit(size).
		Find(&mos).Error; err != nil {
		return nil, err
	}
	var dtos []dto.MultiClusterRepository
	for _, mo := range mos {
		d := dto.MultiClusterRepository{MultiClusterRepository: mo}
		dtos = append(dtos, d)
	}
	p.Items = dtos
	return &p, nil
}

func (m *multiClusterRepositoryService) List() ([]dto.MultiClusterRepository, error) {
	var mos []model.MultiClusterRepository
	if err := db.DB.Model(model.MultiClusterRepository{}).
		Find(&mos).Error; err != nil {
		return nil, err
	}
	var dtos []dto.MultiClusterRepository
	for _, mo := range mos {
		d := dto.MultiClusterRepository{MultiClusterRepository: mo}
		dtos = append(dtos, d)
	}
	return dtos, nil
}

func (m *multiClusterRepositoryService) Get(name string) (*dto.MultiClusterRepository, error) {
	var item model.MultiClusterRepository
	if err := db.DB.Where(model.MultiClusterRepository{Name: name}).First(&item).Error; err != nil {
		return nil, err
	}
	return &dto.MultiClusterRepository{MultiClusterRepository: item}, nil

}

func (m *multiClusterRepositoryService) Create(req dto.MultiClusterRepositoryCreateRequest) (*dto.MultiClusterRepository, error) {
	mo := model.MultiClusterRepository{
		Name:     req.Name,
		Source:   req.Source,
		Username: req.Username,
		Password: req.Password,
		Branch:   req.Branch,
	}
	if err := db.DB.Create(&mo).Error; err != nil {
		return nil, err
	}
	go func() {
		if err := git.CloneRepository(mo.Source, path.Join(constant.DefaultRepositoryDir, mo.Name), mo.Branch, &http.BasicAuth{Username: mo.Username, Password: mo.Password}); err != nil {
			mo.Status = constant.ClusterFailed
			mo.Message = err.Error()
			db.DB.Save(&mo)
		}
		mo.Status = constant.ClusterRunning
		db.DB.Save(&mo)
	}()

	return &dto.MultiClusterRepository{MultiClusterRepository: mo}, nil
}
