package model

import (
	"os"
	"testing"
	"time"

	"github.com/brunoksato/golang-boilerplate/core"
)

func TestJWTTokenExpirationDate(t *testing.T) {
	jwtTime := time.Now().Add(time.Hour * 24).Round(time.Second)
	core.AssertEqual(t, jwtTime, JWTTokenExpirationDate().Round(time.Second))

	os.Setenv("JWT_TOKEN_EXPIRATION", "0")
	jwtTime = time.Now().Add(time.Hour * 0).Round(time.Second)
	core.AssertEqual(t, jwtTime, JWTTokenExpirationDate().Round(time.Second))

	os.Setenv("JWT_TOKEN_EXPIRATION", "2")
	jwtTime = time.Now().Add(time.Hour * 2).Round(time.Second)
	core.AssertEqual(t, jwtTime, JWTTokenExpirationDate().Round(time.Second))

	os.Setenv("JWT_TOKEN_EXPIRATION", "") // defaults to 24h if the conversion is wrong
	jwtTime = time.Now().Add(time.Hour * 24).Round(time.Second)
	core.AssertEqual(t, jwtTime, JWTTokenExpirationDate().Round(time.Second))

	os.Setenv("JWT_TOKEN_EXPIRATION", "24")
}
