package typeparse

import (
	"strings"
)

type QueryCondition struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

func ParseConditionToSql(condition QueryCondition) string {
	switch condition.Operator {
	case "like":
		return condition.Field + " LIKE '%" + condition.Value + "%'"
	case "not like":
		return condition.Field + " NOT LIKE '%" + condition.Value + "%'"
	case "eq":
		return condition.Field + " = '" + condition.Value + "'"
	case "ne":
		return condition.Field + " != '" + condition.Value + "'"
	default:
		return condition.Field + " LIKE '%" + condition.Value + "%'"
	}
}

func ParseConditionQuickToSql(condition QueryCondition, conditionNames ...string) string {
	var returnSQL string
	for _, con := range conditionNames {
		returnSQL += con + " LIKE '%" + condition.Value + "%'" + " OR "
	}
	if strings.Contains(returnSQL, " OR ") {
		returnSQL = returnSQL[0 : len(returnSQL)-4]
	}
	return returnSQL
}
