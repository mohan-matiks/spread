package main

import (
	"log"
	"time"

	"github.com/SwishHQ/spread/config"
	"github.com/SwishHQ/spread/pkg"
	"github.com/SwishHQ/spread/src/controller"
	"github.com/SwishHQ/spread/src/repository"
	"github.com/SwishHQ/spread/src/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recover "github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Fiber instance
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

	app.Post("/app", appController.CreateApp)
	app.Post("/environment", environmentController.CreateEnvironment)

	bundleGroup := app.Group("/bundle")
	bundleGroup.Post("/create", bundleController.CreateNewBundle)
	bundleGroup.Post("/upload", bundleController.UploadBundle)
	bundleGroup.Post("/rollback", bundleController.Rollback)
	// Start server
	log.Println("Server started on port " + config.ServerPort)
	log.Fatal(app.Listen(":" + config.ServerPort))
}
