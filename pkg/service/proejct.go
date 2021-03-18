package service

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
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

	old, _ := p.Get(creation.Name)
	if old.ID != "" {
		return nil, errors.New(ProjectNameExist)
	}

	project := model.Project{
		BaseModel:   common.BaseModel{},
		Name:        creation.Name,
		Description: creation.Description,
	}

	err := p.projectRepo.Save(&project)
	if err != nil {
		return nil, err
	}

	//user, err := p.userService.Get(creation.UserName)
	//if err != nil {
	//	return nil, err
	//}
	//if !user.IsAdmin {
	//	projectMember := model.ProjectMember{
	//		ProjectID: project.ID,
	//		UserID:    user.ID,
	//		Role:      constant.ProjectRoleProjectManager,
	//	}
	//	err := p.projectMemberRepo.Create(&projectMember)
	//	if err != nil {
	//		return nil, err
	//	}
	//}
	return &dto.Project{Project: project}, err
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
