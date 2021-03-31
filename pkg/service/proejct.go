package service

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/jinzhu/gorm"
)

var (
	ProjectNameExist = "NAME_EXISTS"
)

type ProjectService interface {
	Get(name string) (*dto.Project, error)
	List() ([]dto.Project, error)
	Page(num, size int, userId string) (page.Page, error)
	Delete(name string) error
	Create(creation dto.ProjectCreate) (*dto.Project, error)
	Batch(op dto.ProjectOp) error
	Update(name string, update dto.ProjectUpdate) (*dto.Project, error)
	GetResourceTree() ([]dto.ProjectResourceTree, error)
}

type projectService struct {
	projectRepo       repository.ProjectRepository
	userService       UserService
	projectMemberRepo repository.ProjectMemberRepository
}

func NewProjectService() ProjectService {
	return &projectService{
		projectRepo:       repository.NewProjectRepository(),
		userService:       NewUserService(),
		projectMemberRepo: repository.NewProjectMemberRepository(),
	}
}

func (p *projectService) Get(name string) (*dto.Project, error) {
	var projectDTO dto.Project
	mo, err := p.projectRepo.Get(name)
	if err != nil {
		return nil, err
	}
	projectDTO.Project = mo
	return &projectDTO, err
}

func (p *projectService) List() ([]dto.Project, error) {
	var projectDTOs []dto.Project
	mos, err := p.projectRepo.List()
	if err != nil {
		return projectDTOs, err
	}
	for _, mo := range mos {
		projectDTOs = append(projectDTOs, dto.Project{Project: mo})
	}
	return projectDTOs, err
}

func (p *projectService) Create(creation dto.ProjectCreate) (*dto.Project, error) {

	var old model.Project
	if err := db.DB.Where("name = ?", creation.Name).First(&old).Error; !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	if old.ID != "" {
		return nil, errors.New(ProjectNameExist)
	}

	project := model.Project{
		BaseModel:   common.BaseModel{},
		Name:        creation.Name,
		Description: creation.Description,
	}
	if err := db.DB.Create(&project).Error; err != nil {
		return nil, err
	}
	return &dto.Project{Project: project}, nil
}

func (p *projectService) Update(name string, update dto.ProjectUpdate) (*dto.Project, error) {

	var mo model.Project
	if err := db.DB.Where(model.Project{Name: name}).First(&mo).Error; err != nil {
		return nil, err
	}
	if update.Description != "" {
		mo.Description = update.Description
	}
	err := p.projectRepo.Save(&mo)
	if err := db.DB.Save(&mo).Error; err != nil {
		return nil, err
	}

	return &dto.Project{Project: mo}, err
}

func (p *projectService) Page(num, size int, userId string) (page.Page, error) {
	var page page.Page
	var projectDTOS []dto.Project
	total, mos, err := p.projectRepo.Page(num, size, userId)
	if err != nil {
		return page, err
	}
	for _, mo := range mos {
		projectDTOS = append(projectDTOS, dto.Project{Project: mo})
	}
	page.Total = total
	page.Items = projectDTOS
	return page, err
}

func (p *projectService) Delete(name string) error {
	err := p.projectRepo.Delete(name)
	if err != nil {
		return err
	}
	return nil
}

func (p *projectService) Batch(op dto.ProjectOp) error {
	var deleteItems []model.Project
	for _, item := range op.Items {
		deleteItems = append(deleteItems, model.Project{
			BaseModel: common.BaseModel{},
			ID:        item.ID,
			Name:      item.Name,
		})
	}
	err := p.projectRepo.Batch(op.Operation, deleteItems)
	if err != nil {
		return err
	}
	return nil
}

func (p projectService) GetResourceTree() ([]dto.ProjectResourceTree, error) {
	var (
		projects []model.Project
		tree     []dto.ProjectResourceTree
	)

	if err := db.DB.Model(model.Project{}).Order("name").Find(&projects).Error; err != nil {
		return nil, err
	}
	id := 0
	for _, p := range projects {
		id++
		tree = append(tree, dto.ProjectResourceTree{
			ID:    id,
			Label: p.Name,
			Type:  constant.ResourceProject,
		})
	}
	for i, t := range tree {
		var project model.Project
		if err := db.DB.Where("name = ?", t.Label).First(&project).Error; err != nil {
			return nil, err
		}
		var projectResources []model.ProjectResource
		if err := db.DB.Where("project_id = ? AND resource_type = ?", project.ID, constant.ResourceCluster).Find(&projectResources).Error; err != nil {
			return nil, err
		}
		var resourceIds []string
		for _, pr := range projectResources {
			resourceIds = append(resourceIds, pr.ResourceID)
		}
		var clusters []model.Cluster
		if err := db.DB.Model(&model.Cluster{}).
			Where("id in (?)", resourceIds).
			Find(&clusters).Error; err != nil {
			return nil, err
		}
		for _, c := range clusters {
			id++
			tree[i].Children = append(tree[i].Children, dto.ProjectResourceTree{
				ID:    id,
				Label: c.Name,
				Type:  constant.ResourceCluster,
			})
		}
	}
	return tree, nil
}
