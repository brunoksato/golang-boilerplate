package model

import (
	"reflect"
	"testing"

	"github.com/brunoksato/golang-boilerplate/core"
)

func TestIsUpdater(t *testing.T) {
	core.AssertTrue(t, IsUpdater(reflect.TypeOf(&User{})))
	core.AssertFalse(t, IsUpdater(reflect.TypeOf(&Configuration{})))
}

