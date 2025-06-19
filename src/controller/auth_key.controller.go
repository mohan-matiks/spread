package controller

import (
	"github.com/SwishHQ/spread/src/model"
	"github.com/SwishHQ/spread/src/service"
	"github.com/SwishHQ/spread/types"
	"github.com/SwishHQ/spread/utils"
	"github.com/gofiber/fiber/v2"
)

type AuthKeyController interface {
	CreateAuthKey(c *fiber.Ctx) error
	GetAllAuthKeys(c *fiber.Ctx) error
}

type authKeyController struct {
	authKeyService service.AuthKeyService
}

func NewAuthKeyController(authKeyService service.AuthKeyService) AuthKeyController {
	return &authKeyController{authKeyService: authKeyService}
}

func (c *authKeyController) CreateAuthKey(ctx *fiber.Ctx) error {
	authKey := types.CreateAuthKeyRequest{}
	validationErrors := utils.BindAndValidate(ctx, &authKey)
	if len(validationErrors) > 0 {
		return utils.ValidationErrorResponse(ctx, validationErrors)
	}
	user := ctx.Locals("user").(*model.User)
	createdBy := user.Username
	createdAuthKey, err := c.authKeyService.CreateAuthKey(authKey.Name, createdBy)
	if err != nil {
		return utils.ErrorResponse(ctx, err.Error())
	}
	return utils.SuccessResponse(ctx, createdAuthKey)
}

func (c *authKeyController) GetAllAuthKeys(ctx *fiber.Ctx) error {
	authKeys, err := c.authKeyService.GetAllAuthKeys(ctx.Context())
	if err != nil {
		return utils.ErrorResponse(ctx, err.Error())
	}
	return utils.SuccessResponse(ctx, authKeys)
}
