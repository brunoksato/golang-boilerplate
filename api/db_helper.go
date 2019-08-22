package api

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/brunoksato/golang-boilerplate/model"
	"github.com/brunoksato/golang-boilerplate/core"
	"github.com/jinzhu/gorm"
)

func OrderByFor(db *gorm.DB, t reflect.Type) *gorm.DB {
	if model.IsSorter(t) {
		item := reflect.New(t).Interface()
		sorter := item.(model.Sorter)
		return sorter.OrderBy(db)
	}

	tableName := db.NewScope(reflect.New(t).Interface()).TableName()
	tableName = fmt.Sprintf("\"%s\".", tableName)
	order := tableName + "\"id\" ASC"

	_, orderField := model.OrderField(t)
	if orderField != "" {
		order = fmt.Sprintf("%s\"%s\" ASC, %s", tableName, orderField, order)
	}

	return db.Order(order)
}

func JoinsFor(ctx *Context, db *gorm.DB, parentIsSpecific bool) *gorm.DB {
	return db.Scopes(core.DefaultPreloads(db, ctx.Type, ctx.APIType, parentIsSpecific))
}

func ScopesFor(ctx *Context, path string, userID uint, db *gorm.DB) *gorm.DB {
	parts := strings.Split(path, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		switch part {
		//change the name of model
		case "models":
			db = db.Where("user_id = ?", userID)
		}

	}
	return db
}

func FilterByParentFor(db *gorm.DB, pt, t reflect.Type, parentID uint) *gorm.DB {
	userType := reflect.TypeOf(model.User{})

	if pt == userType {
		db = FilterByUserFor(db, pt, t, parentID)
	} else if pt != nil {
		parentField := gorm.ToDBName(pt.Name())
		db = db.Where(parentField+"_id = ?", parentID)
	} else if model.TypeHasParentField(t) {
		_, parentField := model.ParentIdField(t)
		db = db.Where(parentField+" = ?", parentID)
	}

	return db
}

func FilterByUserFor(db *gorm.DB, pt, t reflect.Type, userID uint) *gorm.DB {
	var zeroType reflect.Type

	ptName := ""
	if pt != zeroType {
		ptName = pt.Name()
	}

	switch ptName {
	case "", "User":
		switch t.Name() {
		default:
			_, userIDField := model.UserIDField(t)
			if userIDField == "user_id" {
				db = db.Where("user_id = ?", userID)
			}
		}
	}

	return db
}
