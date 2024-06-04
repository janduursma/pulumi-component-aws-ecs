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

	accountSettingsDefaultConfig, err := getAccountSettingsDefaultConfig(sugar)
	if err != nil {
		sugar.Fatal(err)
	}

	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err = ecs.NewAccountSettingsDefault(ctx, accountSettingsDefaultConfig)
		if err != nil {
			sugar.Error(err)
			return err
		}
		return nil
	})
}

func getAccountSettingsDefaultConfig(sugar *zap.SugaredLogger) ([]ecs.AccountSettingDefaultConfig, error) {
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

	var accountSettingsDefault []ecs.AccountSettingDefaultConfig
	for _, accountSettingDefault := range tempConfig["accountSettingsDefault"].([]interface{}) {
		accountSettingDefaultConfigJSON, err := json.Marshal(accountSettingDefault)
		if err != nil {
			sugar.Error(err)
			return nil, err
		}

		var accountSettingDefaultConfig ecs.AccountSettingDefaultConfig
		err = json.Unmarshal(accountSettingDefaultConfigJSON, &accountSettingDefaultConfig)
		if err != nil {
			sugar.Errorf("Error unmarshaling config file: %v", err)
			return nil, err
		}

		accountSettingsDefault = append(accountSettingsDefault, accountSettingDefaultConfig)
	}

	return accountSettingsDefault, nil
}
