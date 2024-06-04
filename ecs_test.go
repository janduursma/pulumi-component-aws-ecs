package ecs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/autoscaling"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ecs"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const stackName = "dev"

func manageResources(ctx context.Context, stack auto.Stack, sugar *zap.SugaredLogger, t *testing.T) {
	sugar.Infof("Created/Selected stack: %s", stackName)

	workspace := stack.Workspace()

	sugar.Info("Installing the AWS plugin")
	// for inline source programs, we must manage plugins ourselves
	err := workspace.InstallPlugin(ctx, "aws", "v4.0.0")
	assert.NoError(t, err)
	sugar.Info("Successfully installed AWS plugin!")

	sugar.Info("Setting the stack configuration to use region in AWS")
	err = stack.SetConfig(ctx, "aws:region", auto.ConfigValue{Value: "us-west-2"})
	assert.NoError(t, err)
	sugar.Info("Successfully set config!")

	sugar.Info("Starting refresh")
	_, err = stack.Refresh(ctx)
	assert.NoError(t, err)
	sugar.Info("Refresh succeeded!")

	sugar.Info("Starting update")
	// wire up our update to stream progress to stdout
	stdoutStreamerUp := optup.ProgressStreams(os.Stdout)

	// run the equivalent of 'pulumi up'
	result, err := stack.Up(ctx, stdoutStreamerUp)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	sugar.Info("Update succeeded!")

	sugar.Info("Destroying resources!")
	// wire up our update to stream progress to stdout
	stdoutStreamerDestroy := optdestroy.ProgressStreams(os.Stdout)

	// run the equivalent of 'pulumi destroy'
	_, err = stack.Destroy(ctx, stdoutStreamerDestroy)
	assert.NoError(t, err)
	sugar.Info("Destroy succeeded!")
}

func getAccountSettingsDefaultConfig(sugar *zap.SugaredLogger) ([]AccountSettingDefaultConfig, error) {
	configData, err := os.ReadFile("examples/AccountSettingsDefault/config.json")
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

	var accountSettingsDefault []AccountSettingDefaultConfig
	for _, accountSettingDefault := range tempConfig["accountSettingsDefault"].([]interface{}) {
		accountSettingDefaultConfigJSON, err := json.Marshal(accountSettingDefault)
		if err != nil {
			sugar.Error(err)
			return nil, err
		}

		var accountSettingDefaultConfig AccountSettingDefaultConfig
		err = json.Unmarshal(accountSettingDefaultConfigJSON, &accountSettingDefaultConfig)
		if err != nil {
			sugar.Errorf("Error unmarshaling config file: %v", err)
			return nil, err
		}

		accountSettingsDefault = append(accountSettingsDefault, accountSettingDefaultConfig)
	}

	return accountSettingsDefault, nil
}

// TestNewAccountSettingsDefault is an integration test that checks the correctness of the creation of default account settings for AWS ECS.
// It simulates the process of creating default account settings with defined parameters, which can be found in examples/AccountSettingsDefault/config.json, and expected outcomes.
// The test will pass if the account settings are created successfully.
// Otherwise, it will fail providing information about what incidentally caused the failure.
func TestNewAccountSettingsDefault(t *testing.T) {
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

	sugar.Info("Reading ECS default account settings configuration from examples/AccountSettingsDefault/config.json")
	accountSettingsDefaultConfig, err := getAccountSettingsDefaultConfig(sugar)
	assert.NoError(t, err)
	sugar.Info("Successfully read configuration!")

	ctx := context.Background()
	projectName := "test_ecs_account_settings_default"

	stack, err := auto.UpsertStackInlineSource(ctx, stackName, projectName, func(ctx *pulumi.Context) error {
		_, err = NewAccountSettingsDefault(ctx, accountSettingsDefaultConfig)
		if err != nil {
			return err
		}
		return nil
	})
	assert.NoError(t, err)

	// Set config, run 'pulumi up', and afterwards 'pulumi destroy'
	manageResources(ctx, stack, sugar, t)
}

