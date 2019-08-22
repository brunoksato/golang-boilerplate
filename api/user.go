package api

import (
	"net/http"

	"github.com/brunoksato/golang-boilerplate/core"
	log "github.com/brunoksato/golang-boilerplate/log"
	"github.com/brunoksato/golang-boilerplate/model"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Logout(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"status": "OK"})
}

func Me(c echo.Context) error {
	ctx := ServerContext(c)
	ctx.Payload["results"] = ctx.User
	return c.JSON(http.StatusOK, ctx.Payload)
}

func ChangePassword(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	type customPassword struct {
		Password string `json:"password"`
	}

	u := new(customPassword)
	if err := c.Bind(&u); err != nil {
		return log.AddDefaultError(c, core.NewServerError(err.Error()))
	}

	user := model.User{}
	err := db.First(&user, ctx.User.ID).Error
	if err != nil {
		return log.AddDefaultError(c, core.NewServerError(err.Error()))
	}

	user.HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	serr := db.Model(&user).Set("gorm:save_associations", false).UpdateColumn("hashed_password", user.HashedPassword).Error
	if serr != nil {
		return log.AddDefaultError(c, core.NewServerError(serr.Error()))
	}

	return c.JSON(http.StatusOK, user)
}

func UpdateUser(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	user := model.User{}
	err := db.First(&user, ctx.User.ID).Error
	if err != nil {
		return log.AddDefaultError(c, core.NewNotFoundError(err.Error()))
	}

	if err := c.Bind(&user); err != nil {
		return log.AddDefaultError(c, core.NewServerError(err.Error()))
	}

	dberr := db.
		Model(&user).
		Set("gorm:save_associations", false).
		Updates(map[string]interface{}{
			"name":  user.Name,
			"email": user.Email,
			"phone": user.Phone,
			"image": user.Image,
		}).Error
	if dberr != nil {
		return log.AddDefaultError(c, core.NewServerError(dberr.Error()))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"status": "OK"})
}
