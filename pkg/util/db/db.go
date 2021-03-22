package db

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	dbConfig "github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/jinzhu/gorm"
	"strings"
)

func WithConditions(p interface{}, conditions condition.Conditions) (*gorm.DB, error) {

	db := dbConfig.DB.Model(p)
	if !conditions.IsZero() {
		val, ok := conditions["quick"]
		if ok {
			for _, f := range db.NewScope(p).GetStructFields() {
				if !strings.Contains(f.DBName, "id") && f.IsNormal {
					db = db.Or(fmt.Sprintf("%s LIKE ?", f.DBName), "%"+fmt.Sprintf("%v", val.Value)+"%")
				}
			}
			return db, nil
		}
		for _, v := range conditions {
			switch strings.ToLower(v.Operator) {
			case "like":
				db = db.Where(fmt.Sprintf("%s LIKE ?", v.Field), "%"+fmt.Sprintf("%v", v.Value)+"%")
			case "not like":
				db = db.Where(fmt.Sprintf("%s NOT LIKE ?", v.Field), "%"+fmt.Sprintf("%v", v.Value)+"%")
			case "eq":
				db = db.Where(fmt.Sprintf("%s = ?", v.Field), v.Value)
			case "ne":
				db = db.Where(fmt.Sprintf("%s != ?", v.Field), v.Value)
			case "gt":
				db = db.Where(fmt.Sprintf("%s > ?", v.Field), v.Value)
			case "ge":
				db = db.Where(fmt.Sprintf("%s >= ?", v.Field), v.Value)
			case "lt":
				db = db.Where(fmt.Sprintf("%s < ?", v.Field), v.Value)
			case "le":
				db = db.Where(fmt.Sprintf("%s <= ?", v.Field), v.Value)
			case "in":
				val, ok := v.Value.([]interface{})
				if !ok {
					return nil, fmt.Errorf("condition %s must be a list", v.Field)
				}

				db = db.Where(fmt.Sprintf("%s IN (?)", v.Field), val)
			case "not in":
				val, ok := v.Value.([]interface{})
				if !ok {
					return nil, fmt.Errorf("condition %s must be a list", v.Field)
				}

				db = db.Where(fmt.Sprintf("%s NOT IN (?)", v.Field), val)
			case "between":
				val, ok := v.Value.([]interface{})
				if !ok {
					return nil, fmt.Errorf("condition %s must be a list", v.Field)
				}
				if !(len(val) == 2) {
					return nil, fmt.Errorf("condition %s length must be 2", v.Field)
				}
				db = db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", v.Field), val)
			}
		}
	}
	return db, nil
}