func createCapacityProvidersTestDependencies(ctx *pulumi.Context) (*autoscaling.Group, error) {
	launchTemplate, err := ec2.NewLaunchTemplate(ctx, "capacityProviderTestDependency", &ec2.LaunchTemplateArgs{
		NamePrefix:   pulumi.String("capacityProvidersTestDependency"),
		ImageId:      pulumi.String("ami-0eb9d67c52f5c80e5"),
		InstanceType: pulumi.String("t2.micro"),
	})
	if err != nil {
		return nil, err
	}

	autoscalingGroup, err := autoscaling.NewGroup(ctx, "capacityProviderTestDependency", &autoscaling.GroupArgs{
		AvailabilityZones: pulumi.StringArray{
			pulumi.String("us-west-2a"),
		},
		DesiredCapacity: pulumi.Int(0),
		MaxSize:         pulumi.Int(0),
		MinSize:         pulumi.Int(0),
		LaunchTemplate: &autoscaling.GroupLaunchTemplateArgs{
			Id:      launchTemplate.ID(),
			Version: pulumi.String("$Latest"),
		},
		ProtectFromScaleIn: pulumi.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	return autoscalingGroup, nil
}

func getCapacityProvidersConfig(sugar *zap.SugaredLogger) ([]CapacityProviderConfig, error) {
	configData, err := os.ReadFile("examples/CapacityProviders/config.json")
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

	var capacityProviders []CapacityProviderConfig
	for _, capacityProvider := range tempConfig["capacityProviders"].([]interface{}) {
		capacityProviderConfigJSON, err := json.Marshal(capacityProvider)
		if err != nil {
			sugar.Error(err)
			return nil, err
		}

		var capacityProviderConfig CapacityProviderConfig
		err = json.Unmarshal(capacityProviderConfigJSON, &capacityProviderConfig)
		if err != nil {
			sugar.Errorf("Error unmarshaling config file: %v", err)
			return nil, err
		}

		capacityProviders = append(capacityProviders, capacityProviderConfig)
	}

	return capacityProviders, nil
}

// TestNewCapacityProviders is an integration test that checks the correctness of the AWS ECS capacity providers creation.
// It simulates the process of creating capacity providers with defined parameters, which can be found in examples/CapacityProviders/config.json, and expected outcomes.
// The test will pass if the capacity providers are created successfully.
// Otherwise, it will fail providing information about what incidentally caused the failure.
func TestNewCapacityProviders(t *testing.T) {
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

	sugar.Info("Reading ECS capacity providers configuration from examples/CapacityProviders/config.json")
	capacityProvidersConfig, err := getCapacityProvidersConfig(sugar)
	assert.NoError(t, err)
	sugar.Info("Successfully read configuration!")

	ctx := context.Background()
	projectName := "test_ecs_capacity_providers_default"

	stack, err := auto.UpsertStackInlineSource(ctx, stackName, projectName, func(ctx *pulumi.Context) error {
		sugar.Info("Creating dependencies for test")
		autoscalingGroup, err := createCapacityProvidersTestDependencies(ctx)
		if err != nil {
			return err
		}
		sugar.Info("Successfully created dependencies!")

		_, err = autoscalingGroup.Arn.ApplyT(func(arn string) (string, error) {
			for i := range capacityProvidersConfig {
				capacityProvidersConfig[i].AutoscalingGroupProvider.AutoscalingGroupArn = arn
			}
			err = NewCapacityProviders(ctx, capacityProvidersConfig)
			if err != nil {
				return "", err
			}
			return arn, nil
		}).(pulumi.StringOutput), nil
		if err != nil {
			return err
		}
		return nil
	})
	assert.NoError(t, err)

	// Set config, run 'pulumi up', and afterwards 'pulumi destroy'
	manageResources(ctx, stack, sugar, t)
}

func getClusterConfig(sugar *zap.SugaredLogger) (*ClusterConfig, error) {
	configData, err := os.ReadFile("examples/Cluster/config.json")
	if err != nil {
		sugar.Errorf("Error reading config file: %v", err)
		return nil, err
	}

	clusterConfigJSON := make(map[string]*ClusterConfig)

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

// TestNewCluster is an integration test that checks the correctness of an AWS ECS cluster creation.
// It simulates the process of creating a cluster with defined parameters, which can be found in examples/Cluster/config.json, and expected outcomes.
// The test will pass if the cluster is created successfully.
// Otherwise, it will fail providing information about what incidentally caused the failure.
func TestNewCluster(t *testing.T) {
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

	sugar.Info("Reading ECS cluster configuration from examples/Cluster/config.json")
	clusterConfig, err := getClusterConfig(sugar)
	assert.NoError(t, err)
	sugar.Info("Successfully read configuration!")

	ctx := context.Background()
	projectName := "test_ecs_cluster"

	stack, err := auto.UpsertStackInlineSource(ctx, stackName, projectName, func(ctx *pulumi.Context) error {
		_, err = NewCluster(ctx, *clusterConfig)
		if err != nil {
			return err
		}
		return nil
	})
	assert.NoError(t, err)

	// Set config, run 'pulumi up', and afterwards 'pulumi destroy'
	manageResources(ctx, stack, sugar, t)
}

func createClusterCapacityProviderTestDependencies(ctx *pulumi.Context) error {
	_, err := ecs.NewCluster(ctx, "clusterCapacityProviderTestDependency", &ecs.ClusterArgs{
		Name: pulumi.String("my-cluster"),
	})
	if err != nil {
		return err
	}
	return nil
}

func getClusterCapacityProviderConfig(sugar *zap.SugaredLogger) (*ClusterCapacityProviderConfig, error) {
	configData, err := os.ReadFile("examples/ClusterCapacityProvider/config.json")
	if err != nil {
		sugar.Errorf("Error reading config file: %v", err)
		return nil, err
	}

	clusterCapacityProviderConfigJSON := make(map[string]*ClusterCapacityProviderConfig)

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

// TestNewClusterCapacityProvider is an integration test that checks the correctness of an AWS ECS cluster capacity provider creation.
// It simulates the process of creating a cluster capacity provider with defined parameters, which can be found in examples/ClusterCapacityProvider/config.json, and expected outcomes.
// The test will pass if the cluster capacity provider is created successfully.
// Otherwise, it will fail providing information about what incidentally caused the failure.
func TestNewClusterCapacityProvider(t *testing.T) {
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

	sugar.Info("Reading ECS cluster capacity provider configuration from examples/ClusterCapacityProvider/config.json")
	clusterCapacityProviderConfig, err := getClusterCapacityProviderConfig(sugar)
	assert.NoError(t, err)
	sugar.Info("Successfully read configuration!")

	ctx := context.Background()
	projectName := "test_ecs_cluster_capacity_provider"

	stack, err := auto.UpsertStackInlineSource(ctx, stackName, projectName, func(ctx *pulumi.Context) error {
		sugar.Info("Creating dependencies for test")
		err = createClusterCapacityProviderTestDependencies(ctx)
		if err != nil {
			return err
		}
		sugar.Info("Successfully created dependencies!")

		err = NewClusterCapacityProvider(ctx, *clusterCapacityProviderConfig)
		if err != nil {
			return err
		}
		return nil
	})
	assert.NoError(t, err)

	// Set config, run 'pulumi up', and afterwards 'pulumi destroy'
	manageResources(ctx, stack, sugar, t)
}

func createServiceTestDependencies(ctx *pulumi.Context, accountID string) error {
	containerDefinitions, err := json.Marshal([]interface{}{
		map[string]interface{}{
			"name":  "my-container",
			"image": "nginx",
			"portMappings": []map[string]interface{}{
				{
					"containerPort": 80,
					"hostPort":      80,
				},
			},
		},
	})
	if err != nil {
		return err
	}

	_, err = ecs.NewTaskDefinition(ctx, "serviceTestDependency", &ecs.TaskDefinitionArgs{
		Family:               pulumi.String("my-task-definition"),
		ContainerDefinitions: pulumi.String(containerDefinitions),
		Cpu:                  pulumi.String("256"),
		Memory:               pulumi.String("512"),
		NetworkMode:          pulumi.String("awsvpc"),
		RequiresCompatibilities: pulumi.StringArray{
			pulumi.String("FARGATE"),
		},
		TaskRoleArn: pulumi.String("arn:aws:iam::" + accountID + ":role/aws-service-role/ecs.amazonaws.com/AWSServiceRoleForECS"),
	})
	if err != nil {
		return err
	}

	_, err = ecs.NewCluster(ctx, "serviceTestDependency", &ecs.ClusterArgs{
		Name: pulumi.String("my-cluster"),
	})
	if err != nil {
		return err
	}
	return nil
}

func getServiceConfig(sugar *zap.SugaredLogger) (*ServiceConfig, error) {
	configData, err := os.ReadFile("examples/Service/config.json")
	if err != nil {
		sugar.Errorf("Error reading config file: %v", err)
		return nil, err
	}

	serviceConfigJSON := make(map[string]*ServiceConfig)

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

// TestNewService is an integration test that checks the correctness of an AWS ECS service creation.
// It simulates the process of creating a service with defined parameters, which can be found in examples/Service/config.json, and expected outcomes.
// The test will pass if the service is created successfully.
// Otherwise, it will fail providing information about what incidentally caused the failure.
func TestNewService(t *testing.T) {
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

	sugar.Info("Reading ECS service configuration from examples/Service/config.json")
	serviceConfig, err := getServiceConfig(sugar)
	assert.NoError(t, err)
	sugar.Info("Successfully read configuration!")

	ctx := context.Background()
	projectName := "test_ecs_service"

	stack, err := auto.UpsertStackInlineSource(ctx, stackName, projectName, func(ctx *pulumi.Context) error {
		current, err := aws.GetCallerIdentity(ctx, nil, nil)
		assert.NoError(t, err)
		serviceConfig.ClusterArn = strings.Replace(serviceConfig.ClusterArn, "$ACCOUNT_ID", current.AccountId, 1)

		err = createServiceTestDependencies(ctx, current.AccountId)
		if err != nil {
			return err
		}

		_, err = NewService(ctx, *serviceConfig)
		if err != nil {
			return err
		}
		return nil
	})
	assert.NoError(t, err)

	// Set config, run 'pulumi up', and afterwards 'pulumi destroy'
	manageResources(ctx, stack, sugar, t)
}

func getTaskDefinitionConfig(sugar *zap.SugaredLogger) (*TaskDefinitionConfig, error) {
	configData, err := os.ReadFile("examples/TaskDefinition/config.json")
	if err != nil {
		sugar.Errorf("Error reading config file: %v", err)
		return nil, err
	}

	taskDefinitionConfigJSON := make(map[string]*TaskDefinitionConfig)

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

// TestNewTaskDefinition is an integration test that checks the correctness of an AWS ECS task definition creation.
// It simulates the process of creating a task definition with defined parameters, which can be found in examples/TaskDefinition/config.json, and expected outcomes.
// The test will pass if the task definition is created successfully.
// Otherwise, it will fail providing information about what incidentally caused the failure.
func TestNewTaskDefinition(t *testing.T) {
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

	sugar.Info("Reading ECS task definition configuration from examples/TaskDefinition/config.json")
	taskDefinitionConfig, err := getTaskDefinitionConfig(sugar)
	assert.NoError(t, err)
	sugar.Info("Successfully read configuration!")

	ctx := context.Background()
	projectName := "test_ecs_task_definition"

	stack, err := auto.UpsertStackInlineSource(ctx, stackName, projectName, func(ctx *pulumi.Context) error {
		current, err := aws.GetCallerIdentity(ctx, nil, nil)
		assert.NoError(t, err)
		taskDefinitionConfig.ExecutionRoleArn = strings.Replace(taskDefinitionConfig.ExecutionRoleArn, "$ACCOUNT_ID", current.AccountId, 1)
		taskDefinitionConfig.TaskRoleArn = strings.Replace(taskDefinitionConfig.TaskRoleArn, "$ACCOUNT_ID", current.AccountId, 1)

		_, err = NewTaskDefinition(ctx, *taskDefinitionConfig)
		if err != nil {
			return err
		}
		return nil
	})
	assert.NoError(t, err)

	// Set config, run 'pulumi up', and afterwards 'pulumi destroy'
	manageResources(ctx, stack, sugar, t)
}

func createTaskSetsTestDependencies(ctx *pulumi.Context) error {
	containerDefinitions, err := json.Marshal([]interface{}{
		map[string]interface{}{
			"name":  "my-container",
			"image": "nginx",
			"portMappings": []map[string]interface{}{
				{
					"containerPort": 80,
					"hostPort":      80,
				},
			},
		},
	})
	if err != nil {
		return err
	}

	current, err := aws.GetCallerIdentity(ctx, nil, nil)
	if err != nil {
		return err
	}

	_, err = ecs.NewTaskDefinition(ctx, "taskSetsTestDependency", &ecs.TaskDefinitionArgs{
		Family:                  pulumi.String("my-task-definition"),
		ContainerDefinitions:    pulumi.String(containerDefinitions),
		Cpu:                     pulumi.String("256"),
		ExecutionRoleArn:        pulumi.String("arn:aws:iam::" + current.AccountId + ":role/aws-service-role/ecs.amazonaws.com/AWSServiceRoleForECS"),
		Memory:                  pulumi.String("512"),
		NetworkMode:             pulumi.String("awsvpc"),
		RequiresCompatibilities: pulumi.ToStringArray([]string{"FARGATE"}),
	})
	if err != nil {
		return err
	}

	_, err = ecs.NewCluster(ctx, "taskSetsTestDependency", &ecs.ClusterArgs{
		Name: pulumi.String("my-cluster"),
	})
	if err != nil {
		return err
	}

	_, err = ecs.NewService(ctx, "taskSetsTestDependency", &ecs.ServiceArgs{
		Name:    pulumi.String("my-service"),
		Cluster: pulumi.String("arn:aws:ecs:us-west-2:" + current.AccountId + ":cluster/my-cluster"),
		DeploymentController: ecs.ServiceDeploymentControllerArgs{
			Type: pulumi.String("EXTERNAL"),
		},
		DesiredCount: pulumi.Int(2),
	})
	if err != nil {
		return err
	}
	return nil
}

func getTaskSetsConfig(sugar *zap.SugaredLogger) ([]TaskSetConfig, error) {
	configData, err := os.ReadFile("examples/TaskSets/config.json")
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

	var taskSets []TaskSetConfig
	for _, taskSet := range tempConfig["taskSets"].([]interface{}) {
		taskSetConfigJSON, err := json.Marshal(taskSet)
		if err != nil {
			sugar.Error(err)
			return nil, err
		}

		var taskSetConfig TaskSetConfig
		err = json.Unmarshal(taskSetConfigJSON, &taskSetConfig)
		if err != nil {
			sugar.Errorf("Error unmarshaling config file: %v", err)
			return nil, err
		}

		taskSets = append(taskSets, taskSetConfig)
	}

	return taskSets, nil
}

// TestNewTaskSets is an integration test that checks the correctness of the creation of AWS ECS task sets.
// It simulates the process of creating task sets with defined parameters, which can be found in examples/TaskSets/config.json, and expected outcomes.
// The test will pass if the task sets are created successfully.
// Otherwise, it will fail providing information about what incidentally caused the failure.
func TestNewTaskSets(t *testing.T) {
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

	sugar.Info("Reading ECS task set configuration from examples/TaskSets/config.json")
	taskSetsConfig, err := getTaskSetsConfig(sugar)
	assert.NoError(t, err)
	sugar.Info("Successfully read configuration!")

	ctx := context.Background()
	projectName := "test_ecs_task_sets"

	stack, err := auto.UpsertStackInlineSource(ctx, stackName, projectName, func(ctx *pulumi.Context) error {
		err = createTaskSetsTestDependencies(ctx)
		if err != nil {
			return err
		}

		_, err = NewTaskSets(ctx, taskSetsConfig)
		if err != nil {
			return err
		}
		return nil
	})
	assert.NoError(t, err)

	// Set config, run 'pulumi up', and afterwards 'pulumi destroy'
	manageResources(ctx, stack, sugar, t)
}
