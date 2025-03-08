package controller

import (
	"github.com/SwishHQ/spread/logger"
	"github.com/SwishHQ/spread/src/service"
	"github.com/SwishHQ/spread/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type BundleController interface {
	UploadBundle(c *fiber.Ctx) error
}

type bundleControllerImpl struct {
	bundleService service.BundleService
}

func NewBundleController(bundleService service.BundleService) BundleController {
	return &bundleControllerImpl{bundleService: bundleService}
}

func (bundleController *bundleControllerImpl) UploadBundle(c *fiber.Ctx) error {
	logger.L.Info("In UploadBundle: Uploading bundle", zap.Any("filename", c.FormValue("fileName")))
	fileName := c.FormValue("filename")
	if fileName == "" {
		return utils.ErrorResponse(c, "File name is required")
	}

	uploadedFile, err := c.FormFile("file")
	if err != nil {
		return utils.ErrorResponse(c, "Failed to get uploaded file")
	}

	err = bundleController.bundleService.UploadBundle(fileName, uploadedFile)
	if err != nil {
		logger.L.Error("In UploadBundle: Failed to upload bundle", zap.Error(err))
		return utils.ErrorResponse(c, "Failed to upload bundle")
	}
	logger.L.Info("In UploadBundle: Bundle uploaded successfully", zap.Any("fileName", fileName))
	return utils.SuccessResponse(c, nil)
}
