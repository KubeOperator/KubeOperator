package repo

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

var (
	ArmRepoPort       int
	ArmRepositoryPort int

	AmdRepoPort       int
	AmdRepositoryPort int
)

// 集群架构 amd -> x86, arm -> arm, mixed -> arm
// 部署计划 *86
func LoadRegistery() {
	var registerys []model.SystemRegistry
	for i := 0; i < 3; i++ {
		if err := db.DB.Find(&registerys).Error; err != nil {
			logger.Log.Errorf("[retry %d]: load registery of failed, err: %v", i, err)
			continue
		}
		for _, re := range registerys {
			if re.Architecture == constant.ArchitectureOfAMD64 {
				AmdRepoPort = re.RepoPort
				AmdRepositoryPort = re.RegistryPort
			}
			if re.Architecture == constant.ArchitectureOfARM64 {
				ArmRepoPort = re.RepoPort
				ArmRepositoryPort = re.RegistryPort
			}
		}
		break
	}
}
