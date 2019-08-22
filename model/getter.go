package model

import (
	"reflect"

	"github.com/brunoksato/golang-boilerplate/core"
)

type Getter interface {
	GetByID(*ModelCtx, User, uint) (Getter, core.DefaultError)
}

func IsGetter(t reflect.Type) bool {
	modelType := reflect.TypeOf((*Getter)(nil)).Elem()
	return t.Implements(modelType)
}
