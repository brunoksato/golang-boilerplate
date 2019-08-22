package model

import (
	"reflect"
	"testing"

	"github.com/brunoksato/golang-boilerplate/core"
)

func TestIsGetter(t *testing.T) {
	core.AssertFalse(t, IsGetter(reflect.TypeOf(&User{})))
	core.AssertFalse(t, IsGetter(reflect.TypeOf(&Configuration{})))
}
