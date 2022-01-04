package validator_error

import (
	"errors"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entrans "github.com/go-playground/validator/v10/translations/en"
	zhtrans "github.com/go-playground/validator/v10/translations/zh"
	"github.com/kataras/iris/v12/context"
	"reflect"
	"strings"
)

func Tr(ctx context.Context, validate *validator.Validate, err error) error {
	en := en.New()
	zh := zh.New()
	uni := ut.New(en, zh)
	language := ctx.GetLocale().Language()

	var lang string
	if language == "zh-CN" {
		lang = "zh"
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get(lang)
			return name
		})
		trans, _ := uni.GetTranslator(lang)
		_ = zhtrans.RegisterDefaultTranslations(validate, trans)
		errs := err.(validator.ValidationErrors)
		return errors.New(removeStructName(errs.Translate(trans)))
	} else {
		lang = "en"
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get(lang)
			return name
		})
		trans, _ := uni.GetTranslator(lang)
		_ = entrans.RegisterDefaultTranslations(validate, trans)
		errs := err.(validator.ValidationErrors)
		return errors.New(removeStructName(errs.Translate(trans)))
	}
}

func removeStructName(fields map[string]string) string {
	var errMsg string
	result := map[string]string{}
	for field, err := range fields {
		result[field[strings.Index(field, ".")+1:]] = err
	}
	for _, set := range result {
		errMsg = errMsg + set + " "
	}
	return errMsg
}

func RegisterTagNameFunc(ctx context.Context, validate *validator.Validate) {
	language := ctx.GetLocale().Language()
	var lang string
	if language == "zh-CN" {
		lang = "zh"
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get(lang)
			return name
		})
	} else {
		lang = "en"
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get(lang)
			return name
		})
	}
}
