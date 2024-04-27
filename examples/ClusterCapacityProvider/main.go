package main

import (
	"encoding/json"
	"fmt"
	"github.com/janduursma/pulumi-component-aws-ecs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"go.uber.org/zap"
	"os"
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

	clusterCapacityProviderConfig, err := getClusterCapacityProviderConfig(sugar)
	if err != nil {
		sugar.Error(err)
		os.Exit(1)
	}

	pulumi.Run(func(ctx *pulumi.Context) error {
		err = ecs.NewClusterCapacityProvider(ctx, *clusterCapacityProviderConfig)
		if err != nil {
			sugar.Error(err)
			return err
		}
		return nil
	})
}

func getClusterCapacityProviderConfig(sugar *zap.SugaredLogger) (*ecs.ClusterCapacityProviderConfig, error) {
	configData, err := os.ReadFile("config.json")
	if err != nil {
		sugar.Errorf("Error reading config file: %v", err)
		return nil, err
	}

	clusterCapacityProviderConfigJSON := make(map[string]*ecs.ClusterCapacityProviderConfig)

	err = json.Unmarshal(configData, &clusterCapacityProviderConfigJSON)
	if err != nil {
		sugar.Errorf("Error unmarshaling config file: %v", err)
		return nil, err
	}

	clusterCapacityProviderConfig, ok := clusterCapacityProviderConfigJSON["clusterCapacityProvider"]
	if !ok {
		err = fmt.Errorf("'clusterCapacityProvider' key not found in JSON config")
		sugar.Error(err)
		return nil, err
	}

	return clusterCapacityProviderConfig, nil
}
