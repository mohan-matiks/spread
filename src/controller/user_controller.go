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
	GetUser(c *fiber.Ctx) error
	SetupStatus(c *fiber.Ctx) error
	InitUser(c *fiber.Ctx) error
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
	return utils.SuccessResponse(ctx, createdUser)
}

func (c *userController) GetUser(ctx *fiber.Ctx) error {
	// Get userId from the context
	// Get user from context
	user := ctx.Locals("user")
	if user == nil {
		logger.L.Error("In GetUser: User not found in context")
		return utils.ErrorResponse(ctx, "User not found")
	}
	return utils.SuccessResponse(ctx, user)
}

func (c *userController) LoginUser(ctx *fiber.Ctx) error {
	user := types.LoginUserRequest{}
	validationErrors := utils.BindAndValidate(ctx, &user)
	if len(validationErrors) > 0 {
		logger.L.Error("In LoginUser: Validation errors", zap.Any("validationErrors", validationErrors))
		return utils.ValidationErrorResponse(ctx, validationErrors)
	}
	accessToken, err := c.userService.Login(&user)
	if err != nil {
		logger.L.Error("In LoginUser: Error logging in", zap.Error(err))
		return utils.ErrorResponse(ctx, err.Error())
	}
	return utils.SuccessResponse(ctx, fiber.Map{
		"access_token": *accessToken,
	})
}

func (c *userController) SetupStatus(ctx *fiber.Ctx) error {
	count, err := c.userService.Count(ctx.Context())
	if err != nil {
		logger.L.Error("In SetupStatus: Error getting user count", zap.Error(err))
		return utils.ErrorResponse(ctx, err.Error())
	}
	completed := count > 0
	return utils.SuccessResponse(ctx, fiber.Map{
		"completed": completed,
	})
}

func (c *userController) InitUser(ctx *fiber.Ctx) error {
	// Check if any users exist
	count, err := c.userService.Count(ctx.Context())
	if err != nil {
		logger.L.Error("In InitUser: Error getting user count", zap.Error(err))
		return utils.ErrorResponse(ctx, err.Error())
	}

	if count > 0 {
		logger.L.Error("In InitUser: Users already exist, cannot initialize")
		return utils.ErrorResponse(ctx, "Users already exist. Cannot initialize first user.")
	}

	// Parse and validate the request
	user := types.CreateUserRequest{}
	validationErrors := utils.BindAndValidate(ctx, &user)
	if len(validationErrors) > 0 {
		logger.L.Error("In InitUser: Validation errors", zap.Any("validationErrors", validationErrors))
		return utils.ValidationErrorResponse(ctx, validationErrors)
	}

	// Force the first user to be an admin
	user.Roles = []string{"admin"}

	createdUser, err := c.userService.Create(&user)
	if err != nil {
		logger.L.Error("In InitUser: Error creating user", zap.Error(err))
		return utils.ErrorResponse(ctx, err.Error())
	}

	logger.L.Info("In InitUser: First user created successfully", zap.String("username", createdUser.Username))
	return utils.SuccessResponse(ctx, createdUser)
}
