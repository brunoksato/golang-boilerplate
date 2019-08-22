package model

import (
	"reflect"

	"github.com/jinzhu/gorm"
)

type Sorter interface {
	OrderBy(*gorm.DB) *gorm.DB
}

func IsSorter(t reflect.Type) bool {
	modelType := reflect.TypeOf((*Sorter)(nil)).Elem()
	return t.Implements(modelType)
}
