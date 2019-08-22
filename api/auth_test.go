package api_test

import (
	"testing"

	"github.com/brunoksato/golang-boilerplate/core"
	"github.com/stretchr/testify/assert"
)

func TestSignIn(t *testing.T) {
	setup()
	defer teardown()
	router := router()

	u := map[string]interface{}{
		"username": "system",
		"password": "123456",
	}

	rw, req := core.NewTestPost("POST", "/public/signin", u)

	router.ServeHTTP(rw, req)
	core.AssertResponseCode(t, rw, 200)

	actualInt := core.JsonToMap(rw.Body.String())
	actual := actualInt["results"].(map[string]interface{})
	assert.Equal(t, actual["username"], "system")
	assert.Equal(t, actual["email"], "system@model.com")
}

func TestSignInWrong(t *testing.T) {
	setup()
	defer teardown()
	router := router()

	u := map[string]interface{}{
		"email":    "testuser@model.com",
		"password": "password123",
	}

	rw, req := core.NewTestPost("POST", "/public/signin", u)

	router.ServeHTTP(rw, req)
	core.AssertResponseCode(t, rw, 401)
}

func TestSignUp(t *testing.T) {
	setup()
	defer teardown()
	router := router()

	u := map[string]interface{}{
		"name":     "bruno sato",
		"username": "brunoksato",
		"email":    "brunosato@model.com",
		"password": "password",
		"phone":    "12982575000",
	}

	rw, req := core.NewTestPost("POST", "/public/signup", u)

	router.ServeHTTP(rw, req)
	core.AssertResponseCode(t, rw, 201)

	actualInt := core.JsonToMap(rw.Body.String())
	actual := actualInt["results"].(map[string]interface{})
	assert.Equal(t, actual["email"], "brunosato@model.com")
	assert.Equal(t, actual["name"], "bruno sato")
	assert.Equal(t, actual["username"], "brunoksato")
	assert.Equal(t, actual["phone"], "12982575000")
}
