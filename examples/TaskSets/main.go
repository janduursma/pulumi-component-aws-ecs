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

	taskSetsConfig, err := getTaskSetsConfig(sugar)
	if err != nil {
		sugar.Error(err)
		os.Exit(1)
	}

	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err = ecs.NewTaskSets(ctx, taskSetsConfig)
		if err != nil {
			sugar.Error(err)
			return err
		}
		return nil
	})
}

func getTaskSetsConfig(sugar *zap.SugaredLogger) ([]ecs.TaskSetConfig, error) {
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

	var taskSets []ecs.TaskSetConfig
	for _, taskSet := range tempConfig["taskSets"].([]ecs.TaskSetConfig) {
		taskSetConfigJSON, err := json.Marshal(taskSet)
		if err != nil {
			sugar.Error(err)
			return nil, err
		}

		var taskSetConfig ecs.TaskSetConfig
		err = json.Unmarshal(taskSetConfigJSON, &taskSetConfig)
		if err != nil {
			sugar.Errorf("Error unmarshaling config file: %v", err)
			return nil, err
		}

		taskSets = append(taskSets, taskSetConfig)
	}

	return taskSets, nil
}
