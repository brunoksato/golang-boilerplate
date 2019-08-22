package middleware

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

func InitializePayload(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		newv4 := uuid.NewV4()
		c.Set("Payload", make(map[string]interface{}))
		c.Set("Request", make(map[string]interface{}))
		c.Set("AppName", fmt.Sprintf("%s-%s", os.Getenv("SERVER_NAME"), os.Getenv("SERVER_ENV")))
		c.Set("RequestID", newv4.String())
		c.Set("Method", c.Request().Method)
		c.Set("Endpoint", fmt.Sprintf("%s %s", c.Request().Method, c.Request().URL.Path))
		c.Set("Path", c.Request().URL.String())

		//TODO???
		c.Response().After(func() {
		})

		return next(c)
	}
}
