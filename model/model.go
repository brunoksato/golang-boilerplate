package model

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/brunoksato/golang-boilerplate/core"
	"github.com/brunoksato/golang-boilerplate/util"
	"github.com/jinzhu/gorm"
)

type Model struct {
	ID        uint      `json:"id" gorm:"primary_key" settable:"false"`
	CreatedAt time.Time `json:"created_at" sql:"DEFAULT:current_timestamp" settable:"false"`
	UpdatedAt time.Time `json:"updated_at" sql:"DEFAULT:current_timestamp" settable:"false"`
}

func (m Model) OrderBy(db *gorm.DB) *gorm.DB {
	order := "\"created_at\" DESC"
	return db.Order(order)
}

func ParentIdField(t reflect.Type) (field *reflect.StructField, dbFieldName string) {
	elemT := t
	if elemT.Kind() == reflect.Ptr {
		elemT = elemT.Elem()
	}
	for i := 0; i < elemT.NumField(); i++ {
		tag := elemT.Field(i).Tag
		if tag.Get("parent_field") != "" {
			dbFieldName = tag.Get("parent_field")
			fieldRef := elemT.Field(i)
			field = &fieldRef
		}
	}
	return
}

func UserIDField(t reflect.Type) (field *reflect.StructField, dbFieldName string) {
	elemT := t
	if elemT.Kind() == reflect.Ptr {
		elemT = elemT.Elem()
	}

	f, found := elemT.FieldByName("UserID")
	if found {
		field = &f
		dbFieldName = "user_id"
	}
	return
}

func OrderField(t reflect.Type) (field *reflect.StructField, dbFieldName string) {
	elemT := t
	if elemT.Kind() == reflect.Ptr {
		elemT = elemT.Elem()
	}
	for i := 0; i < elemT.NumField(); i++ {
		tag := elemT.Field(i).Tag
		if tag.Get("order_field") != "" {
			dbFieldName = tag.Get("order_field")
			fieldRef := elemT.Field(i)
			field = &fieldRef
		}
	}
	return
}

func SetUserID(item interface{}, id uint) error {
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return errors.New("Cannot set value on nil item")
		}
		v = reflect.ValueOf(item).Elem()
	}
	t := v.Type()
	if !TypeHasUserField(t) {
		return errors.New("Type does not have a reference to a user")
	}
	userField, _ := UserIDField(t)
	f := v.FieldByName(userField.Name)
	if f.Type().Kind() == reflect.Ptr {
		f.Set(reflect.ValueOf(&id))
	} else {
		f.Set(reflect.ValueOf(id))
	}
	return nil
}

func TypeHasParentField(t reflect.Type) bool {
	parentField, _ := ParentIdField(t)
	return parentField != nil
}

func TypeHasUserField(t reflect.Type) bool {
	userField, _ := UserIDField(t)
	return userField != nil
}

func AssertUserCan(t *testing.T, m func(ctx *ModelCtx, u User) (bool, core.DefaultError), ctx *ModelCtx, user User) {
	can, err := m(ctx, user)
	if err != nil {
		t.Errorf("assertUserCan: Error occurred - %s (caller: %s)", err.Error(), util.CallerInfo())
	}
	if !can {
		t.Errorf("assertUserCan: User cannot (caller: %s)", util.CallerInfo())
	}
}

func AssertUserCant(t *testing.T, m func(ctx *ModelCtx, u User) (bool, core.DefaultError), ctx *ModelCtx, user User) {
	can, err := m(ctx, user)
	if err != nil {
		t.Errorf("assertUserCant: Error occurred - %s (caller: %s)", err.Error(), util.CallerInfo())
	}
	if can {
		t.Errorf("assertUserCant: User can (caller: %s)", util.CallerInfo())
	}
}

func AssertUserCanUpdate(t *testing.T, m func(ctx *ModelCtx, u User, s []string) (bool, core.DefaultError), ctx *ModelCtx, user User, fields []string) {
	can, err := m(ctx, user, fields)
	if err != nil {
		t.Errorf("assertUserCanUpdate: Error occurred - %s (caller: %s)", err.Error(), util.CallerInfo())
	}
	if !can {
		t.Errorf("assertUserCanUpdate: User cannot (caller: %s)", util.CallerInfo())
	}
}

func AssertUserCantUpdate(t *testing.T, m func(ctx *ModelCtx, u User, s []string) (bool, core.DefaultError), ctx *ModelCtx, user User, fields []string) {
	can, err := m(ctx, user, fields)
	if err != nil {
		t.Errorf("assertUserCantUpdate: Error occurred - %s (caller: %s)", err.Error(), util.CallerInfo())
	}
	if can {
		t.Errorf("assertUserCantUpdate: User can (caller: %s)", util.CallerInfo())
	}
}
