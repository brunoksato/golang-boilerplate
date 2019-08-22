package middleware

import (
	"github.com/brunoksato/golang-boilerplate/model"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

func LoadConfigurations(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		db := c.Get("Database").(*gorm.DB)
		config := model.Configuration{}
		err := db.First(&config).Error
		if err == nil {
			c.Set("Configuration", config)
		}
		return next(c)
	}
}
