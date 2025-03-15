package controller

import (
	"github.com/SwishHQ/spread/logger"
	"github.com/SwishHQ/spread/src/service"
	"github.com/SwishHQ/spread/types"
	"github.com/SwishHQ/spread/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ClientController interface {
	CheckUpdate(c *fiber.Ctx) error
	ReportStatusDeploy(c *fiber.Ctx) error
	ReportStatusDownload(c *fiber.Ctx) error
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
		logger.L.Error("In CheckUpdate: Error checking update", zap.Error(err))
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

func (c *clientController) ReportStatusDeploy(ctx *fiber.Ctx) error {
	reportStatusRequest := new(types.ReportStatusDeployRequest)
	validationErrors := utils.BindAndValidate(ctx, reportStatusRequest)
	if len(validationErrors) > 0 {
		logger.L.Error("In ReportStatusDeploy: Validation errors", zap.Any("validationErrors", validationErrors))
		return ctx.Status(fiber.StatusOK).Send([]byte("OK"))
	}
	logger.L.Info("In ReportStatusDeploy", zap.Any("reportStatusRequest", reportStatusRequest))
	err := c.clientService.ReportStatusDeploy(reportStatusRequest)
	if err != nil {
		logger.L.Error("In ReportStatusDeploy: Error reporting status", zap.Error(err))
	}
	return ctx.Status(fiber.StatusOK).Send([]byte("OK"))
}

func (c *clientController) ReportStatusDownload(ctx *fiber.Ctx) error {
	reportStatusRequest := new(types.ReportStatusDownloadRequest)
	validationErrors := utils.BindAndValidate(ctx, reportStatusRequest)
	if len(validationErrors) > 0 {
		logger.L.Error("In ReportStatusDownload: Validation errors", zap.Any("validationErrors", validationErrors))
		return ctx.Status(fiber.StatusOK).Send([]byte("OK"))
	}
	logger.L.Info("In ReportStatusDownload", zap.Any("reportStatusRequest", reportStatusRequest))
	err := c.clientService.ReportStatusDownload(reportStatusRequest)
	if err != nil {
		logger.L.Error("In ReportStatusDownload: Error reporting status", zap.Error(err))
	}
	return ctx.Status(fiber.StatusOK).Send([]byte("OK"))
}
