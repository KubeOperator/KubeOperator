package service

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	dbUtil "github.com/KubeOperator/KubeOperator/pkg/util/db"
	"github.com/jinzhu/gorm"
)

var (
	ProjectNameExist = "NAME_EXISTS"
)

type ProjectService interface {
	Get(name string) (*dto.Project, error)
	List(user dto.SessionUser, conditions condition.Conditions) ([]dto.Project, error)
	Page(num, size int, user dto.SessionUser, conditions condition.Conditions) (*page.Page, error)
	Delete(name string) error
	Create(creation dto.ProjectCreate) (*dto.Project, error)
	Batch(op dto.ProjectOp) error
	Update(name string, update dto.ProjectUpdate) (*dto.Project, error)
	GetResourceTree(user dto.SessionUser) ([]dto.ProjectResourceTree, error)
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

func (p *projectService) List(user dto.SessionUser, conditions condition.Conditions) ([]dto.Project, error) {

	var (
		pa          page.Page
		projectDTOS []dto.Project
		projects    []model.Project
	)

	d := db.DB.Model(model.Project{})
	if err := dbUtil.WithConditions(&d, model.Project{}, conditions); err != nil {
		return nil, err
	}

	if user.IsAdmin {
		if err := d.Count(&pa.Total).Order("created_at desc").Find(&projects).Error; err != nil {
			return nil, err
		}
	} else {

		var projectMembers []model.ProjectMember
		var projectIds []string
		var projectResources []model.ProjectResource

		err := db.DB.Where("user_id = ?", user.UserId).Find(&projectMembers).Error
		if err != nil {
			return nil, err
		}
		for _, pm := range projectMembers {
			projectIds = append(projectIds, pm.ProjectID)
		}
		err = db.DB.Raw("SELECT DISTINCT project_id  FROM ko_project_resource WHERE resource_type = 'CLUSTER' AND resource_id in (SELECT DISTINCT cluster_id FROM ko_cluster_member WHERE user_id = ?)", user.UserId).Scan(&projectResources).Error
		if err != nil {
			return nil, err
		}
		for _, pm := range projectResources {
			projectIds = append(projectIds, pm.ProjectID)
		}
		err = d.Count(&pa.Total).Order("created_at desc").Where("id in (?)", projectIds).Find(&projects).Error
		if err != nil {
			return nil, err
		}
	}

	for _, mo := range projects {
		projectDTOS = append(projectDTOS, dto.Project{Project: mo})
	}

	return projectDTOS, nil
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
	mo.Description = update.Description
	err := p.projectRepo.Save(&mo)
	if err := db.DB.Save(&mo).Error; err != nil {
		return nil, err
	}

	return &dto.Project{Project: mo}, err
}

func (p *projectService) Page(num, size int, user dto.SessionUser, conditions condition.Conditions) (*page.Page, error) {

	var (
		pa          page.Page
		projectDTOS []dto.Project
		projects    []model.Project
	)

	d := db.DB.Model(model.Project{})
	if err := dbUtil.WithConditions(&d, model.Project{}, conditions); err != nil {
		return nil, err
	}

	if user.IsAdmin {
		if err := d.Count(&pa.Total).Order("created_at desc").Offset((num - 1) * size).Limit(size).Find(&projects).Error; err != nil {
			return nil, err
		}
	} else {
		if user.IsRole(constant.RoleProjectManager) {

			var projectResources []model.ProjectMember
			err := db.DB.Where("user_id = ?", user.UserId).Find(&projectResources).Error
			if err != nil {
				return nil, err
			}
			var projectIds []string
			for _, pm := range projectResources {
				projectIds = append(projectIds, pm.ProjectID)
			}
			err = d.Count(&pa.Total).Order("created_at desc").Where("id in (?)", projectIds).Offset((num - 1) * size).Limit(size).Find(&projects).Error
			if err != nil {
				return nil, err
			}
		} else {
			var projectResources []model.ProjectResource
			err := db.DB.Raw("SELECT DISTINCT project_id  FROM ko_project_resource WHERE resource_type = 'CLUSTER' AND resource_id in (SELECT DISTINCT cluster_id FROM ko_cluster_member WHERE user_id = ?)", user.UserId).Scan(&projectResources).Error
			if err != nil {
				return nil, err
			}
			var projectIds []string
			for _, pm := range projectResources {
				projectIds = append(projectIds, pm.ProjectID)
			}
			err = d.Count(&pa.Total).Order("created_at desc").Where("id in (?)", projectIds).Offset((num - 1) * size).Limit(size).Find(&projects).Error
			if err != nil {
				return nil, err
			}
		}
	}

	for _, mo := range projects {
		projectDTOS = append(projectDTOS, dto.Project{Project: mo})
	}
	pa.Items = projectDTOS
	return &pa, nil
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

func (p projectService) GetResourceTree(user dto.SessionUser) ([]dto.ProjectResourceTree, error) {
	var (
		projects []model.Project
		tree     []dto.ProjectResourceTree
	)

	if user.IsAdmin {
		if err := db.DB.Model(model.Project{}).Order("created_at ASC").Find(&projects).Error; err != nil {
			return nil, err
		}
	} else if user.IsRole(constant.RoleProjectManager) {
		if err := db.DB.Model(model.Project{}).Where("name = ?", user.CurrentProject).Order("created_at ASC").Find(&projects).Error; err != nil {
			return nil, err
		}
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
		var clusters []model.Cluster
		err := db.DB.Raw("SELECT * FROM ko_cluster WHERE id IN (SELECT resource_id FROM ko_project_resource WHERE project_id"+
			" = ( select distinct ID FROM ko_project WHERE name  = ?) And resource_type = 'CLUSTER')", t.Label).Scan(&clusters).Error
		if err != nil {
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
