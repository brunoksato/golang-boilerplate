package server

import (
	"time"

	config "github.com/brunoksato/golang-boilerplate/config"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/olivere/elastic"
)

var RW_DB_POOL *gorm.DB
var ES *elastic.Client

func Start() *echo.Echo {
	config.Init()
	RW_DB_POOL = config.InitDB()
	RW_DB_POOL.DB().SetConnMaxLifetime(time.Second * 30)
	RW_DB_POOL.DB().SetMaxIdleConns(40)
	RW_DB_POOL.DB().SetMaxOpenConns(40)
	RW_DB_POOL.LogMode(true)
	ES = config.InitElasticSearchAndLogger()
	root := SetupRouter(ProductionMiddlewareConfigurer{})

	return root
}
