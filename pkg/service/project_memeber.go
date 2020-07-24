package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
)

type ProjectMemberService interface {
	PageByProjectId(num, size int, projectId string) (page.Page, error)
	Batch(op dto.ProjectMemberOP) error
}

type projectMemberService struct {
	projectMemberRepo repository.ProjectMemberRepository
}

func NewProjectMemberService() ProjectMemberService {
	return &projectMemberService{
		projectMemberRepo: repository.NewProjectMemberRepository(),
	}
}

func (p projectMemberService) PageByProjectId(num, size int, projectId string) (page.Page, error) {
	var page page.Page
	total, mos, err := p.projectMemberRepo.PageByProjectId(num, size, projectId)
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
		opItems = append(opItems, model.ProjectMember{
			BaseModel: common.BaseModel{},
			ID:        item.ID,
			UserID:    item.UserID,
			ProjectID: item.ProjectID,
			Role:      item.Role,
		})
	}
	return p.projectMemberRepo.Batch(op.Operation, opItems)
}
