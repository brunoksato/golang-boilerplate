package api

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/brunoksato/golang-boilerplate/model"
	"github.com/brunoksato/golang-boilerplate/core"
	log "github.com/brunoksato/golang-boilerplate/log"
	"github.com/brunoksato/golang-boilerplate/util"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

func List(c echo.Context) error {
	ctx := ServerContext(c)

	db := ctx.Database

	db, err := DefaultListQuery(c, ctx, db)
	if err != nil {
		return log.AddDefaultError(c, err)
	}

	db = DefaultJoins(c, ctx, db)
	db = DefaultScopes(c, ctx, db)
	db = DefaultPaging(c, ctx, db)
	db = DefaultOrder(c, ctx, db)

	err = AddListToPayload(ctx, db)
	if err != nil {
		return log.AddDefaultError(c, err)
	}

	return c.JSON(http.StatusOK, ctx.Payload)
}

func Create(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	item := reflect.New(ctx.Type).Interface()
	if err := c.Bind(item); err != nil {
		return log.AddDefaultError(c, core.NewServerError(err.Error()))
	}

	userID := ActiveUserID(c, ctx)
	if userID > 0 {
		model.SetUserID(item, uint(userID))
	}

	if model.IsCreator(reflect.PtrTo(ctx.Type)) {
		creator := item.(model.Creator)
		err := creator.Create(ArgonContext(c), ctx.User)
		if err != nil {
			return log.AddDefaultError(c, err)
		}
		item = creator
	} else {
		err := DefaultValidationForCreate(c, ctx, item)
		if err != nil {
			return log.AddDefaultError(c, err)
		}

		dberr := db.Set("gorm:save_associations", false).Create(item).Error
		if dberr != nil {
			return log.AddDefaultError(c, core.NewServerError(dberr.Error()))
		}
	}

	db.First(item)

	ctx.Payload["results"] = item
	return c.JSON(http.StatusCreated, ctx.Payload)
}

func Get(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return log.AddDefaultError(c, core.NewNotFoundError(err.Error()))
	}

	item := reflect.New(ctx.Type).Interface()

	if model.IsGetter(ctx.Type) && err != nil {
		var merr core.DefaultError
		getter := item.(model.Getter)
		item, merr = getter.GetByID(ArgonContext(c), ctx.User, uint(id))
		if merr != nil {
			return log.AddDefaultError(c, merr)
		}
	} else {
		db = DefaultJoins(c, ctx, db)
		db = DefaultScopes(c, ctx, db)
		err = db.First(item, id).Error

		if err != nil {
			return log.AddDefaultError(c, core.NewNotFoundError(err.Error()))
		}

		merr := DefaultValidationForGet(c, item)
		if merr != nil {
			return log.AddDefaultError(c, merr)
		}
	}

	ctx.Payload["results"] = item
	return c.JSON(http.StatusOK, ctx.Payload)
}

func Update(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return log.AddDefaultError(c, core.NewNotFoundError(err.Error()))
	}

	item := reflect.New(ctx.Type).Interface()
	err = db.First(item, id).Error
	if err != nil {
		return log.AddDefaultError(c, core.NewNotFoundError(err.Error()))
	}

	if err := c.Bind(item); err != nil {
		return log.AddDefaultError(c, core.NewServerError(err.Error()))
	}

	requestMap := core.ModelToJsonMap(item)

	fields := make([]string, 0)
	for fieldName, newVal := range requestMap {
		typeField, _ := core.GetFieldByJsonTag(item, fieldName)
		if typeField != nil && shouldUpdateField(*typeField) {
			core.SetByJsonTag(item, fieldName, newVal)
			fields = append(fields, typeField.Name)
		}
	}

	if model.IsUpdater(reflect.PtrTo(ctx.Type)) {
		updater := item.(model.Updater)
		merr := updater.Update(ArgonContext(c), ctx.User)
		if merr != nil {
			return log.AddDefaultError(c, merr)
		}
		item = updater
	} else {
		merr := DefaultValidationForUpdate(c, ctx, item, fields)
		if merr != nil {
			return log.AddDefaultError(c, merr)
		}

		dberr := db.Set("gorm:save_associations", false).Save(item).Error
		if dberr != nil {
			return log.AddDefaultError(c, core.NewServerError("Error saving data: "+dberr.Error()))
		}
	}

	db.First(item)

	ctx.Payload["results"] = item
	return c.JSON(http.StatusAccepted, ctx.Payload)
}

