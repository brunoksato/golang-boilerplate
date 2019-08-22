package log

import (
	"github.com/Sirupsen/logrus"
	"github.com/brunoksato/golang-boilerplate/model"
	"github.com/labstack/echo/v4"
)

func Logger(c echo.Context) *logrus.Entry {
	return LoggerForParams(c, nil)
}

func LoggerForParams(c echo.Context, params map[string]interface{}) *logrus.Entry {
	fields := make(map[string]interface{})
	if params != nil {
		for k, v := range params {
			fields[k] = v
		}
	}

	var user model.User
	if c.Get("User") != nil {
		user = c.Get("User").(model.User)
		fields["id-u"] = user.ID
	} else {
		fields["id-u"] = 0
	}

	fields["@application"] = c.Get("AppName").(string)
	fields["id-req"] = c.Get("RequestID").(string)
	fields["method"] = c.Get("Method").(string)
	fields["endpoint"] = c.Get("Endpoint").(string)
	fields["path"] = c.Get("Path").(string)

	if fields["system"] == nil {
		fields["system"] = "api"
	}

	return logrus.WithFields(logrus.Fields(fields))
}
