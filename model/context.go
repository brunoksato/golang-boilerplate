package model

import (
	"github.com/Sirupsen/logrus"
	"github.com/brunoksato/golang-boilerplate/core"
	"github.com/jinzhu/gorm"
)

type ModelCtx struct {
	RequestID     string
	APIType       core.APIType
	Database      *gorm.DB
	User          User
	Configuration Configuration
	Logger        *logrus.Entry
}
