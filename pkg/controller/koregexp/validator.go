package koregexp

import (
	"fmt"
	"regexp"
	"unicode"

	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/go-playground/validator/v10"
)

var log = logger.Default

func CheckNamePattern(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	result, err := regexp.MatchString("^[a-zA-Z\u4e00-\u9fa5]{1}[a-zA-Z0-9_\u4e00-\u9fa5]{0,30}$", value)
	if err != nil {
		log.Errorf("regexp matchString failed, %v", err)
	}
	return result
}

func CheckClusterNamePattern(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	result, err := regexp.MatchString("[a-z]([-a-z0-9]*[a-z0-9])?(\\\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*", value)
	if err != nil {
		log.Errorf("regexp check cluster name matchString failed, %v", err)
	}
	return result
}

func CheckCommonNamePattern(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	result, err := regexp.MatchString(`[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*`, value)
	if err != nil {
		log.Errorf("regexp check common name matchString failed, %v", err)
	}
	return result
}

func CheckHostNamePattern(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	result, err := regexp.MatchString(`^[a-zA-Z0-9.]{0,30}$`, value)
	if err != nil {
		log.Errorf("regexp check host name matchString failed, %v", err)
	}
	return result
}

func CheckIpPattern(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	result, err := regexp.MatchString(`^((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}$`, value)
	if err != nil {
		log.Errorf("regexp check ip matchString failed, %v", err)
	}
	return result
}

func CheckPasswordPattern(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if len(value) < 8 || len(value) > 30 {
		return false
	}

	hasNum := false
	hasLetter := false
	for _, r := range value {
		if unicode.IsLetter(r) && !hasLetter {
			hasLetter = true
		}
		if unicode.IsNumber(r) && !hasNum {
			hasNum = true
		}
		if hasLetter && hasNum {
			return true
		}
	}

	return false
}

func CheckVmConfigPattern(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	result, err := regexp.MatchString("^[a-zA-Z0-9]{1}[a-zA-Z0-9]{0,30}$", value)
	fmt.Println(err)
	return result
}
