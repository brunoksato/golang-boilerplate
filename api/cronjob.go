package api

import (
	"net/http"

	"github.com/brunoksato/golang-boilerplate/core"
	log "github.com/brunoksato/golang-boilerplate/log"
	"github.com/labstack/echo/v4"
)

func CronJobSample(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	tx := db.Begin()
	if err := tx.Error; err != nil {
		tx.Rollback()
		return log.AddDefaultError(c, core.NewServerError(err.Error()))
	}

	tx.Commit()

	return c.JSON(http.StatusOK, map[string]interface{}{"status": "ok"})
}
