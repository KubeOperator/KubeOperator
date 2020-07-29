package migrate

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	uuid "github.com/satori/go.uuid"
	"time"
)

const (
	phaseName = "migrate"
)

type InitMigrateDBPhase struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

func (i *InitMigrateDBPhase) Init() error {
	var log = logger.Default
	for _, m := range model.Models {
		log.Infof("migrate table: %s", m.TableName())
		db.DB.AutoMigrate(m)
	}
	for _, d := range model.InitData {
		switch v := d.(type) {
		case model.User:
			op, ok := d.(model.User)
			if ok {
				user := model.User{}
				db.DB.Model(model.User{}).Where("name = ?", op.Name).First(&user)
				if db.DB.NewRecord(user) {
					db.DB.Create(d)
				}
			}
		case model.Credential:
			op, ok := d.(model.Credential)
			if ok {
				credential := model.Credential{}
				db.DB.Model(model.Credential{}).Where("name = ?", op.Name).First(&credential)
				if db.DB.NewRecord(credential) {
					db.DB.Create(d)
				}
			}
		case model.Project:
			op, ok := d.(model.Project)
			if ok {
				project := model.Project{}
				db.DB.Model(model.Project{}).Where("name = ?", op.Name).First(&project)
				if db.DB.NewRecord(project) {
					db.DB.Create(d)
					var clusters []model.Cluster
					db.DB.Model(model.Cluster{}).Find(&clusters)
					for _, cluster := range clusters {
						err := db.DB.Create(model.ProjectResource{
							ProjectID:    op.ID,
							ResourceId:   cluster.ID,
							ResourceType: constant.ResourceCluster,
							BaseModel: common.BaseModel{
								UpdatedAt: time.Now(),
								CreatedAt: time.Now(),
							},
							ID: uuid.NewV4().String(),
						})
						fmt.Println(err)
					}
					var hosts []model.Host
					db.DB.Model(model.Host{}).Where("cluster_id != ?", "''").Find(&hosts)
					for _, host := range hosts {
						db.DB.Create(model.ProjectResource{
							ProjectID:    op.ID,
							ResourceId:   host.ID,
							ResourceType: constant.ResourceHost,
							BaseModel: common.BaseModel{
								UpdatedAt: time.Now(),
								CreatedAt: time.Now(),
							},
							ID: uuid.NewV4().String(),
						})
					}
				}
			}

		default:
			log.Infof("insert data failed: %s", v)
		}
	}
	return nil
}

func (i *InitMigrateDBPhase) PhaseName() string {
	return phaseName
}
