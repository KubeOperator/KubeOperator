package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type Demo struct {
	model.Demo
	Order int
}

type CreateDemo struct {
	Name  string
	Order int
}
