package model

import (
	"strings"
	"testing"

	"github.com/brunoksato/golang-boilerplate/core"
)

func TestUserValidateForCreate(t *testing.T) {
	u := User{}
	u.Name = "Bruno Sato"
	u.Username = "brunoksato"
	u.Email = "bruno@model.com"
	u.Password = "123456"
	err := u.ValidateForCreate()
	core.AssertNoError(t, err)

	u.Email = "32"
	err = u.ValidateForCreate()
	core.AssertBusinessError(t, "email: 32 does not validate as email;", err)
	core.AssertEqual(t, core.ERROR_SUBCODE_EMAIL_FORMAT, err.Subcode())

	u.Email = "bruno@model.com"
	err = u.ValidateForCreate()
	core.AssertNoError(t, err)

	u.Name = "u2"
	err = u.ValidateForCreate()
	core.AssertBusinessError(t, "name: u2 does not validate as length(3|255);", err)
	core.AssertEqual(t, core.ERROR_SUBCODE_NAME_LENGTH, err.Subcode())

	u.Name = "Proper_Username"
	u.Password = "N"
	err = u.ValidateForCreate()
	core.AssertEqual(t, "password: N does not validate as length(5|64);", err.Error())
	core.AssertEqual(t, core.ERROR_SUBCODE_PASSWORD_LENGTH, err.Subcode())

	u.Password = "A valid passWord"

	u.Username = "u1"
	err = u.ValidateForCreate()
	core.AssertBusinessError(t, "username: u1 does not validate as length(3|15);", err)
	core.AssertEqual(t, core.ERROR_SUBCODE_USERNAME_LENGTH, err.Subcode())

	u.Username = "brunoksato"

	err = u.ValidateForCreate()
	core.AssertNil(t, err)
}

func TestUserValidateForUpdate(t *testing.T) {
	u := User{}
	u.Name = "Bruno Sato"
	u.Username = "brunoksato"
	u.Email = "bruno@model.com"
	u.Password = "123456"
	err := u.ValidateForUpdate()
	core.AssertNoError(t, err)
}

func TestUserValidateForDelete(t *testing.T) {
	u := User{}
	err := u.ValidateForDelete(CTX)
	core.AssertNoError(t, err)
}

func TestUserRestrictor(t *testing.T) {
	setupDB()
	defer teardownDB()

	owner := User{Model: Model{ID: 1}, Name: "User 1", Email: "user1@model.com"}
	notOwner := User{Model: Model{ID: 1}, Name: "User 1", Email: "user1@model.com"}
	view := User{Model: Model{ID: 2}, Name: "User 2", Email: "user2@model.com"}

	AssertUserCan(t, owner.UserCanView, CTX, view)
	AssertUserCan(t, owner.UserCanCreate, CTX, view)
	AssertUserCantUpdate(t, owner.UserCanUpdate, CTX, view, []string{})
	AssertUserCant(t, owner.UserCanDelete, CTX, view)

	AssertUserCan(t, owner.UserCanView, CTX, owner)
	AssertUserCan(t, owner.UserCanCreate, CTX, owner)
	AssertUserCanUpdate(t, owner.UserCanUpdate, CTX, owner, []string{})
	AssertUserCant(t, owner.UserCanDelete, CTX, owner)

	AssertUserCan(t, notOwner.UserCanView, CTX, view)
	AssertUserCan(t, notOwner.UserCanCreate, CTX, view)
	AssertUserCantUpdate(t, notOwner.UserCanUpdate, CTX, view, []string{})
	AssertUserCant(t, notOwner.UserCanDelete, CTX, view)
}

