package model

import (
	"fmt"
	"reflect"

	"github.com/asaskevich/govalidator"
	"github.com/brunoksato/golang-boilerplate/core"
	"github.com/brunoksato/golang-boilerplate/util"
)

type Validator interface {
	ValidateForCreate() core.DefaultError
	ValidateForUpdate() core.DefaultError
	ValidateForDelete(*ModelCtx) core.DefaultError
	ValidateField(string) core.DefaultError
}

func IsValidator(t reflect.Type) bool {
	modelType := reflect.TypeOf((*Validator)(nil)).Elem()
	return t.Implements(modelType)
}

func ValidateStruct(item interface{}) core.DefaultError {
	_, err := govalidator.ValidateStruct(item)
	if err != nil {
		return core.NewBusinessError(err.Error())
	}
	return nil
}

func ValidateStructField(item interface{}, f string) core.DefaultError {
	_, err := govalidator.ValidateStruct(item)
	errStr := govalidator.ErrorByField(err, f)

	if errStr != "" {
		return core.NewBusinessError(fmt.Sprintf("%s: %s;", f, errStr))
	}

	return nil
}

func ValidateStructFields(item interface{}, fs []string) core.DefaultError {
	_, err := govalidator.ValidateStruct(item)
	if err == nil {
		return nil
	}

	result := ""
	for _, f := range fs {
		errStr := govalidator.ErrorByField(err, f)
		if errStr != "" {
			result = fmt.Sprintf("%s%s: %s;", result, f, errStr)
		}
	}

	if result == "" {
		return nil
	}

	return core.NewBusinessError(result)
}

func ValidateRequiredFields(item interface{}, fields []string) core.DefaultError {
	itemVal := reflect.ValueOf(item)
	errStr := ""

	for _, fName := range fields {
		f := itemVal.FieldByName(fName)
		if util.IsEmptyValue(f) {
			errStr = fmt.Sprintf("%s%s: non zero value required;", errStr, fName)
		}
	}

	if errStr != "" {
		return core.NewBusinessError(errStr)
	}

	return nil
}
