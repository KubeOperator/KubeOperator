package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
)

type DemoService interface {
	Save(demo dto.CreateDemo) error
	List() ([]dto.Demo, error)
	Get(name string) (dto.Demo, error)
}

type demoService struct {
	repo repository.DemoRepository
}

func NewDemoService() *demoService {
	return &demoService{repo: repository.NewDemoRepository()}
}

func (d *demoService) Save(demo dto.CreateDemo) error {
	return d.repo.Save(model.Demo{Name: demo.Name})
}

func (d *demoService) Get(name string) (dto.Demo, error) {
	var demo dto.Demo
	m, err := d.repo.Get(name)
	if err != nil {
		return demo, err
	}
	demo.Demo = m
	demo.Order = 1
	return demo, nil
}

func (d *demoService) List() ([]dto.Demo, error) {
	var demos []dto.Demo
	ms, err := d.repo.List()
	if err != nil {
		return nil, err
	}
	for i := range ms {
		demos = append(demos, dto.Demo{
			Demo:  ms[i],
			Order: 0,
		})
	}
	return demos, nil
}
