package model

import (
	"reflect"
	"testing"

	"github.com/brunoksato/golang-boilerplate/core"
)

func TestIsDeleter(t *testing.T) {
	core.AssertFalse(t, IsDeleter(reflect.TypeOf(&User{})))
	core.AssertFalse(t, IsDeleter(reflect.TypeOf(&Configuration{})))
}
