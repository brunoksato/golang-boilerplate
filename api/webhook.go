package api

import (
	"net/http"

	"github.com/brunoksato/golang-boilerplate/core"
	log "github.com/brunoksato/golang-boilerplate/log"
	"github.com/labstack/echo/v4"
)

func WebhookSample(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database
	tx := db.Begin()

	ctx.Logger.Info("WebhookSample Start")

	if err := tx.Error; err != nil {
		tx.Rollback()
		return log.AddDefaultError(c, core.NewServerError(err.Error()))
	}

	event := map[string]interface{}{}
	if err := c.Bind(&event); err != nil {
		return log.AddDefaultError(c, core.NewServerError(err.Error()))
	}

	tx.Commit()

	ctx.Logger.Info("WebhookSample End")

	return c.JSON(http.StatusOK, map[string]interface{}{"status": "ok"})
}
