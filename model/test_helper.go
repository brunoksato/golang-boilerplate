package model

import (
	"fmt"
	"os"

	"github.com/brunoksato/golang-boilerplate/core"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var INITDB *gorm.DB
var TESTDB *gorm.DB
var CTX *ModelCtx

func startTransaction() {
	TESTDB = INITDB.Begin()
	CTX.Database = TESTDB
}

func teardownDB() {
	TESTDB = TESTDB.Rollback()
}

func startDBLog() {
	TESTDB.LogMode(true)
}

func stopDBLog() {
	TESTDB.LogMode(false)
}

func InitTestDB() *gorm.DB {
	if os.Getenv("TEST_DB") == "" {
		os.Setenv("TEST_DB", "user=postgres dbname=server_test sslmode=disable")
	} else {
		os.Setenv("TEST_DB", os.Getenv("TEST_DB"))
	}

	var err error
	var db *gorm.DB
	if db, err = core.OpenTestConnection(); err != nil {
		fmt.Println("No error should happen when connecting to test database, but got", err)
	}

	if os.Getenv("TEST_DB_LOGMODE") == "true" {
		fmt.Println("Setting logmode to true")
		db.LogMode(true)
	} else {
		db.LogMode(false)
	}

	db.DB().SetMaxIdleConns(1)
	db.DB().SetMaxOpenConns(1)

	runMigration(db)

	return db
}

func runMigration(db *gorm.DB) {
	values := []interface{}{
		&Configuration{},
		&User{},
	}

	for _, value := range values {
		derr := db.DropTableIfExists(value).Error
		if derr != nil {
			panic(fmt.Sprintf("Error dropping table %+v ", derr))
		}
	}

	if err := db.AutoMigrate(values...).Error; err != nil {
		panic(fmt.Sprintf("No error should happen when create table, but got %+v", err))
	}

	db.Exec("CREATE UNIQUE INDEX idx_users_email ON users USING btree (email);")
	db.Exec("CREATE UNIQUE INDEX idx_users_username ON users USING btree (username);")
	db.Exec("CREATE UNIQUE INDEX idx_lower_case_username ON users ((lower(username)));")

	SeedDatabase(db)
}

func DeleteAllCommitedEntities(db *gorm.DB) {
	err := db.Delete(&Configuration{}).Error
	if err != nil {
		fmt.Println("Error deleting Configuration", err)
	}
	err = db.Unscoped().Delete(&User{}).Error
	if err != nil {
		fmt.Println("Error deleting User", err)
	}
}

func SeedDatabase(db *gorm.DB) {
	config := Configuration{
		Model:       Model{ID: 1},
		MinValueBuy: 25.0,
	}
	db.Create(&config)

	password := "123456"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := User{Model: Model{ID: 999}, Name: "System", Email: "system@model.com", Username: "system", HashedPassword: hashedPassword}
	db.Create(&user)
}
