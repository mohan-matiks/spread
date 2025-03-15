package cmd

import (
	"log"
	"time"

	"github.com/SwishHQ/spread/config"
	"github.com/SwishHQ/spread/middleware"
	"github.com/SwishHQ/spread/pkg"
	"github.com/SwishHQ/spread/src/controller"
	"github.com/SwishHQ/spread/src/repository"
	"github.com/SwishHQ/spread/src/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs Spread server",
	Run:   serve,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	db, errMongoConnection := pkg.MongoConnection()
	if errMongoConnection != nil {
		log.Fatal(errMongoConnection)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service":   "spread",
			"status":    "running",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	authKeyRepository := repository.NewAuthKeyRepository(db)
	authKeyService := service.NewAuthKeyService(authKeyRepository)
	authKeyController := controller.NewAuthKeyController(authKeyService)

	appRepository := repository.NewAppRepository(db)
	appService := service.NewAppService(appRepository)
	appController := controller.NewAppController(appService)

	environmentRepository := repository.NewEnvironmentRepository(db)
	environmentService := service.NewEnvironmentService(appService, environmentRepository)
	environmentController := controller.NewEnvironmentController(environmentService)

	versionRepository := repository.NewVersionRepository(db)
	versionService := service.NewVersionService(versionRepository)

	bundleRepository := repository.NewBundleRepository(db)
	bundleService := service.NewBundleService(appService, versionService, environmentService, bundleRepository)
	bundleController := controller.NewBundleController(bundleService)

	clientService := service.NewClientService(appService, environmentService, bundleService, versionService)
	clientController := controller.NewClientController(clientService)

	// public endpoints
	app.Post("/login", userController.LoginUser)
	app.Get("/apps", appController.GetApps)
	app.Post("/user/create", userController.CreateUser)

	// code-push compatible endpoints
	app.Get("/v0.1/public/codepush/update_check", clientController.CheckUpdate)
	app.Post("/v0.1/public/codepush/report_status/deploy", clientController.ReportStatusDeploy)
	app.Post("/v0.1/public/codepush/report_status/download", clientController.ReportStatusDownload)

	// protected endpoints
	coreGroup := app.Group("/core", func(c *fiber.Ctx) error {
		return middleware.AuthMiddleware(c, userService)
	})
	coreGroup.Post("/environment", environmentController.CreateEnvironment)
	coreGroup.Post("/app", appController.CreateApp)
	coreGroup.Post("/auth-key/create", authKeyController.CreateAuthKey)

	// auth key protected endpoints
	bundleGroup := app.Group("/bundle", func(c *fiber.Ctx) error {
		return middleware.AuthKeyMiddleware(c, authKeyService)
	})
	bundleGroup.Post("/create", bundleController.CreateNewBundle)
	bundleGroup.Post("/upload", bundleController.UploadBundle)
	bundleGroup.Post("/rollback", bundleController.Rollback)

	// Start server
	log.Println("Server started on port " + config.ServerPort)
	log.Fatal(app.Listen(":" + config.ServerPort))
}
