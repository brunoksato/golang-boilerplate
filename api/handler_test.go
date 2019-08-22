package api_test

import (
	"fmt"
	"os"

	"github.com/brunoksato/golang-boilerplate/api"
	middle "github.com/brunoksato/golang-boilerplate/middleware"
	"github.com/brunoksato/golang-boilerplate/model"
	"github.com/brunoksato/golang-boilerplate/server"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
)

var CTX *api.Context
var INITDB *gorm.DB
var TESTDB *gorm.DB
var ROLE string = "user"

func init() {
	os.Setenv("TEST_ON", "true")
	INITDB = model.InitTestDB()
}

func setup() {
	TESTDB = INITDB.Begin()
	model.DeleteAllCommitedEntities(TESTDB)
	model.SeedDatabase(TESTDB)
}

func setAdminRole() {
	ROLE = "admin"
}

func startTransaction() {
	TESTDB = INITDB.Begin()
	CTX.Database = TESTDB
}

func teardown() {
	TESTDB = TESTDB.Rollback()
}

func startDBLog() {
	TESTDB.LogMode(true)
}

func stopDBLog() {
	TESTDB.LogMode(false)
}

func router() *echo.Echo {
	return server.SetupRouter(TestMiddlewareConfigurer{})
}

func createTestUser() {
	password := "123456"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := model.User{
		Model:          model.Model{ID: 999},
		Email:          "system@model.com",
		Name:           "system",
		Username:       "system",
		Phone:          "12982573000",
		HashedPassword: hashedPassword,
	}
	err := TESTDB.Create(&user).Error
	if err != nil {
		panic(fmt.Sprintf("Error creating test user: %s", err))
	}
}

type TestMiddlewareConfigurer struct{}

func (mc TestMiddlewareConfigurer) ConfigureDefaultApiMiddleware(root *echo.Echo) *echo.Echo {
	root.Use(middleware.Recover())
	root.Use(middleware.CORS())
	root.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		ContentSecurityPolicy: "default-src 'self'",
	}))
	root.Use(middle.DBMiddleware(TESTDB))
	root.Use(middle.ElasticMiddleware(nil))
	root.Use(middle.DetermineType)
	root.Use(middle.InitializePayload)
	root.Use(middle.LoadConfigurations)

	return root
}

func (mc TestMiddlewareConfigurer) ConfigurePublicApiMiddleware(root *echo.Echo) *echo.Group {
	api := mc.ConfigureDefaultApiMiddleware(root)
	public := api.Group("/public")
	public.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}}))
	public.Use(middle.Session)

	return public
}

func (mc TestMiddlewareConfigurer) ConfigurePrivateApiMiddleware(root *echo.Echo) *echo.Group {
	api := mc.ConfigureDefaultApiMiddleware(root)
	private := api.Group("/api")
	private.Use(middleware.Gzip())
	private.Use(middle.SettingHeaders)
	private.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}}))
	private.Use(SetUpTestUser(ROLE))
	private.Use(SetTestUserToken)
	private.Use(middle.Session)

	return private
}

func (mc TestMiddlewareConfigurer) ConfigureCronJobApiMiddleware(root *echo.Echo) *echo.Group {
	api := mc.ConfigureDefaultApiMiddleware(root)
	private := api.Group("/cronjob")
	private.Use(middleware.CORS())
	private.Use(middleware.Gzip())
	private.Use(middle.SettingHeaders)
	private.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}}))
	private.Use(SetUpTestUser(ROLE))
	private.Use(SetTestUserToken)
	private.Use(middle.Session)

	return private
}

func (mc TestMiddlewareConfigurer) ConfigureAdminApiMiddleware(root *echo.Echo) *echo.Group {
	api := mc.ConfigureDefaultApiMiddleware(root)
	private := api.Group("/admin")
	private.Use(middleware.CORS())
	private.Use(middleware.Gzip())
	private.Use(middle.SettingHeaders)
	private.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}}))
	private.Use(SetUpTestUser(ROLE))
	private.Use(SetTestUserToken)
	private.Use(middle.Session)

	return private
}

func SetUpTestUser(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			password := "password"
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

			user := model.User{}
			dbFind := TESTDB.Unscoped().First(&user, 999)
			if dbFind.RecordNotFound() {
				user := model.User{
					Model:          model.Model{ID: 999},
					Email:          "system@model.com",
					Name:           "system",
					Username:       "system",
					Phone:          "12982573000",
					HashedPassword: hashedPassword,
				}
				err := TESTDB.Create(&user).Error
				if err != nil {
					panic(fmt.Sprintf("Error creating test user: %s", err))
				}
			} else {
				user.HashedPassword = hashedPassword
				TESTDB.Unscoped().Save(&user)
			}

			c.Set("User", user)

			return next(c)
		}
	}
}

func SetTestUserToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("User").(model.User)
		expireAt := model.JWTTokenExpirationDate()
		jwt, _ := model.IssueJWToken(user.ID, []string{"user"}, expireAt)
		c.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))
		return next(c)
	}
}
