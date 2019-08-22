package middleware

import (
	"reflect"
	"strings"

	"github.com/brunoksato/golang-boilerplate/model"
	"github.com/labstack/echo/v4"
)

func DetermineType(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var tType reflect.Type
		var parentType reflect.Type

		parts := PathParts(c.Path())
		var pathType string
		for i := 0; i < len(parts); i++ {
			pathType = parts[i]
			t := StringToType(pathType)
			if t != nil {
				if tType != nil {
					parentType = tType
				}
				tType = t
			}
		}

		c.Set("ParentType", parentType)
		c.Set("Type", tType)

		return next(c)
	}
}

func PathParts(path string) []string {
	return strings.Split(strings.Trim(path, " /"), "/")
}

func StringToType(typeName string) (t reflect.Type) {
	switch typeName {
	case "users":
		var m model.User
		t = reflect.TypeOf(m)
	case "configurations":
		var m model.Configuration
		t = reflect.TypeOf(m)
	}
	return
}
