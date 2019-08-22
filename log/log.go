package log

import (
	"fmt"
	"net/http"

	"github.com/brunoksato/golang-boilerplate/core"
	"github.com/labstack/echo/v4"
)

func AddPayloadWarning(c echo.Context, code int, message string) error {
	// payload["warning"] = map[string]interface{}{"code": code, "message": message}
	return nil
}

func AddPayloadError(c echo.Context, code int, message string) error {
	return c.JSON(http.StatusBadRequest, map[string]interface{}{"code": code, "message": message})
}

func AddPermissionError(c echo.Context, code int, message string) error {
	return c.JSON(http.StatusForbidden, map[string]interface{}{"code": code, "message": message})
}

func AddNotFoundError(c echo.Context, code int, message string) error {
	return c.JSON(http.StatusNotFound, map[string]interface{}{"code": code, "message": message})
}

func AddServerError(c echo.Context, code int, message string) error {
	return c.JSON(http.StatusInternalServerError, map[string]interface{}{"code": code, "message": message})
}

func AddDefaultError(c echo.Context, errModel core.DefaultError) error {
	var err error
	msg := fmt.Sprintf("%s (caller: %s)", errModel.Error(), errModel.Location())

	code := errModel.Code()
	if errModel.Subcode() != 0 {
		code = errModel.Subcode()
	}

	params := errModel.Data()
	params["code"] = errModel.Code()
	params["subcode"] = errModel.Subcode()

	logger := LoggerForParams(c, params)

	switch errModel.Code() {
	case 300:
		logger.Warning("Warning: " + msg)
		err = AddPayloadWarning(c, code, errModel.Error())
	case 400:
		logger.Info("Business Error: " + msg)
		err = AddPayloadError(c, code, errModel.Error())
	case 403:
		logger.Warning("Permission Error: " + msg)
		err = AddPermissionError(c, code, errModel.Error())
	case 404:
		logger.Info("Not Found: " + msg)
		err = AddNotFoundError(c, code, errModel.Error())
	default:
		logger.Error("Server Error: " + msg)
		err = AddServerError(c, code, errModel.Error())
	}

	return err
}
