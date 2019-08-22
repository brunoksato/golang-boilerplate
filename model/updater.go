package model

import (
	"reflect"

	"github.com/brunoksato/golang-boilerplate/core"
)

type Updater interface {
	Update(*ModelCtx, User) core.DefaultError
}

func IsUpdater(t reflect.Type) bool {
	modelType := reflect.TypeOf((*Updater)(nil)).Elem()
	return t.Implements(modelType)
}
