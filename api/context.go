package api

import (
	"reflect"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/brunoksato/golang-boilerplate/model"
	"github.com/brunoksato/golang-boilerplate/core"
	log "github.com/brunoksato/golang-boilerplate/log"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/olivere/elastic"
)

type Context struct {
	Database      *gorm.DB
	Elastic       *elastic.Client
	Logger        *logrus.Entry
	Payload       map[string]interface{}
	Request       map[string]interface{}
	Type          reflect.Type
	ParentType    reflect.Type
	User          model.User
	AppName       string
	Method        string
	Path          string
	Endpoint      string
	RequestID     string
	Configuration model.Configuration
	APIType       core.APIType
	ModelCtx      *model.ModelCtx
}

func ServerContext(c echo.Context) *Context {
	var t reflect.Type
	var parentType reflect.Type
	var user model.User
	var APIType core.APIType
	var configuration model.Configuration

	if c.Get("User") != nil {
		user = c.Get("User").(model.User)
	}

	if c.Get("Type") != nil {
		t = c.Get("Type").(reflect.Type)
	}

	if c.Get("ParentType") != nil {
		parentType = c.Get("ParentType").(reflect.Type)
	}

	if c.Get("APIType") != nil {
		APIType = c.Get("APIType").(core.APIType)
	}

	if c.Get("Configuration") != nil {
		configuration = c.Get("Configuration").(model.Configuration)
	}

	return &Context{
		RequestID:     c.Get("RequestID").(string),
		Database:      c.Get("Database").(*gorm.DB),
		Elastic:       c.Get("Elastic").(*elastic.Client),
		User:          user,
		Configuration: configuration,
		Logger:        log.Logger(c),
		Type:          t,
		ParentType:    parentType,
		Payload:       make(map[string]interface{}),
		Request:       make(map[string]interface{}),
		APIType:       APIType,
	}
}

func ArgonContext(c echo.Context) *model.ModelCtx {
	ctx := ServerContext(c)
	if ctx.ModelCtx == nil {
		c.Logger().Debug("Creating a new ArgonContext.")
		ctx.ModelCtx = &model.ModelCtx{
			RequestID:     ctx.RequestID,
			APIType:       ctx.APIType,
			Database:      ctx.Database,
			Configuration: c.Get("Configuration").(model.Configuration),
			User:          c.Get("User").(model.User),
			Logger:        log.Logger(c),
		}
	}

	return ctx.ModelCtx
}

func ArgonContextTransaction(c echo.Context, db *gorm.DB) *model.ModelCtx {
	ctx := ServerContext(c)
	if ctx.ModelCtx == nil {
		c.Logger().Debug("Creating a new ArgonContext.")
		ctx.ModelCtx = &model.ModelCtx{
			RequestID:     ctx.RequestID,
			APIType:       ctx.APIType,
			Database:      db,
			Configuration: c.Get("Configuration").(model.Configuration),
			User:          c.Get("User").(model.User),
			Logger:        log.Logger(c),
		}
	}

	ctx.Database = db

	return ctx.ModelCtx
}

func ActiveUserID(c echo.Context, ctx *Context) uint {
	userID := ctx.User.ID
	uid, _ := strconv.Atoi(c.Param("userId"))
	if uid > 0 {
		userID = uint(uid)
	}
	return userID
}