func Delete(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return log.AddDefaultError(c, core.NewNotFoundError(err.Error()))
	}

	item := reflect.New(ctx.Type).Interface()
	err = db.First(item, id).Error
	if err != nil {
		return log.AddDefaultError(c, core.NewNotFoundError(err.Error()))
	}

	var merr core.DefaultError
	if model.IsDeleter(reflect.PtrTo(ctx.Type)) {
		deleter := item.(model.Deleter)
		merr = deleter.Delete(ArgonContext(c), ctx.User)
	} else {
		merr := DefaultValidationForDelete(c, ctx, item)
		if merr != nil {
			return log.AddDefaultError(c, merr)
		}

		merr = model.DefaultDelete(ArgonContext(c), item)
		if merr != nil {
			return log.AddDefaultError(c, merr)
		}
	}

	if merr != nil {
		return log.AddDefaultError(c, merr)
	}

	ctx.Payload["results"] = item
	return c.JSON(http.StatusOK, ctx.Payload)
}

func DefaultListQuery(c echo.Context, ctx *Context, db *gorm.DB) (*gorm.DB, core.DefaultError) {
	var parentID, userID int
	strParentID := c.Param("parentId")
	if strParentID != "" {
		var err error
		parentID, err = strconv.Atoi(strParentID)
		if err != nil {
			return db, core.NewNotFoundError(fmt.Sprintf("Invalid id: %s", strParentID))
		}
		if parentID <= 0 {
			return db, core.NewNotFoundError(fmt.Sprintf("Invalid id: %s", strParentID))
		}
	}
	strUserID := c.Param("userId")
	if strUserID != "" {
		var err error
		userID, err = strconv.Atoi(strUserID)
		if err != nil {
			return db, core.NewNotFoundError(fmt.Sprintf("Invalid id: %s", strUserID))
		}
		if userID <= 0 {
			return db, core.NewNotFoundError(fmt.Sprintf("Invalid id: %s", strUserID))
		}
	}

	if parentID > 0 {
		db = FilterByParentFor(db, ctx.ParentType, ctx.Type, uint(parentID))
	}

	if userID > 0 {
		db = FilterByUserFor(db, ctx.ParentType, ctx.Type, uint(userID))
	}

	switch ctx.APIType {
	case core.USER_API:
		userID := ActiveUserID(c, ctx)
		if ctx.ParentType == reflect.TypeOf(model.User{}) {
			db = FilterByUserFor(db, ctx.ParentType, ctx.Type, userID)
		}
	case core.ADMIN_API:
		// no filter
	}
	return db, nil
}

func DefaultJoins(c echo.Context, ctx *Context, db *gorm.DB) *gorm.DB {
	parentID, _ := strconv.Atoi(c.Param("parentId"))

	if parentID > 0 {
		db = JoinsFor(ctx, db, true)
	} else {
		db = JoinsFor(ctx, db, false)
	}
	return db
}

func DefaultScopes(c echo.Context, ctx *Context, db *gorm.DB) *gorm.DB {
	userID := ActiveUserID(c, ctx)
	path := c.Path()
	db = ScopesFor(ctx, path, userID, db)
	return db
}

func DefaultPaging(c echo.Context, ctx *Context, db *gorm.DB, opts ...bool) *gorm.DB {
	queryTC := true
	if len(opts) > 0 {
		queryTC = opts[0]
	}

	st := c.QueryParam("start")
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if limit > 0 && queryTC {
		queryTotalCount(ctx, db)
	}

	if st != "" {
		startIdx, _ := strconv.Atoi(st)
		if startIdx > 0 {
			db = db.Offset(startIdx)
		}
	}

	if limit > 0 {
		db = limitQueryByConfig(ctx, db, "", limit)
	}

	return db
}

func queryTotalCount(ctx *Context, db *gorm.DB) {
	item := reflect.New(ctx.Type).Interface()
	var n int

	db.Model(item).
		Select("COUNT(*)").
		Row().
		Scan(&n)

	ctx.Payload["ct"] = n
}

func limitQueryByConfig(c *Context, db *gorm.DB, key string, requestLimit int) *gorm.DB {
	dbLimit := requestLimit
	limitStr := os.Getenv(key)
	limit, err := strconv.Atoi(limitStr)
	if err == nil {
		if dbLimit <= 0 || (limit > 0 && limit < dbLimit) {
			dbLimit = limit
		}
	}
	if dbLimit > 0 {
		db = db.Limit(dbLimit)
	}
	return db
}

func DefaultOrder(c echo.Context, ctx *Context, db *gorm.DB) *gorm.DB {
	sort := c.FormValue("sort")
	if sort != "" {
		fields, ascending := util.ConvertQueryTermToOrderTerm(sort)
		tableName := db.NewScope(reflect.New(ctx.Type).Interface()).TableName()
		order := ""
		for i, field := range fields {
			term := ""
			switch field {
			case "email":
				term = fmt.Sprintf("\"%s\".email", tableName)
			case "id":
				term = fmt.Sprintf("\"%s\".id", tableName)
			case "created_at":
				term = fmt.Sprintf("\"%s\".id", tableName)
			case "updated_at":
				term = fmt.Sprintf("\"%s\".id", tableName)
			default:
				item := reflect.New(ctx.Type).Interface()
				structField, err := core.GetFieldByJsonTag(item, field)
				if err == nil {
					term = fmt.Sprintf("\"%s\".%s", tableName, gorm.ToDBName(structField.Name))
				} else {
					continue
				}
			}

			if ascending[i] {
				term = fmt.Sprintf("%s ASC", term)
			} else {
				term = fmt.Sprintf("%s DESC", term)
			}

			order = fmt.Sprintf("%s %s,", order, term)
		}

		order = fmt.Sprintf("%s %s.id ASC", order, tableName)
		db = db.Order(order)
	} else {
		db = OrderByFor(db, ctx.Type)
	}
	return db
}

