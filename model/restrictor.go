package model

import (
	"reflect"

	"github.com/brunoksato/golang-boilerplate/core"
)

type Restrictor interface {
	UserCanView(ctx *ModelCtx, u User) (bool, core.DefaultError)
	UserCanCreate(ctx *ModelCtx, u User) (bool, core.DefaultError)
	UserCanUpdate(ctx *ModelCtx, u User, fields []string) (bool, core.DefaultError)
	UserCanDelete(ctx *ModelCtx, u User) (bool, core.DefaultError)
}

func IsRestrictor(t reflect.Type) bool {
	modelType := reflect.TypeOf((*Restrictor)(nil)).Elem()
	return t.Implements(modelType)
}
