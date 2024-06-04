package main

import (
	"encoding/json"
	"fmt"
	"os"

	ecs "github.com/janduursma/pulumi-component-aws-ecs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to sync logger: %v", err)
		}
	}()

	sugar := logger.Sugar()

	serviceConfig, err := getServiceConfig(sugar)
	if err != nil {
		sugar.Fatal(err)
	}

	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err = ecs.NewService(ctx, *serviceConfig)
		if err != nil {
			sugar.Error(err)
			return err
		}
		return nil
	})
}

func getServiceConfig(sugar *zap.SugaredLogger) (*ecs.ServiceConfig, error) {
	configData, err := os.ReadFile("config.json")
	if err != nil {
		sugar.Errorf("Error reading config file: %v", err)
		return nil, err
	}

	serviceConfigJSON := make(map[string]*ecs.ServiceConfig)

	err = json.Unmarshal(configData, &serviceConfigJSON)
	if err != nil {
		sugar.Errorf("Error unmarshaling config file: %v", err)
		return nil, err
	}

	serviceConfig, ok := serviceConfigJSON["service"]
	if !ok {
		err = fmt.Errorf("'service' key not found in JSON config")
		sugar.Error(err)
		return nil, err
	}

	return serviceConfig, nil
}
