package server

import (
	"net/http"
	"os"

	"github.com/brunoksato/golang-boilerplate/api"
	middle "github.com/brunoksato/golang-boilerplate/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.elastic.co/apm/module/apmechov4"
)

var cors []string

type MiddlewareConfigurer interface {
	ConfigureDefaultApiMiddleware(*echo.Echo) *echo.Echo
	ConfigurePublicApiMiddleware(*echo.Echo) *echo.Group
	ConfigurePrivateApiMiddleware(*echo.Echo) *echo.Group
	ConfigureCronJobApiMiddleware(*echo.Echo) *echo.Group
	ConfigureAdminApiMiddleware(*echo.Echo) *echo.Group
}

func SetupRouter(mc MiddlewareConfigurer) *echo.Echo {
	root := echo.New()
	root.Use(apmechov4.Middleware())
	root.GET("/", func(c echo.Context) error {
		hello := map[string]interface{}{"status": "API OK"}
		return c.JSON(http.StatusOK, hello)
	})

	root.POST("/webhook/sample", api.WebhookSample)

	if os.Getenv("ENV") == "production" {
		cors = []string{"*"}
	} else {
		cors = []string{"*"}
	}

	//
	// PUBLIC ENDPOINTS
	//
	public := mc.ConfigurePublicApiMiddleware(root)

	public.POST("/signin", api.SignIn)
	public.POST("/signup", api.SignUp)
	public.GET("/recover/:email", api.RecoverPassword)
	public.PUT("/change_password", api.ChangePasswordExternal)

	//
	// CRONJOB ENDPOINTS
	//
	cronjob := mc.ConfigureCronJobApiMiddleware(root)
	cronjob.GET("/sample", api.CronJobSample)
	//
	// PRIVATE ENDPOINTS
	//
	private := mc.ConfigurePrivateApiMiddleware(root)

	/* General */
	private.GET("/logout", api.Logout)

	/* User */
	private.GET("/users/me", api.Me)

	private.PUT("/users", api.UpdateUser)
	private.PUT("/users/password", api.ChangePassword)

	//
	// ADMIN ENDPOINTS
	//
	admin := mc.ConfigureAdminApiMiddleware(root)

	admin.GET("/users", api.List)

	return root
}

type ProductionMiddlewareConfigurer struct{}

func (mc ProductionMiddlewareConfigurer) ConfigureDefaultApiMiddleware(root *echo.Echo) *echo.Echo {
	root.Use(middleware.Logger())
	root.Use(middleware.Recover())
	root.Use(middleware.CORS())
	root.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		ContentSecurityPolicy: "default-src 'self'",
	}))
	root.Use(middle.DBMiddleware(RW_DB_POOL))
	root.Use(middle.ElasticMiddleware(ES))
	root.Use(middle.DetermineType)
	root.Use(middle.InitializePayload)
	root.Use(middle.LoadConfigurations)

	return root
}

func (mc ProductionMiddlewareConfigurer) ConfigurePublicApiMiddleware(root *echo.Echo) *echo.Group {
	api := mc.ConfigureDefaultApiMiddleware(root)
	public := api.Group("/public")
	public.Use(middleware.CORS())
	public.Use(middle.SettingHeaders)
	public.Use(middle.Session)

	return public
}

func (mc ProductionMiddlewareConfigurer) ConfigurePrivateApiMiddleware(root *echo.Echo) *echo.Group {
	api := mc.ConfigureDefaultApiMiddleware(root)
	private := api.Group("/api")
	private.Use(middleware.Gzip())
	private.Use(middleware.CORS())
	private.Use(middle.SettingHeaders)
	private.Use(middle.Session)

	return private
}

func (mc ProductionMiddlewareConfigurer) ConfigureCronJobApiMiddleware(root *echo.Echo) *echo.Group {
	api := mc.ConfigureDefaultApiMiddleware(root)
	private := api.Group("/cronjob")
	private.Use(middleware.CORS())
	private.Use(middleware.Gzip())
	private.Use(middle.SettingHeaders)
	private.Use(middle.Session)

	return private
}

func (mc ProductionMiddlewareConfigurer) ConfigureAdminApiMiddleware(root *echo.Echo) *echo.Group {
	api := mc.ConfigureDefaultApiMiddleware(root)
	private := api.Group("/admin")
	private.Use(middleware.CORS())
	private.Use(middleware.Gzip())
	private.Use(middle.SettingHeaders)
	private.Use(middle.Session)

	return private
}
