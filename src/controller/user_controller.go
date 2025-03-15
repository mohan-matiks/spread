package controller

import (
	"github.com/SwishHQ/spread/logger"
	"github.com/SwishHQ/spread/src/service"
	"github.com/SwishHQ/spread/types"
	"github.com/SwishHQ/spread/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserController interface {
	CreateUser(c *fiber.Ctx) error
	LoginUser(c *fiber.Ctx) error
}

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &userController{userService: userService}
}

func (c *userController) CreateUser(ctx *fiber.Ctx) error {
	user := types.CreateUserRequest{}
	validationErrors := utils.BindAndValidate(ctx, &user)
	if len(validationErrors) > 0 {
		logger.L.Error("In CreateUser: Validation errors", zap.Any("validationErrors", validationErrors))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": validationErrors,
		})
	}
	createdUser, err := c.userService.Create(&user)
	if err != nil {
		logger.L.Error("In CreateUser: Error creating user", zap.Error(err))
		return utils.ErrorResponse(ctx, err.Error())
	}
	return ctx.Status(fiber.StatusCreated).JSON(createdUser)
}

func (c *userController) LoginUser(ctx *fiber.Ctx) error {
	user := types.LoginUserRequest{}
	validationErrors := utils.BindAndValidate(ctx, &user)
	if len(validationErrors) > 0 {
		logger.L.Error("In LoginUser: Validation errors", zap.Any("validationErrors", validationErrors))
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": validationErrors,
		})
	}
	accessToken, err := c.userService.Login(&user)
	if err != nil {
		logger.L.Error("In LoginUser: Error logging in", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": *accessToken,
	})
}
