package db

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/jinzhu/gorm"
	"strings"
)

func WithProjectResource(db **gorm.DB, projectName string, resourceType string) error {
	if projectName != "" {
		var (
			p   model.Project
			res []model.ProjectResource
		)
		if err := (*db).
			Where(model.Project{Name: projectName}).
			First(&p).Error; err != nil {
			return err
		}

		if err := (*db).Where(model.ProjectResource{
			ResourceType: resourceType,
			ProjectID:    p.ID,
		}).Find(&res).Error; err != nil {
			return err
		}
		if len(res) > 0 {
			resIds := func() []string {
				var r []string
				for i := range res {
					r = append(r, res[i].ResourceID)
				}
				return r
			}()
			*db = (*db).Where("id IN (?)", resIds)
		}
	}
	return nil
}

func WithConditions(db **gorm.DB, model interface{}, conditions condition.Conditions) error {
	if !conditions.IsZero() {
		val, ok := conditions["quick"]
		if ok {
			for _, f := range (*db).NewScope(model).GetStructFields() {
				if !strings.Contains(strings.ToLower(f.Name), "id") && f.IsNormal {
					*db = (*db).Or(fmt.Sprintf("%s LIKE ?", f.DBName), "%"+fmt.Sprintf("%v", val.Value)+"%")
				}
			}
			return nil
		}
		for _, v := range conditions {
			switch strings.ToLower(v.Operator) {
			case "like":
				*db = (*db).Where(fmt.Sprintf("%s LIKE ?", v.Field), "%"+fmt.Sprintf("%v", v.Value)+"%")
			case "not like":
				*db = (*db).Where(fmt.Sprintf("%s NOT LIKE ?", v.Field), "%"+fmt.Sprintf("%v", v.Value)+"%")
			case "eq":
				*db = (*db).Where(fmt.Sprintf("%s = ?", v.Field), v.Value)
			case "ne":
				*db = (*db).Where(fmt.Sprintf("%s != ?", v.Field), v.Value)
			case "gt":
				*db = (*db).Where(fmt.Sprintf("%s > ?", v.Field), v.Value)
			case "ge":
				*db = (*db).Where(fmt.Sprintf("%s >= ?", v.Field), v.Value)
			case "lt":
				*db = (*db).Where(fmt.Sprintf("%s < ?", v.Field), v.Value)
			case "le":
				*db = (*db).Where(fmt.Sprintf("%s <= ?", v.Field), v.Value)
			case "in":
				val, ok := v.Value.([]interface{})
				if !ok {
					return fmt.Errorf("condition %s must be a list", v.Field)
				}

				*db = (*db).Where(fmt.Sprintf("%s IN (?)", v.Field), val)
			case "not in":
				val, ok := v.Value.([]interface{})
				if !ok {
					return fmt.Errorf("condition %s must be a list", v.Field)
				}

				*db = (*db).Where(fmt.Sprintf("%s NOT IN (?)", v.Field), val)
			case "between":
				val, ok := v.Value.([]interface{})
				if !ok {
					return fmt.Errorf("condition %s must be a list", v.Field)
				}
				if !(len(val) == 2) {
					return fmt.Errorf("condition %s length must be 2", v.Field)
				}
				*db = (*db).Where(fmt.Sprintf("%s BETWEEN ? AND ?", v.Field), val)
			}
		}
	}
	return nil
}
