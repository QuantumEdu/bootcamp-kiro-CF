package main

import (
	"context"
	"log"

	"github.com/QuantumEdu/bootcamp-kiro-CF/internal/bootstrap"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/config"
	"github.com/akrylysov/algnhsa"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func main() {
	ctx := context.Background()

	// Load AWS SDK config from the Lambda execution environment.
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("loading AWS config: %v", err)
	}

	// Load secrets from Secrets Manager (DB URL, session key, AI config).
	smClient := secretsmanager.NewFromConfig(awsCfg)
	loader := config.NewSecretsLoader(smClient)
	lambdaCfg, err := loader.LoadConfig(ctx)
	if err != nil {
		log.Fatalf("loading secrets: %v", err)
	}

	// Build bootstrap config for Lambda mode.
	cfg := bootstrap.Config{
		AppEnv:              "lambda",
		DatabaseURL:         lambdaCfg.DatabaseURL,
		SessionSecret:       lambdaCfg.SessionSecret,
		BedrockModelID:      lambdaCfg.BedrockModelID,
		BedrockRegion:       lambdaCfg.BedrockRegion,
		MaxTokens:           lambdaCfg.MaxTokens,
		Temperature:         lambdaCfg.Temperature,
		QueryTimeoutSeconds: 5,
	}

	// Set templates directory relative to WORKDIR (/var/task in the Lambda container).
	bootstrap.TemplatesDir = "templates"

	// Build the chi router with all dependencies.
	router, cleanup, err := bootstrap.BuildRouter(cfg)
	if err != nil {
		log.Fatalf("bootstrap failed: %v", err)
	}
	defer cleanup()

	// Start the Lambda handler — converts API Gateway events to HTTP requests.
	algnhsa.ListenAndServe(router, nil)
}
