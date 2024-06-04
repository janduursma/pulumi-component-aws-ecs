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

	capacityProvidersConfig, err := getCapacityProvidersConfig(sugar)
	if err != nil {
		sugar.Fatal(err)
	}

	pulumi.Run(func(ctx *pulumi.Context) error {
		err = ecs.NewCapacityProviders(ctx, capacityProvidersConfig)
		if err != nil {
			sugar.Error(err)
			return err
		}
		return nil
	})
}

func getCapacityProvidersConfig(sugar *zap.SugaredLogger) ([]ecs.CapacityProviderConfig, error) {
	configData, err := os.ReadFile("config.json")
	if err != nil {
		sugar.Errorf("Error reading config file: %v", err)
		return nil, err
	}

	var tempConfig map[string]interface{}
	err = json.Unmarshal(configData, &tempConfig)
	if err != nil {
		sugar.Error(err)
		return nil, err
	}

	var capacityProviders []ecs.CapacityProviderConfig
	for _, capacityProvider := range tempConfig["capacityProviders"].([]interface{}) {
		capacityProviderConfigJSON, err := json.Marshal(capacityProvider)
		if err != nil {
			sugar.Error(err)
			return nil, err
		}

		var capacityProviderConfig ecs.CapacityProviderConfig
		err = json.Unmarshal(capacityProviderConfigJSON, &capacityProviderConfig)
		if err != nil {
			sugar.Errorf("Error unmarshaling config file: %v", err)
			return nil, err
		}

		capacityProviders = append(capacityProviders, capacityProviderConfig)
	}

	return capacityProviders, nil
}