func TestByUserEmail(t *testing.T) {
	setupDB()
	defer teardownDB()

	u1 := User{Name: "User1", Email: "user1@model.com", Username: "user1"}
	u2 := User{Name: "User2", Email: "user2@model.com", Username: "user2"}
	u3 := User{Name: "User3", Email: "user3@model.com", Username: "user3"}
	u4 := User{Name: "User4", Email: "user4@model.com", Username: "user4"}
	u5 := User{Name: "User5"}
	TESTDB.Create(&u1)
	TESTDB.Create(&u2)
	TESTDB.Create(&u3)
	TESTDB.Create(&u4)
	TESTDB.Create(&u5)

	u := User{}
	err := TESTDB.Scopes(ByUserEmail(u1.Email)).First(&u).Error
	core.AssertNoError(t, err)
	core.AssertEqual(t, u1.ID, u.ID)

	u = User{}
	err = TESTDB.Scopes(ByUserEmail(strings.ToUpper(u2.Email))).First(&u).Error
	core.AssertNoError(t, err)
	core.AssertEqual(t, u2.ID, u.ID)

	u = User{}
	err = TESTDB.Scopes(ByUserEmail(strings.ToUpper(u3.Email))).First(&u).Error
	core.AssertNoError(t, err)
	core.AssertEqual(t, u3.ID, u.ID)

	u = User{}
	err = TESTDB.Scopes(ByUserEmail(strings.ToLower(u4.Email))).First(&u).Error
	core.AssertNoError(t, err)
	core.AssertEqual(t, u4.ID, u.ID)

	u = User{}
	err = TESTDB.Scopes(ByUserEmail("USeR_5")).First(&u).Error
	core.AssertEqual(t, "record not found", err.Error())
}

func TestByUserUsername(t *testing.T) {
	setupDB()
	defer teardownDB()

	u1 := User{Name: "User1", Email: "user1@model.com", Username: "user1"}
	u2 := User{Name: "User2", Email: "user2@model.com", Username: "user2"}
	u3 := User{Name: "User3", Email: "user3@model.com", Username: "user3"}
	u4 := User{Name: "User4", Email: "user4@model.com", Username: "user4"}
	u5 := User{Name: "User5"}
	TESTDB.Create(&u1)
	TESTDB.Create(&u2)
	TESTDB.Create(&u3)
	TESTDB.Create(&u4)
	TESTDB.Create(&u5)

	u := User{}
	err := TESTDB.Scopes(ByUserUsername(u1.Username)).First(&u).Error
	core.AssertNoError(t, err)
	core.AssertEqual(t, u1.ID, u.ID)
}

func TestCreateUser(t *testing.T) {
	setupDB()
	defer teardownDB()

	us := []User{}
	TESTDB.Find(&us)
	core.AssertEqual(t, 0, len(us))

	u := User{}
	u.Name = "Bruno Sato"
	u.Username = "brunoksato"
	u.Email = "bruno@model.com"
	u.Password = "123456"
	err := u.Create(CTX, User{})
	core.AssertNoError(t, err)

	us = []User{}
	TESTDB.Find(&us)
	core.AssertEqual(t, 1, len(us))

	u = User{Name: "User1", Email: "bruno@model.com", Username: "user1"}
	err = u.Create(CTX, User{})
	core.AssertBusinessError(t, "email: already used;", err)
	core.AssertEqual(t, core.ERROR_SUBCODE_EMAIL_TAKEN, err.Subcode())

	u = User{Name: "User1", Email: "bruno1@model.com", Username: "brunoksato"}
	err = u.Create(CTX, User{})
	core.AssertBusinessError(t, "username: already used;", err)
	core.AssertEqual(t, core.ERROR_SUBCODE_EMAIL_TAKEN, err.Subcode())
}

func TestUserIsAdmin(t *testing.T) {
	setupDB()
	defer teardownDB()

	u := User{}
	u.Email = "bruno@model.com"
	u.Name = "bruno"
	u.Password = "password"
	TESTDB.Create(&u)
	core.AssertFalse(t, u.IsAdmin())

	u.Admin = true
	TESTDB.Save(&u)
	core.AssertTrue(t, u.IsAdmin())
}

func TestUserIsUser(t *testing.T) {
	setupDB()
	defer teardownDB()

	u := User{}
	u.Email = "bruno@model.com"
	u.Name = "bruno"
	u.Password = "password"
	u.Admin = true
	TESTDB.Create(&u)
	core.AssertFalse(t, u.IsUser())

	u.Admin = false
	TESTDB.Save(&u)
	core.AssertTrue(t, u.IsUser())
}
