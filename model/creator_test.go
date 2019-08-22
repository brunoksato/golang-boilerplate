package model

import (
	"reflect"
	"testing"

	"github.com/brunoksato/golang-boilerplate/core"
)

func TestIsCreator(t *testing.T) {
	core.AssertTrue(t, IsCreator(reflect.TypeOf(&User{})))
	core.AssertFalse(t, IsCreator(reflect.TypeOf(&Configuration{})))
}