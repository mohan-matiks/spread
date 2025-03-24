package controller

import (
	"github.com/SwishHQ/spread/src/service"
	"github.com/SwishHQ/spread/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VersionController interface {
	GetAll(c *fiber.Ctx) error
	GetByVersionId(c *fiber.Ctx) error
}

type versionController struct {
	versionService service.VersionService
}

func NewVersionController(versionService service.VersionService) VersionController {
	return &versionController{versionService: versionService}
}

func (v *versionController) GetAll(c *fiber.Ctx) error {
	environmentId := c.Query("environmentId")
	environmentIdObjectID, err := primitive.ObjectIDFromHex(environmentId)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}
	versions, err := v.versionService.GetAllVersionsByEnvironmentId(c.Context(), environmentIdObjectID)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}
	return utils.SuccessResponse(c, versions)
}

func (v *versionController) GetByVersionId(c *fiber.Ctx) error {
	versionId := c.Params("versionId")
	versionIdObjectID, err := primitive.ObjectIDFromHex(versionId)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}
	version, err := v.versionService.GetByVersionId(c.Context(), versionIdObjectID)
	if err != nil {
		return utils.ErrorResponse(c, err.Error())
	}
	return utils.SuccessResponse(c, version)
}
