package core

import (
	"errors"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
)

func DefaultPreloads(db *gorm.DB, t reflect.Type, api APIType, skipParent bool) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		elemT := t
		if elemT.Kind() == reflect.Ptr {
			elemT = elemT.Elem()
		}
		for i := 0; i < elemT.NumField(); i++ {
			f := elemT.Field(i)
			ft := f.Type
			if ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
			}

			preload := false
			all := false
			skip := false
			fieldName := f.Name
			tag := elemT.Field(i).Tag

			if skipParent {
				if _, ok := tag.Lookup("parent"); ok {
					continue
				}
			}

			fetch, _ := tag.Lookup("fetch")
			fetchConfigs := strings.Split(fetch, ",")
			for _, config := range fetchConfigs {
				switch config {
				case "user":
					if api == USER_API {
						preload = true
					}
				case "admin":
					if api == ADMIN_API {
						preload = true
					}
				case "eager":
					preload = true
				case "all":
					all = true
				case "parent":
					skip = skipParent
				default:
					fieldName = config
				}
			}
			if skip {
				preload = false
			}
			if preload {
				if all {
					db = db.Preload(fieldName, func(db *gorm.DB) *gorm.DB {
						return db.Unscoped()
					})
				} else {
					db = db.Preload(fieldName)
				}
			}

		}
		return db
	}
}

func PageQueryResults(start, limit *uint) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if start != nil && *start != 0 {
			db = db.Offset(*start)
		}
		if limit != nil && *limit != 0 {
			db = db.Limit(*limit)
		}
		return db
	}
}

func UpdateField(db *gorm.DB, item interface{}, field string, value interface{}) DefaultError {
	if item == nil {
		return NewServerError("Trying to update field on nil object")
	}

	if reflect.ValueOf(item).Kind() != reflect.Ptr {
		return NewServerError("Trying to update field on non-pointer object")
	}

	id, err := GetID(item)
	if err != nil {
		return NewServerError(err.Error())
	}

	if id > 0 {
		err = db.Set("gorm:save_associations", false).Model(item).Update(field, value).Error
		if err != nil {
			return NewServerError(err.Error())
		}
		return nil
	}

	return NewServerError("Trying to update field on empty object")
}

func UpdateFields(db *gorm.DB, item interface{}, transitive interface{}) DefaultError {
	if item == nil {
		return NewServerError("Trying to update fields on nil object")
	}

	if reflect.ValueOf(item).Kind() != reflect.Ptr {
		return NewServerError("Trying to update fields on non-pointer object")
	}

	id, err := GetID(item)
	if err != nil {
		return NewServerError(err.Error())
	}

	if id > 0 {
		err = db.Set("gorm:save_associations", false).Model(item).Updates(transitive).Error
		if err != nil {
			return NewServerError(err.Error())
		}
		return nil
	}

	return NewServerError("Trying to update fields on empty object")
}

func GetID(item interface{}) (uint, error) {
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return 0, errors.New("Cannot get value from nil item")
		}
		v = reflect.ValueOf(item).Elem()
	}
	f := v.FieldByName("ID")
	parentID := uint(f.Uint())
	return parentID, nil
}

func TableNameFor(t reflect.Type) string {
	name := gorm.ToDBName(t.Name())
	if strings.HasSuffix(name, "y") {
		name = name[:len(name)-1] + "ies"
	} else if strings.HasSuffix(name, "ss") {
		name = name + "es"
	} else {
		name = name + "s"
	}
	return name
}
