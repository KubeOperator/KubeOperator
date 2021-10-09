package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/jinzhu/gorm"

	originalDB "github.com/KubeOperator/KubeOperator/pkg/db"
)

func WithProjectResource(db **gorm.DB, projectName string, resourceType string) ([]model.ProjectResource, error) {
	var (
		p      model.Project
		res    []model.ProjectResource
		resIds []string
	)
	if err := originalDB.DB.
		Where("name = ?", projectName).
		First(&p).Error; err != nil {
		return res, err
	}

	if err := originalDB.DB.Where("resource_type = ? AND project_id = ?", resourceType, p.ID).Find(&res).Error; err != nil {
		return res, err
	}
	for i := range res {
		resIds = append(resIds, res[i].ResourceID)
	}
	*db = (*db).Where("id IN (?)", resIds)
	return res, nil
}

// 后续如果有非created_at字段需要加时间选择，需要在conditions中加上time字段
func WithConditions(db **gorm.DB, model interface{}, conditions condition.Conditions) error {
	if !conditions.IsZero() {
		val, ok := conditions["quick"]
		var (
			keys   []string
			values []interface{}
		)
		if ok {
			for _, f := range (*db).NewScope(model).GetStructFields() {
				if !strings.Contains(strings.ToLower(f.Name), "id") && f.IsNormal {
					keys = append(keys, fmt.Sprintf("%s LIKE ?", dealReservedWord(f.DBName)))
					values = append(values, "%"+fmt.Sprintf("%v", val.Value)+"%")
				}
			}
			var sql string
			for i := range keys {
				if i != 0 {
					sql += " OR "
				}
				sql += keys[i]
			}
			*db = (*db).Where(sql, values...)
			return nil
		}
		for _, v := range conditions {
			if v.Field == "created_at" {
				switch strings.ToLower(v.Operator) {
				case "eq":
					tm1, err := time.Parse("2006-01-02", fmt.Sprintf("%s", v.Value))
					if err != nil {
						return fmt.Errorf("wrong time format")
					}
					tm2 := tm1.AddDate(0, 0, 1)
					*db = (*db).Where(fmt.Sprintf("%s > ? AND %s < ?", v.Field, v.Field), tm1, tm2)
				case "gt":
					tm1, err := time.Parse("2006-01-02", fmt.Sprintf("%s", v.Value))
					tm2 := tm1.AddDate(0, 0, 1)
					if err != nil {
						return fmt.Errorf("wrong time format")
					}
					*db = (*db).Where(fmt.Sprintf("%s > ?", v.Field), tm2)
				case "ge":
					tm1, err := time.Parse("2006-01-02", fmt.Sprintf("%s", v.Value))
					if err != nil {
						return fmt.Errorf("wrong time format")
					}
					*db = (*db).Where(fmt.Sprintf("%s >= ?", v.Field), tm1)
				case "lt":
					tm1, err := time.Parse("2006-01-02", fmt.Sprintf("%s", v.Value))
					if err != nil {
						return fmt.Errorf("wrong time format")
					}
					*db = (*db).Where(fmt.Sprintf("%s < ?", v.Field), tm1)
				case "le":
					tm1, err := time.Parse("2006-01-02", fmt.Sprintf("%s", v.Value))
					tm2 := tm1.AddDate(0, 0, 1)
					if err != nil {
						return fmt.Errorf("wrong time format")
					}
					*db = (*db).Where(fmt.Sprintf("%s <= ?", v.Field), tm2)
				case "between":
					val, ok := v.Value.([]interface{})
					if !ok {
						return fmt.Errorf("condition %s must be a list", v.Field)
					}
					if !(len(val) == 2) {
						return fmt.Errorf("condition %s length must be 2", v.Field)
					}
					tm1, err := time.Parse("2006-01-02", fmt.Sprintf("%s", val[0]))
					if err != nil {
						return fmt.Errorf("wrong time format")
					}
					tm2, err := time.Parse("2006-01-02", fmt.Sprintf("%s", val[1]))
					if err != nil {
						return fmt.Errorf("wrong time format")
					}
					tm3 := tm2.AddDate(0, 0, 1)
					*db = (*db).Where(fmt.Sprintf("%s BETWEEN ? AND ?", v.Field), tm1, tm3)
				}
			} else {
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
					*db = (*db).Where(fmt.Sprintf("%s BETWEEN ? AND ?", v.Field), val[0], val[1])
				}
			}
		}
	}
	return nil
}

func dealReservedWord(name string) string {
	reservedWord := []string{"memory"}
	for _, word := range reservedWord {
		if name == word {
			return "`" + name + "`"
		}
	}
	return name
}
