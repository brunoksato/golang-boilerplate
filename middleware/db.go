package middleware

import (
	"github.com/brunoksato/golang-boilerplate/config"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

func DBMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if db == nil {
				db = config.InitDB()
			}

			c.Set("Database", db)

			return next(c)
		}
	}
}
