package controller

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12/context"
)

type MultiClusterRepositoryController struct {
	Ctx                           context.Context
	MultiClusterRepositoryService service.MultiClusterRepositoryService
	MultiClusterSyncLogService    service.MultiClusterSyncLogService
}

func NewMultiClusterRepositoryController() *MultiClusterRepositoryController {
	return &MultiClusterRepositoryController{
		MultiClusterRepositoryService: service.NewMultiClusterRepositoryService(),
		MultiClusterSyncLogService:    service.NewMultiClusterSyncLogService(),
	}
}

func (m *MultiClusterRepositoryController) Get() (*page.Page, error) {
	pg, _ := m.Ctx.Values().GetBool("page")
	if pg {
		num, _ := m.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := m.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return m.MultiClusterRepositoryService.Page(num, size)
	} else {
		var pg page.Page
		items, err := m.MultiClusterRepositoryService.List()
		if err != nil {
			return nil, err
		}
		pg.Items = items
		pg.Total = len(items)
		return &pg, nil
	}
}

func (m *MultiClusterRepositoryController) GetBy(name string) (*dto.MultiClusterRepository, error) {
	return m.MultiClusterRepositoryService.Get(name)
}
func (m *MultiClusterRepositoryController) Post() (*dto.MultiClusterRepository, error) {
	var req dto.MultiClusterRepositoryCreateRequest
	if err := m.Ctx.ReadJSON(&req); err != nil {
		return nil, err
	}
	return m.MultiClusterRepositoryService.Create(req)
}

func (m *MultiClusterRepositoryController) PatchBy(name string) (*dto.MultiClusterRepository, error) {
	var req dto.MultiClusterRepositoryUpdateRequest
	if err := m.Ctx.ReadJSON(&req); err != nil {
		return nil, err
	}
	return m.MultiClusterRepositoryService.Update(name, req)
}

func (m *MultiClusterRepositoryController) DeleteBy(name string) error {
	return m.MultiClusterRepositoryService.Delete(name)
}

func (m *MultiClusterRepositoryController) GetRelationsBy(name string) ([]dto.ClusterRelation, error) {
	return m.MultiClusterRepositoryService.GetClusterRelations(name)
}

func (m *MultiClusterRepositoryController) PostRelationsBy(name string) error {
	var req dto.UpdateRelationRequest
	if err := m.Ctx.ReadJSON(&req); err != nil {
		return nil
	}
	return m.MultiClusterRepositoryService.UpdateClusterRelations(name, req)
}

func (m *MultiClusterRepositoryController) GetLogsBy(name string) (*page.Page, error) {
	pg, _ := m.Ctx.Values().GetBool("page")
	if pg {
		num, _ := m.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := m.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return m.MultiClusterSyncLogService.Page(name, num, size)
	} else {
		return nil, fmt.Errorf("this resource must has page params")
	}
}

func (m *MultiClusterRepositoryController) GetLogsDetailBy(name string, logId string) (*dto.MultiClusterSyncLogDetail, error) {
	return m.MultiClusterSyncLogService.Detail(name, logId)
}

func (m *MultiClusterRepositoryController) PostBatch() error {
	var req dto.MultiClusterRepositoryBatch
	if err := m.Ctx.ReadJSON(&req); err != nil {
		return err
	}
	return m.MultiClusterRepositoryService.Batch(req)
}
