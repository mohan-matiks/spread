package pkg

import (
	"bytes"
	"context"

	"github.com/SwishHQ/spread/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	AWSConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	s3Client *s3.Client
	bucket   string
}

func NewR2Service() (*S3Service, error) {
	account := config.CloudflareR2AccountID
	bucket := config.CloudflareR2Bucket
	accessKey := config.CloudflareR2AccessKeyID
	secretKey := config.CloudflareR2SecretAccessKey

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: "https://" + account + ".r2.cloudflarestorage.com",
		}, nil
	})

	cfg, err := AWSConfig.LoadDefaultConfig(context.TODO(),
		func(o *AWSConfig.LoadOptions) error {
			o.Region = "apac"
			o.EndpointResolverWithOptions = r2Resolver
			o.Credentials = credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)

	return &S3Service{
		s3Client: s3Client,
		bucket:   bucket,
	}, nil
}

func (s *S3Service) UploadFileToR2(ctx context.Context, key string, file []byte) error {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(file),
		ContentType: aws.String("application/zip"), // Set the content type to application/zip for .zip files
	}

	// Upload the file to Cloudflare R2 Storage
	_, err := s.s3Client.PutObject(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