func AddListToPayload(ctx *Context, db *gorm.DB) core.DefaultError {
	result, err := getListFromQuery(ctx, db)
	if err != nil {
		if err.IsWarning() {
			ctx.Payload["results"] = result
		}
		return err
	}

	ctx.Payload["results"] = result
	return nil
}

func DefaultValidationForGet(c echo.Context, item interface{}) core.DefaultError {
	ctx := ServerContext(c)

	switch ctx.Type.Name() {
	case "Example":
		// Don't restrict the basic types
	default:
		if model.IsRestrictor(ctx.Type) {
			restrictor := item.(model.Restrictor)
			canView, err := restrictor.UserCanView(ArgonContext(c), ctx.User)
			if err != nil {
				return err
			}
			if !canView {
				return core.NewPermissionError("You do not have permission",
					core.ERROR_SUBCODE_USER_LACKS_PERMISSION)
			}
		}
	}
	return nil
}

func DefaultValidationForCreate(c echo.Context, ctx *Context, item interface{}) core.DefaultError {
	if model.IsRestrictor(ctx.Type) {
		restrictor := item.(model.Restrictor)
		canCreate, err := restrictor.UserCanCreate(ArgonContext(c), ctx.User)
		if err != nil {
			return err
		}
		if !canCreate {
			return core.NewPermissionError("You do not have permission",
				core.ERROR_SUBCODE_USER_LACKS_PERMISSION)
		}
	}

	if model.IsValidator(ctx.Type) {
		validator := item.(model.Validator)
		return validator.ValidateForCreate()
	}

	return nil
}

func DefaultValidationForUpdate(c echo.Context, ctx *Context, item interface{}, fields []string) core.DefaultError {
	if model.IsRestrictor(ctx.Type) {
		restrictor := item.(model.Restrictor)
		canUpdate, err := restrictor.UserCanUpdate(ArgonContext(c), ctx.User, fields)
		if err != nil {
			return err
		}
		if !canUpdate {
			return core.NewPermissionError("You do not have permission to make these changes",
				core.ERROR_SUBCODE_USER_LACKS_PERMISSION)
		}
	}

	if model.IsValidator(ctx.Type) {
		err := model.ValidateStructFields(item, fields)
		if err != nil {
			return core.NewBusinessError(err.Error())
		}
		return nil
	}

	return nil
}

func DefaultValidationForDelete(c echo.Context, ctx *Context, item interface{}) core.DefaultError {
	if model.IsRestrictor(ctx.Type) {
		restrictor := item.(model.Restrictor)
		canDelete, err := restrictor.UserCanDelete(ArgonContext(c), ctx.User)
		if err != nil {
			return err
		}
		if !canDelete {
			return core.NewPermissionError("You do not have permission",
				core.ERROR_SUBCODE_USER_LACKS_PERMISSION)
		}
	}

	if model.IsValidator(ctx.Type) {
		validator := item.(model.Validator)
		return validator.ValidateForDelete(ArgonContext(c))
	}

	return nil
}

func getListFromQuery(ctx *Context, db *gorm.DB) (interface{}, core.DefaultError) {
	var err error

	items := util.NewSliceForType(ctx.Type)
	err = db.Find(items).Error
	if err != nil {
		return nil, core.NewServerError(err.Error())
	}

	items = util.ItemsOrEmptySlice(ctx.Type, items)

	return items, nil
}

func shouldUpdateField(field reflect.StructField) bool {
	tag := field.Tag
	if tag.Get("settable") == "false" {
		return false
	}

	timeType := reflect.TypeOf(time.Time{})
	timestampType := reflect.TypeOf(core.Timestamp{})
	nullableTimestampType := reflect.TypeOf(core.NullableTimestamp{})

	should := true
	typ := field.Type
	kind := typ.Kind()

	if kind == reflect.Ptr {
		typ = field.Type.Elem()
		kind = typ.Kind()
	}

	should = should && kind != reflect.Array
	should = should && kind != reflect.Slice
	should = should && kind != reflect.Struct
	should = should || typ == timeType
	should = should || typ == timestampType
	should = should || typ == nullableTimestampType
	return should
}
