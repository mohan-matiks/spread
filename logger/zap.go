package logger

import (
	"os"

	"github.com/SwishHQ/spread/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var L *zap.Logger

func init() {
	L = createLogger()
}

func createLogger() *zap.Logger {
	// for local development
	if config.ENV == "local" {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.Encoding = "console"
		logger, _ := config.Build()
		defer logger.Sync()
		return logger
	}

	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "/var/log/layout/application.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	})

	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	var core zapcore.Core

	if config.ENV == "prod" {
		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "timestamp"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

		config := zap.Config{
			Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
			Development:       false,
			DisableCaller:     false,
			DisableStacktrace: false,
			Sampling:          nil,
			Encoding:          "json",
			EncoderConfig:     encoderCfg,
			OutputPaths: []string{
				"stderr",
			},
			ErrorOutputPaths: []string{
				"stderr",
			},
			InitialFields: map[string]interface{}{
				"pid": os.Getpid(),
			},
		}
		return zap.Must(config.Build())
	}

	core = zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)
	return zap.New(core)
}
