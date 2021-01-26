package client

import (
	"github.com/KubeOperator/KubeOperator/pkg/logger"
)

var log = logger.Default

type DatastoreResult struct {
	Name      string `json:"name"`
	Capacity  int    `json:"capacity"`
	FreeSpace int    `json:"freeSpace"`
}
