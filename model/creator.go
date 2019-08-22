package model

import (
	"reflect"

	"github.com/brunoksato/golang-boilerplate/core"
)

type Creator interface {
	Create(*ModelCtx, User) core.DefaultError
}

func IsCreator(t reflect.Type) bool {
	modelType := reflect.TypeOf((*Creator)(nil)).Elem()
	return t.Implements(modelType)
}
