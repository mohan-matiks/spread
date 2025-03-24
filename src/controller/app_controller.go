package controller

import (
	"context"

	"github.com/SwishHQ/spread/logger"
	"github.com/SwishHQ/spread/src/service"
	"github.com/SwishHQ/spread/types"
	"github.com/SwishHQ/spread/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AppController interface {
	CreateApp(c *fiber.Ctx) error
	GetApps(c *fiber.Ctx) error
	GetAppById(c *fiber.Ctx) error
}

type appControllerImpl struct {
	appService service.AppService
}

func NewAppController(appService service.AppService) AppController {
	return &appControllerImpl{appService: appService}
}

func (controller *appControllerImpl) CreateApp(c *fiber.Ctx) error {
	var appRequest types.CreateAppRequest
	validationErrors := utils.BindAndValidate(c, &appRequest)
	if len(validationErrors) > 0 {
		return utils.ValidationErrorResponse(c, validationErrors)
	}
	logger.L.Info("Creating app", zap.Any("appRequest", appRequest))
	app, err := controller.appService.CreateApp(context.Background(), appRequest.AppName, appRequest.OS)
	if err != nil {
		logger.L.Error("Error creating app", zap.Error(err))
		return utils.ErrorResponse(c, "Error creating app")
	}
	logger.L.Info("App created", zap.Any("app", app))
	return utils.SuccessResponse(c, fiber.Map{
		"name": app.Name,
		"os":   app.OS,
	})
}

func (controller *appControllerImpl) GetApps(c *fiber.Ctx) error {
	apps, err := controller.appService.GetApps(context.Background())
	if err != nil {
		return utils.ErrorResponse(c, "Error getting apps")
	}
	return utils.SuccessResponse(c, apps)
}

func (controller *appControllerImpl) GetAppById(c *fiber.Ctx) error {
	id := c.Params("id")
	app, err := controller.appService.GetAppById(context.Background(), id)
	if err != nil {
		logger.L.Error("In GetAppById: Error getting app", zap.Error(err))
		return utils.ErrorResponse(c, err.Error())
	}
	return utils.SuccessResponse(c, app)
}
