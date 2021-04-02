package service

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/errorf"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/jinzhu/gorm"
)

type ProjectMemberService interface {
	Page(projectName string, num, size int) (*page.Page, error)
	Batch(op dto.ProjectMemberOP) error
	GetUsers(name string) (dto.AddMemberResponse, error)
	Create(projectName string, request dto.ProjectMemberCreate) ([]dto.ProjectMember, error)
	Get(name string, projectName string) (*dto.ProjectMember, error)
	Delete(name, projectName string) error
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

func (p *projectMemberService) Page(projectName string, num, size int) (*page.Page, error) {

	var (
		pa                page.Page
		projectMembers    []model.ProjectMember
		project           model.Project
		projectMemberDTOs []dto.ProjectMember
	)

	if err := db.DB.Where("name = ?", projectName).First(&project).Error; err != nil {
		return nil, err
	}

	err := db.DB.Model(&model.ProjectMember{}).
		Where("project_id = ?", project.ID).
		Preload("User").
		Count(&pa.Total).
		Offset((num - 1) * size).
		Limit(size).
		Find(&projectMembers).Error

	for _, mo := range projectMembers {
		projectMemberDTOs = append(projectMemberDTOs, toProjectMemberDTO(mo))
	}
	pa.Items = projectMemberDTOs
	return &pa, err
}

func (p *projectMemberService) Batch(op dto.ProjectMemberOP) error {
	var opItems []model.ProjectMember
	//for _, item := range op.Items {
	//	id := ""
	//	user, err := NewUserService().Get(item.Username)
	//	if err != nil {
	//		return err
	//	}
	//	project, err := NewProjectService().Get(item.ProjectName)
	//	if err != nil {
	//		return err
	//	}
	//	if op.Operation == constant.BatchOperationUpdate || op.Operation == constant.BatchOperationDelete {
	//		var pm model.ProjectMember
	//		err := db.DB.Where("user_id = ? AND project_id = ?", user.ID, project.ID).First(&pm).Error
	//		if err != nil {
	//			return err
	//		}
	//		id = pm.ID
	//	}
	//
	//	opItems = append(opItems, model.ProjectMember{
	//		BaseModel: common.BaseModel{},
	//		ID:        id,
	//		UserID:    user.ID,
	//		ProjectID: project.ID,
	//		Role:      item.Role,
	//	})
	//}
	return p.projectMemberRepo.Batch(op.Operation, opItems)
}

func (p *projectMemberService) GetUsers(name string) (dto.AddMemberResponse, error) {
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

func (p *projectMemberService) Create(projectName string, request dto.ProjectMemberCreate) ([]dto.ProjectMember, error) {

	var (
		project model.Project
		errs    errorf.CErrFs
		result  []dto.ProjectMember
	)
	if err := db.DB.Model(model.Project{}).Where("name = ?", projectName).First(&project).Error; err != nil {
		return nil, err
	}

	for _, name := range request.Usernames {
		var user model.User
		if err := db.DB.Model(model.User{}).Where("name = ?", name).First(&user).Error; err != nil {
			errs = errs.Add(errorf.New("USER_IS_NOT_FOUND", name))
			continue
		} else {
			var oldPm dto.ProjectMember
			if err := db.DB.Where("user_id = ? AND project_id = ?", user.ID, project.ID).Find(&oldPm).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
				errs = errs.Add(errorf.New(err.Error()))
				continue
			}
			if oldPm.ID != "" {
				errs = errs.Add(errorf.New("USER_IS_ADD", name))
				continue
			}
			pm := model.ProjectMember{
				UserID:    user.ID,
				Role:      constant.ProjectRoleProjectManager,
				ProjectID: project.ID,
			}
			if err := db.DB.Create(&pm).Error; err != nil {
				errs = errs.Add(errorf.New(err.Error()))
			}
			d := toProjectMemberDTO(pm)
			result = append(result, d)
		}
	}
	if len(errs) > 0 {
		return result, errs
	} else {
		return result, nil
	}
}

func (p *projectMemberService) Get(name string, projectName string) (*dto.ProjectMember, error) {
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

func (p *projectMemberService) Delete(name, projectName string) error {

	var (
		project model.Project
		pm      model.ProjectMember
	)
	user, err := p.userService.Get(name)
	if err != nil {
		return err
	}
	if err := db.DB.Model(model.Project{}).Where("name = ?", projectName).First(&project).Error; err != nil {
		return err
	}
	if err := db.DB.Model(model.ProjectMember{}).Where("project_id = ? AND user_id = ?", project.ID, user.ID).Find(&pm).Error; err != nil {
		return err
	}
	if err := db.DB.Delete(&pm).Error; err != nil {
		return err
	}
	return nil
}

func toProjectMemberDTO(mo model.ProjectMember) dto.ProjectMember {
	d := dto.ProjectMember{
		ProjectMember: mo,
		Username:      mo.User.Name,
		UserStatus: func() string {
			if mo.User.IsActive {
				return constant.UserStatusActive
			}
			return constant.UserStatusPassive
		}(),
		Email: mo.User.Email,
	}
	return d
}
