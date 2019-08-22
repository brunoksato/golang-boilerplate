package model

import (
	"reflect"
	"testing"

	"github.com/brunoksato/golang-boilerplate/core"
)

func TestIsRestrictor(t *testing.T) {
	core.AssertTrue(t, IsRestrictor(reflect.TypeOf(User{})))
	core.AssertTrue(t, IsRestrictor(reflect.TypeOf(Configuration{})))
}
