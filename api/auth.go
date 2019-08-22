package api

import (
	"net/http"
	"os"

	"github.com/brunoksato/golang-boilerplate/core"
	log "github.com/brunoksato/golang-boilerplate/log"
	"github.com/brunoksato/golang-boilerplate/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database
	tx := db.Begin()

	if err := tx.Error; err != nil {
		tx.Rollback()
		return log.AddDefaultError(c, core.NewServerError(err.Error()))
	}

	user := model.User{}
	if err := c.Bind(&user); err != nil {
		tx.Rollback()
		return log.AddDefaultError(c, core.NewServerError(err.Error()))
	}

	cerr := user.Create(ArgonContextTransaction(c, tx), ctx.User)
	if cerr != nil {
		tx.Rollback()
		return log.AddDefaultError(c, cerr)
	}

	err := tx.Where("email = ?", user.Email).Find(&user).Error
	if err != nil {
		tx.Rollback()
		return log.AddDefaultError(c, core.NewServerError(err.Error()))
	}

	tx.Commit()

	ctx.Payload["results"] = user
	return c.JSON(http.StatusCreated, ctx.Payload)
}

func SignIn(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	u := model.User{}
	if err := c.Bind(&u); err != nil {
		return log.AddDefaultError(c, core.NewServerError(err.Error()))
	}

	if u.Username != "" {
		if err := db.Where("username = ?", u.Username).Find(&ctx.User).Error; err != nil {
			return log.AddDefaultError(c, core.NewServerError(err.Error()))
		}

		if ctx.APIType == core.ADMIN_API && !ctx.User.IsAdmin() {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"status": "Not Authorized"})
		}

		if ok, _ := ctx.User.VerifyPassword(u.Password); ok {
			expireAt := model.JWTTokenExpirationDate()

			jwt, dberr := model.IssueJWToken(ctx.User.ID, []string{"user"}, expireAt)
			if dberr != nil {
				return log.AddDefaultError(c,
					core.NewServerError(
						dberr.Error(),
						map[string]interface{}{
							"user_id": ctx.User.ID,
						},
					),
				)
			}

			ctx.Payload["results"] = ctx.User
			ctx.Payload["token"] = jwt
			return c.JSON(http.StatusOK, ctx.Payload)
		}
	}

	return c.JSON(http.StatusUnauthorized, map[string]interface{}{"status": "Not Authorized"})
}

func RecoverPassword(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	email := c.Param("email")

	user := model.User{}
	count := 0

	err := db.Where("LOWER(email) = LOWER(?)", email).Find(&user).Count(&count).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": "Email Not Found"})
	}

	if count > 0 {
		expireAt := model.JWTTokenExpirationDate()
		jwt, dberr := model.IssueJWTTokenForEmail(user.ID, user.Email, expireAt)
		if dberr != nil {
			return log.AddDefaultError(c,
				core.NewServerError(
					dberr.Error(),
					map[string]interface{}{
						"user_id": ctx.User.ID,
					},
				),
			)
		}

		user.ResetPasswordEmail(jwt)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"status": "OK"})
}

func ChangePasswordExternal(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	type customPasswordExternal struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	u := new(customPasswordExternal)
	if err := c.Bind(&u); err != nil {
		return log.AddDefaultError(c, core.NewServerError(err.Error()))
	}

	secretKey := os.Getenv("JWT_KEY_EMAIL")
	token, err := model.VerifyJWTToken(u.Token, secretKey)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Token invalid"})
	}

	if err != nil {
		log.LoggerForParams(c,
			map[string]interface{}{
				"user_id": ctx.User.ID,
				"status":  "Token invalid",
				"code":    http.StatusBadRequest,
			},
		).Info("ChangePasswordExternal: Token inválido")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Token invalid"})
	}

	claims := token.Claims.(jwt.MapClaims)
	if token.Valid && claims["iss"] == model.JWT_ISS {
		email := claims["email"].(string)

		user := model.User{}
		err := db.Where("email = ?", email).First(&user).Error
		if err != nil {
			return log.AddDefaultError(c, core.NewServerError(err.Error()))
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		uerr := db.Model(&user).Set("gorm:save_associations", false).UpdateColumn("hashed_password", hashedPassword).Error
		if uerr != nil {
			return log.AddDefaultError(c, core.NewServerError(uerr.Error()))
		}
	} else {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Token invalid or expired"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"status": "OK"})
}
