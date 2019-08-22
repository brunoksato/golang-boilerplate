package middleware

import (
	"net/http"

	"github.com/brunoksato/golang-boilerplate/core"
	"github.com/labstack/echo/v4"
)

func SettingHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		xCompany := c.Request().Header.Get("X-Company")
		switch xCompany {
		case "Office":
			c.Set("APIType", core.USER_API)
		case "CronJob":
			c.Set("APIType", core.CRONJOB_API)
			cCronjob := c.Request().Header.Get("X-Cronjob")
			if cCronjob != "youpassword" {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}
		case "Admin":
			c.Set("APIType", core.ADMIN_API)
		default:
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		return next(c)
	}
}
