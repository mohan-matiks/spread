package config

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	ENV                         = GetEnv("ENV", "dev")
	ServerPort                  = GetEnv("PORT", "3000")
	AppName                     = GetEnv("APP_NAME", "")
	MongoUrl                    = GetEnv("MONGODB_URL", "")
	TokenSecret                 = GetEnv("TOKEN_SECRET", "")
	MongoDatabase               = GetEnv("MONGODB_DATABASE", "")
	CloudflareR2AccountID       = GetEnv("CLOUDFLARE_R2_ACCOUNT_ID", "")
	CloudflareR2Bucket          = GetEnv("CLOUDFLARE_R2_BUCKET", "")
	CloudflareR2AccessKeyID     = GetEnv("CLOUDFLARE_R2_ACCESS_KEY_ID", "")
	CloudflareR2SecretAccessKey = GetEnv("CLOUDFLARE_R2_SECRET_ACCESS_KEY", "")
)

func GetEnv(key, defaultValue string) string {

	if _, exists := os.LookupEnv(key); !exists {
		err := godotenv.Load(".env")
		if err != nil {
			zap.Error(err)
		}
		value := os.Getenv(key)
		if len(value) == 0 {
			return defaultValue
		}
		return value
	}

	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
