package service

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/SwishHQ/spread/pkg"
	"github.com/SwishHQ/spread/src/repository"
)

type BundleService interface {
	UploadBundle(fileName string, file *multipart.FileHeader) error
}

type bundleServiceImpl struct {
	bundleRepository repository.BundleRepository
}

func NewBundleService(bundleRepository repository.BundleRepository) BundleService {
	return &bundleServiceImpl{bundleRepository: bundleRepository}
}

func (bundleService *bundleServiceImpl) UploadBundle(fileName string, file *multipart.FileHeader) error {
	r2Service, err := pkg.NewR2Service()
	if err != nil {
		return err
	}

	fileBytes, err := file.Open()
	if err != nil {
		return err
	}
	defer fileBytes.Close()

	// Read file into byte slice
	buffer, err := io.ReadAll(fileBytes)
	if err != nil {
		return err
	}

	err = r2Service.UploadFileToR2(context.Background(), fileName, buffer)
	if err != nil {
		return err
	}
	return nil
}
