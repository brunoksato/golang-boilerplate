package api_test

import (
	"testing"

	"github.com/brunoksato/golang-boilerplate/core"
	"github.com/brunoksato/golang-boilerplate/model"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserDirects(t *testing.T) {
	setup()
	defer teardown()
	router := router()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)

	u1 := model.User{Name: "User1", Phone: "12982573000", Email: "user1@model.com", Username: "user1", HashedPassword: hashedPassword}
	u2 := model.User{Name: "User2", Phone: "12982573000", Email: "user2@model.com", Username: "user2", HashedPassword: hashedPassword}
	u3 := model.User{Name: "User3", Phone: "12982573000", Email: "user3@model.com", Username: "user3", HashedPassword: hashedPassword}
	u4 := model.User{Name: "User4", Phone: "12982573000", Email: "user4@model.com", Username: "user4", HashedPassword: hashedPassword}
	TESTDB.Create(&u1)
	TESTDB.Create(&u2)
	TESTDB.Create(&u3)
	TESTDB.Create(&u4)
	u5 := model.User{Name: "User5", Email: "user5@model.com", Username: "user5", HashedPassword: hashedPassword}
	u6 := model.User{Name: "User6", Email: "user6@model.com", Username: "user6", HashedPassword: hashedPassword}
	TESTDB.Create(&u5)
	TESTDB.Create(&u6)

	rw, req := core.NewTestRequest("GET", "/api/users/directs")
	router.ServeHTTP(rw, req)
	core.AssertResponseCode(t, rw, 200)

	actualInt := core.JsonToMap(rw.Body.String())
	actual := actualInt["results"].([]interface{})
	assert.Equal(t, 4, len(actual))
	assert.Equal(t, actual[0].(map[string]interface{})["username"], "user1")
	assert.Equal(t, actual[1].(map[string]interface{})["username"], "user2")
	assert.Equal(t, actual[2].(map[string]interface{})["username"], "user3")
	assert.Equal(t, actual[3].(map[string]interface{})["username"], "user4")
}
