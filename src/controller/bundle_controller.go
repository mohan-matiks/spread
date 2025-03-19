package controller

import (
	"github.com/SwishHQ/spread/logger"
	"github.com/SwishHQ/spread/src/model"
	"github.com/SwishHQ/spread/src/service"
	"github.com/SwishHQ/spread/types"
	"github.com/SwishHQ/spread/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type BundleController interface {
	UploadBundle(c *fiber.Ctx) error
	CreateNewBundle(c *fiber.Ctx) error
	Rollback(c *fiber.Ctx) error
}

type bundleControllerImpl struct {
	bundleService service.BundleService
}

func NewBundleController(bundleService service.BundleService) BundleController {
	return &bundleControllerImpl{bundleService: bundleService}
}

func (bundleController *bundleControllerImpl) UploadBundle(c *fiber.Ctx) error {
	authKey := c.Locals("authKey").(*model.AuthKey)
	keyUser := authKey.CreatedBy
	logger.L.Info("In UploadBundle: Uploading bundle", zap.Any("filename", c.FormValue("fileName")), zap.Any("keyUser", keyUser))
	fileName := c.FormValue("filename")
	if fileName == "" {
		logger.L.Error("In UploadBundle: File name is required")
		return utils.ErrorResponse(c, "File name is required")
	}

	uploadedFile, err := c.FormFile("file")
	if err != nil {
		logger.L.Error("In UploadBundle: no file found", zap.Error(err))
		return utils.ErrorResponse(c, "No file found")
	}

	err = bundleController.bundleService.UploadBundle(fileName, uploadedFile)
	if err != nil {
		logger.L.Error("In UploadBundle: Failed to upload bundle", zap.Error(err))
		return utils.ErrorResponse(c, "Failed to upload bundle")
	}
	logger.L.Info("In UploadBundle: Bundle uploaded successfully", zap.Any("fileName", fileName))
	return utils.SuccessResponse(c, nil)
}

func (bundleController *bundleControllerImpl) CreateNewBundle(c *fiber.Ctx) error {
	var createNewBundleRequest types.CreateNewBundleRequest
	validationErrors := utils.BindAndValidate(c, &createNewBundleRequest)
	if len(validationErrors) > 0 {
		logger.L.Error("In CreateNewBundle: Validation errors", zap.Any("validationErrors", validationErrors))
		return utils.ValidationErrorResponse(c, validationErrors)
	}
	authKey := c.Locals("authKey").(*model.AuthKey)
	createdBy := authKey.CreatedBy
	bundle, err := bundleController.bundleService.CreateNewBundle(&createNewBundleRequest, createdBy)
	logger.L.Info("In CreateNewBundle: Creating new bundle", zap.Any("keyUser", createdBy), zap.Any("createNewBundleRequest", createNewBundleRequest))
	if err != nil {
		logger.L.Error("In CreateNewBundle: Failed to create new bundle", zap.Error(err))
		return utils.ErrorResponse(c, err.Error())
	}
	logger.L.Info("In CreateNewBundle: New bundle created successfully", zap.Any("bundle", bundle))
	return utils.SuccessResponse(c, bundle)
}

func (bundleController *bundleControllerImpl) Rollback(c *fiber.Ctx) error {
	var rollbackRequest types.RollbackRequest
	validationErrors := utils.BindAndValidate(c, &rollbackRequest)
	if len(validationErrors) > 0 {
		logger.L.Error("In Rollback: Validation errors", zap.Any("validationErrors", validationErrors))
		return utils.ValidationErrorResponse(c, validationErrors)
	}
	// log key user as audit log
	authKey := c.Locals("authKey").(*model.AuthKey)
	keyUser := authKey.CreatedBy
	logger.L.Info("In Rollback", zap.Any("keyUser", keyUser), zap.Any("rollbackRequest", rollbackRequest))

	rollbackBundle, err := bundleController.bundleService.Rollback(&rollbackRequest)
	if err != nil {
		logger.L.Error("In Rollback: Failed to rollback", zap.Error(err))
		return utils.ErrorResponse(c, err.Error())
	}
	if rollbackBundle == nil {
		return utils.SuccessResponse(c, fiber.Map{
			"success": true,
			"message": "Rollback successful",
		})
	}
	logger.L.Info("In Rollback: Rollback bundle found", zap.Any("rollbackBundle", rollbackBundle))
	return utils.SuccessResponse(c, rollbackBundle)
}
