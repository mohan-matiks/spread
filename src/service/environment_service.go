package service

import (
	"context"
	"errors"
	"time"

	"github.com/SwishHQ/spread/src/model"
	"github.com/SwishHQ/spread/src/repository"
	"github.com/SwishHQ/spread/types"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EnvironmentService interface {
	CreateEnvironment(environmentRequest *types.CreateEnvironmentRequest) (*model.Environment, error)
	GetEnvironmentByAppIdAndName(ctx context.Context, appId primitive.ObjectID, environmentName string) (*model.Environment, error)
	GetEnvironmentByKey(ctx context.Context, key string) (*model.Environment, error)
}

type environmentServiceImpl struct {
	appService            AppService
	environmentRepository repository.EnvironmentRepository
}

func NewEnvironmentService(appService AppService, environmentRepository repository.EnvironmentRepository) EnvironmentService {
	return &environmentServiceImpl{appService: appService, environmentRepository: environmentRepository}
}

func (environmentService *environmentServiceImpl) GetEnvironmentByAppIdAndName(ctx context.Context, appId primitive.ObjectID, environmentName string) (*model.Environment, error) {
	environment, err := environmentService.environmentRepository.GetByAppIdAndName(ctx, appId, environmentName)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return environment, nil
}

// Deployments are environments, every app (app-ios, app-android) can have multiple deployments (dev, staging, prod)
func (environmentService *environmentServiceImpl) CreateEnvironment(environmentRequest *types.CreateEnvironmentRequest) (*model.Environment, error) {
	app, err := environmentService.appService.GetAppByName(context.Background(), environmentRequest.AppName)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("app not found")
		}
		return nil, err
	}
	existingEnvironment, err := environmentService.environmentRepository.GetByAppIdAndName(context.Background(), app.Id, environmentRequest.EnvironmentName)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if existingEnvironment != nil {
		return nil, errors.New("environment with name " + environmentRequest.EnvironmentName + " already exists for app " + app.Name)
	}
	key := uuid.New().String()
	environment := model.Environment{
		AppId:     app.Id,
		Name:      environmentRequest.EnvironmentName,
		Key:       key,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
	createdEnvironment, err := environmentService.environmentRepository.Insert(context.Background(), &environment)
	if err != nil {
		return nil, err
	}
	return createdEnvironment, nil
}

func (environmentService *environmentServiceImpl) GetEnvironmentByKey(ctx context.Context, key string) (*model.Environment, error) {
	environment, err := environmentService.environmentRepository.GetByKey(ctx, key)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return environment, nil
}
