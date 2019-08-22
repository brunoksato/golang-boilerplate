package middleware

import (
	"github.com/brunoksato/golang-boilerplate/config"
	"github.com/labstack/echo/v4"
	"github.com/olivere/elastic"
)

func ElasticMiddleware(es *elastic.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if es == nil {
				es = config.InitElasticSearchAndLogger()
			}

			c.Set("Elastic", es)

			return next(c)
		}
	}
}
