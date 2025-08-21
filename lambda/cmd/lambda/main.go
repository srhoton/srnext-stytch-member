package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/steverhoton/srnext-stytch-member/lambda/internal/config"
	"github.com/steverhoton/srnext-stytch-member/lambda/internal/handler"
	"github.com/steverhoton/srnext-stytch-member/lambda/internal/models"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Initialize logger
	logger, err := initLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		_ = logger.Sync()
	}()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Create handler
	h, err := handler.NewHandler(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to create handler", zap.Error(err))
	}

	// Start Lambda
	logger.Info("Starting Lambda function")
	lambda.Start(func(ctx context.Context, req models.ALBRequest) (models.ALBResponse, error) {
		return h.HandleRequest(ctx, req)
	})
}

// initLogger initializes the logger
func initLogger() (*zap.Logger, error) {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	var level zapcore.Level
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	// Customize encoder config for better Lambda/CloudWatch compatibility
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.NameKey = "logger"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.FunctionKey = zapcore.OmitKey
	config.EncoderConfig.StacktraceKey = "stacktrace"

	return config.Build()
}
