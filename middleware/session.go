package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/brunoksato/golang-boilerplate/model"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

func Session(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		db := c.Get("Database").(*gorm.DB)
		user := model.User{}
		isPrivate := strings.Contains(c.Path(), "api")
		isAdmin := strings.Contains(c.Path(), "admin")
		if isPrivate || isAdmin {
			authorization := c.Request().Header.Get("Authorization")
			secretKey := os.Getenv("JWT_KEY_SIGNIN")
			if authorization != "" {
				tokenSlice := strings.Split(authorization, " ")
				if len(tokenSlice) == 2 && tokenSlice[0] == "Bearer" {
					token, err := model.VerifyJWTToken(tokenSlice[1], secretKey)
					if err != nil {
						return echo.NewHTTPError(http.StatusUnauthorized)
					}

					claims := token.Claims.(jwt.MapClaims)
					if token.Valid && claims["iss"] == model.JWT_ISS {
						uidParse := claims["user"].(float64)
						uid := uint(uidParse)
						err := db.
							First(&user, uid).
							Error
						if err != nil {
							return echo.NewHTTPError(http.StatusUnauthorized)
						}

						if isAdmin {
							if !user.IsAdmin() {
								return echo.NewHTTPError(http.StatusUnauthorized)
							}
						}
					} else {
						return echo.NewHTTPError(http.StatusUnauthorized)
					}
				} else {
					return echo.NewHTTPError(http.StatusUnauthorized)
				}
			} else {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}
		} else {
			user.ID = 0
		}

		c.Set("User", user)

		return next(c)
	}
}
