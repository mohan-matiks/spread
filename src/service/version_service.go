package service

import (
	"context"
	"errors"
	"sort"

	"github.com/SwishHQ/spread/src/model"
	"github.com/SwishHQ/spread/src/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type VersionService interface {
	CreateVersion(ctx context.Context, version *model.Version) (*model.Version, error)
	UpdateVersionCurrentBundleIdByVersionId(ctx context.Context, id primitive.ObjectID, currentBundleId primitive.ObjectID) (*model.Version, error)
	GetVersionByEnvironmentIdAndAppVersion(ctx context.Context, environmentId primitive.ObjectID, appVersion string) (*model.Version, error)
	GetVersionByEnvironmentIdAndVersionId(ctx context.Context, versionId primitive.ObjectID, environmentId primitive.ObjectID) (*model.Version, error)
	GetLatestVersionByEnvironmentId(ctx context.Context, environmentId primitive.ObjectID) (*model.Version, error)
	GetAllVersionsByEnvironmentId(ctx context.Context, environmentId primitive.ObjectID) ([]*model.Version, error)
	GetByVersionId(ctx context.Context, versionId primitive.ObjectID) (*model.Version, error)
}

type versionService struct {
	versionRepository repository.VersionRepository
}

func NewVersionService(versionRepository repository.VersionRepository) VersionService {
	return &versionService{versionRepository: versionRepository}
}

func (v *versionService) GetVersionByEnvironmentIdAndAppVersion(ctx context.Context, environmentId primitive.ObjectID, appVersion string) (*model.Version, error) {
	version, err := v.versionRepository.GetByEnvironmentIdAndAppVersion(ctx, environmentId, appVersion)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return version, nil
}

func (v *versionService) CreateVersion(ctx context.Context, version *model.Version) (*model.Version, error) {
	version, err := v.versionRepository.Create(ctx, version)
	if err != nil {
		return nil, err
	}
	return version, nil
}

func (v *versionService) UpdateVersionCurrentBundleIdByVersionId(ctx context.Context, id primitive.ObjectID, currentBundleId primitive.ObjectID) (*model.Version, error) {
	version, err := v.versionRepository.UpdateCurrentBundleId(ctx, id, currentBundleId)
	if err != nil {
		return nil, err
	}
	return version, nil
}

func (v *versionService) GetLatestVersionByEnvironmentId(ctx context.Context, environmentId primitive.ObjectID) (*model.Version, error) {
	version, err := v.versionRepository.GetLatestVersionByEnvironmentId(ctx, environmentId)
	if err != nil {
		return nil, err
	}
	return version, nil
}

func (v *versionService) GetAllVersionsByEnvironmentId(ctx context.Context, environmentId primitive.ObjectID) ([]*model.Version, error) {
	versions, err := v.versionRepository.GetAllByEnvironmentId(ctx, environmentId)
	if err != nil {
		return nil, err
	}
	// sort versions by createdAt in descending order
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].CreatedAt.After(versions[j].CreatedAt)
	})
	return versions, nil
}

func (v *versionService) GetByVersionId(ctx context.Context, versionId primitive.ObjectID) (*model.Version, error) {
	version, err := v.versionRepository.GetById(ctx, versionId)
	if err != nil {
		return nil, err
	}
	if version == nil {
		return nil, errors.New("version not found")
	}
	return version, nil
}

func (v *versionService) GetVersionByEnvironmentIdAndVersionId(ctx context.Context, versionId primitive.ObjectID, environmentId primitive.ObjectID) (*model.Version, error) {
	version, err := v.versionRepository.GetByIdAndEnvironmentId(ctx, versionId, environmentId)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if version == nil {
		return nil, errors.New("version not found")
	}
	return version, nil
}
