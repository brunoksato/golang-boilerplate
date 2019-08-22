package model

import (
	"reflect"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/brunoksato/golang-boilerplate/core"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func init() {
	db := InitTestDB()
	INITDB = db
}

func setupDB() {
	TESTDB = INITDB.Begin()
	CTX = &ModelCtx{}
	CTX.Database = TESTDB
	CTX.Logger = logrus.WithFields(logrus.Fields{})
	CTX.APIType = core.USER_API
}

func TestTableNameForUser(t *testing.T) {
	expected := "users"
	ty := reflect.TypeOf(User{})
	actual := core.TableNameFor(ty)
	core.AssertEqual(t, expected, actual)
}


func TestTableNameForConfiguration(t *testing.T) {
	expected := "configurations"
	ty := reflect.TypeOf(Configuration{})
	actual := core.TableNameFor(ty)
	core.AssertEqual(t, expected, actual)
}