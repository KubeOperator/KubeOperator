package service

import (
	"errors"
	"fmt"

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
	UserIsAdd = "USER_IS_ADD"
)

type ProjectMemberService interface {
	PageByProjectName(num, size int, projectName string) (page.Page, error)
	Batch(op dto.ProjectMemberOP) error
	GetUsers(name string) (dto.AddMemberResponse, error)
	Create(request dto.ProjectMemberCreate) (*dto.ProjectMember, error)
	Get(name string, projectName string) (*dto.ProjectMember, error)
	CheckUserProjectPermissionByName(userId string, projectNames []string) (bool, error)
	CheckUserProjectPermissionByID(userId string, projectIDs []string) (bool, error)
}

type projectMemberService struct {
	projectMemberRepo repository.ProjectMemberRepository
	userService       UserService
	projectRepo       repository.ProjectRepository
}

func NewProjectMemberService() ProjectMemberService {
	return &projectMemberService{
		projectMemberRepo: repository.NewProjectMemberRepository(),
		userService:       NewUserService(),
		projectRepo:       repository.NewProjectRepository(),
	}
}

func (p projectMemberService) PageByProjectName(num, size int, projectName string) (page.Page, error) {
	var page page.Page
	project, err := p.projectRepo.Get(projectName)
	if err != nil {
		return page, err
	}
	total, mos, err := p.projectMemberRepo.PageByProjectId(num, size, project.ID)
	if err != nil {
		return page, err
	}

	var result []dto.ProjectMember
	for _, mo := range mos {
		result = append(result, dto.ProjectMember{ProjectMember: mo, UserName: mo.User.Name})
	}
	page.Items = result
	page.Total = total
	return page, err
}

func (p projectMemberService) Batch(op dto.ProjectMemberOP) error {
	var opItems []model.ProjectMember
	for _, item := range op.Items {
		id := ""
		user, err := NewUserService().Get(item.Username)
		if err != nil {
			return err
		}
		project, err := p.projectRepo.Get(item.ProjectName)
		if err != nil {
			return err
		}
		if op.Operation == constant.BatchOperationUpdate || op.Operation == constant.BatchOperationDelete {
			var pm model.ProjectMember
			err := db.DB.Where("user_id = ? AND project_id = ?", user.ID, project.ID).First(&pm).Error
			if err != nil {
				return err
			}
			id = pm.ID
		}

		opItems = append(opItems, model.ProjectMember{
			BaseModel: common.BaseModel{},
			ID:        id,
			UserID:    user.ID,
			ProjectID: project.ID,
			Role:      item.Role,
		})
	}
	return p.projectMemberRepo.Batch(op.Operation, opItems)
}

func (p projectMemberService) GetUsers(name string) (dto.AddMemberResponse, error) {
	var result dto.AddMemberResponse
	var users []model.User
	err := db.DB.Select("name").Where("is_admin = 0 AND name LIKE ?", "%"+name+"%").Find(&users).Error
	if err != nil {
		return result, err
	}
	var addUsers []string
	for _, user := range users {
		addUsers = append(addUsers, user.Name)
	}
	result.Items = addUsers
	return result, nil
}

func (p projectMemberService) Create(request dto.ProjectMemberCreate) (*dto.ProjectMember, error) {
	var projectMember dto.ProjectMember
	user, err := p.userService.Get(request.Username)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, UserNotFound
		} else {
			return nil, err
		}
	}
	project, err := p.projectRepo.Get(request.ProjectName)
	if err != nil {
		return nil, err
	}
	var oldPm dto.ProjectMember
	err = db.DB.Where("user_id = ? AND project_id = ?", user.ID, project.ID).Find(&oldPm).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	if oldPm.ID != "" {
		return nil, errors.New(UserIsAdd)
	}
	pm := model.ProjectMember{
		UserID:    user.ID,
		Role:      request.Role,
		ProjectID: project.ID,
	}
	err = p.projectMemberRepo.Create(&pm)
	if err != nil {
		return nil, err
	}
	projectMember.ProjectMember = pm
	return &projectMember, nil
}

func (p projectMemberService) Get(name string, projectName string) (*dto.ProjectMember, error) {
	u, err := p.userService.Get(name)
	if err != nil {
		return nil, err
	}
	pj, err := p.projectRepo.Get(projectName)
	if err != nil {
		return nil, err
	}

	var pm model.ProjectMember
	notFound := db.DB.Where("user_id = ? AND project_id = ?", u.ID, pj.ID).First(&pm).RecordNotFound()
	if notFound {
		return nil, fmt.Errorf("project member: %s not found in project %s", name, projectName)
	}
	return &dto.ProjectMember{
		ProjectMember: pm,
	}, nil
}

func (p projectMemberService) CheckUserProjectPermissionByName(userId string, projectNames []string) (bool, error) {
	if userId == "" {
		return true, nil
	}

	var projectMembers []model.ProjectMember
	if err := db.DB.Where("user_id = ?", userId).Find(&projectMembers).Error; err != nil {
		return false, err
	}
	var projects []model.Project
	if err := db.DB.Where("name in (?)", projectNames).Find(&projects).Error; err != nil {
		return false, err
	}
	for _, project := range projects {
		hasPower := false
		for _, member := range projectMembers {
			if project.ID == member.ProjectID {
				hasPower = true
				break
			}
		}
		if !hasPower {
			return false, nil
		}
	}
	return true, nil
}

func (p projectMemberService) CheckUserProjectPermissionByID(userId string, projectIDs []string) (bool, error) {
	if userId == "" {
		return true, nil
	}

	var projectMembers []model.ProjectMember
	if err := db.DB.Where("user_id = ?", userId).Find(&projectMembers).Error; err != nil {
		return false, err
	}
	for _, id := range projectIDs {
		hasPower := false
		for _, member := range projectMembers {
			if id == member.ProjectID {
				hasPower = true
				break
			}
		}
		if !hasPower {
			return false, nil
		}
	}
	return true, nil
}
