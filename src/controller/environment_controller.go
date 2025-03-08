package controller

import (
	"github.com/SwishHQ/spread/logger"
	"github.com/SwishHQ/spread/src/service"
	"github.com/SwishHQ/spread/types"
	"github.com/SwishHQ/spread/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type EnvironmentController interface {
	CreateEnvironment(c *fiber.Ctx) error
}

type environmentControllerImpl struct {
	environmentService service.EnvironmentService
}

func NewEnvironmentController(environmentService service.EnvironmentService) EnvironmentController {
	return &environmentControllerImpl{environmentService: environmentService}
}

func (environmentController *environmentControllerImpl) CreateEnvironment(c *fiber.Ctx) error {
	var environmentRequest types.CreateEnvironmentRequest
	validationErrors := utils.BindAndValidate(c, &environmentRequest)
	if len(validationErrors) > 0 {
		return utils.ValidationErrorResponse(c, validationErrors)
	}
	logger.L.Info("Creating environment", zap.Any("environmentRequest", environmentRequest))
	environment, err := environmentController.environmentService.CreateEnvironment(&environmentRequest)
	if err != nil {
		logger.L.Error("Error creating environment", zap.Error(err))
		return utils.ErrorResponse(c, err.Error())
	}
	logger.L.Info("Environment created", zap.Any("environment", environment))
	return utils.SuccessResponse(c, fiber.Map{
		"key":  environment.Key,
		"name": environment.Name,
	})
}
