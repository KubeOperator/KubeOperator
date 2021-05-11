package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/errorf"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/jinzhu/gorm"
)

type ClusterMemberService interface {
	Page(clusterName string, num, size int) (*page.Page, error)
	GetUsers(name string) (dto.UsersResponse, error)
	Create(clusterName string, request dto.ClusterMemberCreate) ([]dto.ClusterMember, error)
	Delete(name, clusterName string) error
}

type clusterMemberService struct {
	userService UserService
}

func NewClusterMemberService() ClusterMemberService {
	return &clusterMemberService{
		userService: NewUserService(),
	}
}

func (c *clusterMemberService) Page(clusterName string, num, size int) (*page.Page, error) {
	var (
		pa                page.Page
		clusterMembers    []model.ClusterMember
		cluster           model.Cluster
		clusterMemberDTOs []dto.ClusterMember
	)

	if err := db.DB.Where("name = ?", clusterName).First(&cluster).Error; err != nil {
		return nil, err
	}
	err := db.DB.Model(&model.ClusterMember{}).Where("cluster_id = ?", cluster.ID).
		Preload("User").
		Count(&pa.Total).
		Order("created_at desc").
		Offset((num - 1) * size).
		Limit(size).
		Find(&clusterMembers).Error

	for _, mo := range clusterMembers {
		clusterMemberDTOs = append(clusterMemberDTOs, toClusterMemberDTO(mo, cluster.Name))
	}
	pa.Items = clusterMemberDTOs
	return &pa, err
}

func (c *clusterMemberService) GetUsers(name string) (dto.UsersResponse, error) {
	var (
		result dto.UsersResponse
		users  []model.User
	)
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

func (c *clusterMemberService) Create(clusterName string, request dto.ClusterMemberCreate) ([]dto.ClusterMember, error) {
	var (
		cluster model.Cluster
		errs    errorf.CErrFs
		result  []dto.ClusterMember
	)
	if err := db.DB.Model(model.Cluster{}).Where("name = ?", clusterName).First(&cluster).Error; err != nil {
		return nil, err
	}

	for _, name := range request.Usernames {
		var user model.User
		if err := db.DB.Model(model.User{}).Where("name = ?", name).First(&user).Error; err != nil {
			errs = errs.Add(errorf.New("USER_IS_NOT_FOUND", name))
			continue
		} else {
			var oldCm dto.ClusterMember
			if err := db.DB.Where("user_id = ? AND cluster_id = ?", user.ID, cluster.ID).Find(&oldCm).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
				errs = errs.Add(errorf.New(err.Error()))
				continue
			}
			if oldCm.ID != "" {
				errs = errs.Add(errorf.New("USER_IS_ADD", name))
				continue
			}
			cm := model.ClusterMember{
				UserID:    user.ID,
				Role:      constant.ProjectRoleClusterManager,
				ClusterID: cluster.ID,
			}
			if err := db.DB.Create(&cm).Error; err != nil {
				errs = errs.Add(errorf.New(err.Error()))
			}
			d := toClusterMemberDTO(cm, cluster.Name)
			result = append(result, d)
		}
	}
	if len(errs) > 0 {
		return result, errs
	} else {
		return result, nil
	}
}

func (c *clusterMemberService) Delete(name, clusterName string) error {
	var (
		cluster model.Cluster
		cm      model.ClusterMember
		pr      model.ProjectResource
	)
	user, err := c.userService.Get(name)
	if err != nil {
		return err
	}
	if err := db.DB.Model(model.Cluster{}).Where("name = ?", clusterName).First(&cluster).Error; err != nil {
		return err
	}
	if err := db.DB.Model(model.ClusterMember{}).Where("cluster_id = ? AND user_id = ?", cluster.ID, user.ID).Find(&cm).Error; err != nil {
		return err
	}
	if err := db.DB.Delete(&cm).Error; err != nil {
		return err
	}
	if err := db.DB.Debug().Model(model.ProjectResource{}).Where("resource_id = ? AND resource_type = 'CLUSTER'", cluster.ID).Find(&pr).Error; err != nil {
		return err
	}
	if user.CurrentProjectID == pr.ProjectID {
		user.User.CurrentProjectID = ""
		db.DB.Save(&user.User)
	}
	return nil
}

func toClusterMemberDTO(mo model.ClusterMember, clusterName string) dto.ClusterMember {
	d := dto.ClusterMember{
		ClusterMember: mo,
		Username:      mo.User.Name,
		ClusterName:   clusterName,
		Email:         mo.User.Email,
	}
	return d
}
