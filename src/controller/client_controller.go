package controller

import (
	"github.com/SwishHQ/spread/logger"
	"github.com/SwishHQ/spread/src/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ClientController interface {
	CheckUpdate(c *fiber.Ctx) error
}

type clientController struct {
	clientService service.ClientService
}

func NewClientController(clientService service.ClientService) ClientController {
	return &clientController{
		clientService: clientService,
	}
}
func (c *clientController) CheckUpdate(ctx *fiber.Ctx) error {
	environmentKey := ctx.Query("deployment_key")
	appVersion := ctx.Query("app_version")
	bundleHash := ctx.Query("package_hash")
	label := ctx.Query("label")
	clientUniqueId := ctx.Query("client_unique_id")

	logger.L.Info("In CheckUpdate", zap.String("environmentKey", environmentKey), zap.String("appVersion", appVersion), zap.String("bundleHash", bundleHash), zap.String("label", label), zap.String("clientUniqueId", clientUniqueId))
	updateInfo, err := c.clientService.CheckUpdate(environmentKey, appVersion, bundleHash)
	if err != nil {
		// when update_info is nil, it means there is no update available
		// to be safe of sdk crashing upon any irrelavent error, we return nil
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"update_info": nil,
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"update_info": updateInfo,
	})
}
