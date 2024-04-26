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

	clusterConfig, err := getClusterConfig(sugar)
	if err != nil {
		sugar.Error(err)
		os.Exit(1)
	}

	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err = ecs.NewCluster(ctx, *clusterConfig)
		if err != nil {
			sugar.Error(err)
			return err
		}
		return nil
	})
}

func getClusterConfig(sugar *zap.SugaredLogger) (*ecs.ClusterConfig, error) {
	configData, err := os.ReadFile("config.json")
	if err != nil {
		sugar.Errorf("Error reading config file: %v", err)
		return nil, err
	}

	clusterConfigJSON := make(map[string]*ecs.ClusterConfig)

	err = json.Unmarshal(configData, &clusterConfigJSON)
	if err != nil {
		sugar.Errorf("Error unmarshaling config file: %v", err)
		return nil, err
	}

	clusterConfig, ok := clusterConfigJSON["cluster"]
	if !ok {
		err = fmt.Errorf("'cluster' key not found in JSON config")
		sugar.Error(err)
		return nil, err
	}

	return clusterConfig, nil
}
