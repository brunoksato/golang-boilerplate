package model

import (
	"fmt"
	"reflect"

	"github.com/brunoksato/golang-boilerplate/core"
)

type Deleter interface {
	Delete(*ModelCtx, User) core.DefaultError
	ValidateForDelete(*ModelCtx) core.DefaultError
}

func IsDeleter(t reflect.Type) bool {
	modelType := reflect.TypeOf((*Deleter)(nil)).Elem()
	return t.Implements(modelType)
}

func DefaultDelete(ctx *ModelCtx, item interface{}) core.DefaultError {
	itemID, err := core.GetID(item)
	if err != nil || itemID == 0 {
		return core.NewNotFoundError(fmt.Sprintf("A %v must exist in the database to be deleted", reflect.TypeOf(item)))
	}

	db := ctx.Database
	err = db.Set("gorm:save_associations", false).Delete(item).Error
	if err != nil {
		return core.NewServerError("Database error while deleting: " + err.Error())
	}
	return nil
}
