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

	taskDefinitionConfig, err := getTaskDefinitionConfig(sugar)
	if err != nil {
		sugar.Fatal(err)
	}

	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err = ecs.NewTaskDefinition(ctx, *taskDefinitionConfig)
		if err != nil {
			sugar.Error(err)
			return err
		}
		return nil
	})
}

func getTaskDefinitionConfig(sugar *zap.SugaredLogger) (*ecs.TaskDefinitionConfig, error) {
	configData, err := os.ReadFile("config.json")
	if err != nil {
		sugar.Errorf("Error reading config file: %v", err)
		return nil, err
	}

	taskDefinitionConfigJSON := make(map[string]*ecs.TaskDefinitionConfig)

	err = json.Unmarshal(configData, &taskDefinitionConfigJSON)
	if err != nil {
		sugar.Errorf("Error unmarshaling config file: %v", err)
		return nil, err
	}

	taskDefinitionConfig, ok := taskDefinitionConfigJSON["taskDefinition"]
	if !ok {
		err = fmt.Errorf("'taskDefinition' key not found in JSON config")
		sugar.Error(err)
		return nil, err
	}

	return taskDefinitionConfig, nil
}
