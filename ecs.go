/*
Package ecs is a package that's designed to simplify the creation of AWS ECS (Elastic Container Service) resources using Pulumi.

To get started with this package, import it like so:

	import "github.com/janduursma/pulumi-component-aws-ecs"

Then, create resources by calling functions provided by the package.

Descriptions of the provided functions and their parameters are available in the function doc comments in this package.

Example usage of these functions can be found in the /examples directory at the root of this package repository. These examples provide full, working demonstrations of creating various AWS ECS resources using this package.

For more information, consult the package README or the /examples directory.
*/
package ecs

import (
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ecs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// AccountSettingDefaultConfig defines arguments for AWS ECS account settings.
type AccountSettingDefaultConfig struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// CapacityProviderConfig defines arguments for AWS ECS capacity provider.
type CapacityProviderConfig struct {
	AutoscalingGroupProvider struct {
		AutoscalingGroupArn string  `json:"autoscalingGroupArn"`
		ManagedDraining     *string `json:"managedDraining,omitempty"`
		ManagedScaling      *struct {
			InstanceWarmupPeriod   *int    `json:"instanceWarmupPeriod,omitempty"`
			MaximumScalingStepSize *int    `json:"maximumScalingStepSize,omitempty"`
			MinimumScalingStepSize *int    `json:"minimumScalingStepSize,omitempty"`
			Status                 *string `json:"status,omitempty"`
			TargetCapacity         *int    `json:"targetCapacity,omitempty"`
		} `json:"managedScaling"`
		ManagedTerminationProtection *string `json:"managedTerminationProtection,omitempty"`
	} `json:"autoscalingGroupProvider"`
	Name string            `json:"name"`
	Tags map[string]string `json:"tags,omitempty"`
}

// ClusterConfig defines arguments for creating an AWS ECS cluster.
type ClusterConfig struct {
	Configuration *struct {
		ExecuteCommand struct {
			KmsKeyID         *string `json:"kmsKeyId,omitempty"`
			LogConfiguration *struct {
				CloudWatchEncryptionEnabled *bool   `json:"cloudWatchEncryptionEnabled,omitempty"`
				CloudWatchLogGroupName      *string `json:"cloudWatchLogGroupName,omitempty"`
				S3BucketEncryptionEnabled   *bool   `json:"s3BucketEncryptionEnabled,omitempty"`
				S3BucketName                *string `json:"s3BucketName,omitempty"`
				S3KeyPrefix                 *string `json:"s3KeyPrefix,omitempty"`
			} `json:"logConfiguration"`
			Logging *string `json:"logging,omitempty"`
		} `json:"executeCommand"`
	} `json:"configuration"`
	Name                   string `json:"name"`
	ServiceConnectDefaults *struct {
		Namespace string `json:"namespace"`
	} `json:"serviceConnectDefaults"`
	Settings []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"settings"`
	Tags map[string]string `json:"tags,omitempty"`
}

// clusterOutput defines outputs from the AWS ECS cluster creation.
type clusterOutput struct {
	clusterArn pulumi.StringInput
	id         pulumi.StringInput
}

// ClusterCapacityProviderConfig defines arguments for setting up capacity providers for an AWS ECS cluster.
type ClusterCapacityProviderConfig struct {
	CapacityProviders                 []string `json:"capacityProviders,omitempty"`
	ClusterName                       string   `json:"clusterName"`
	DefaultCapacityProviderStrategies []struct {
		Base             *int   `json:"base,omitempty"`
		CapacityProvider string `json:"capacityProvider"`
		Weight           *int   `json:"weight,omitempty"`
	} `json:"defaultCapacityProviderStrategies"`
}

// ServiceConfig defines arguments for creating an AWS ECS service.
type ServiceConfig struct {
	Alarms *struct {
		AlarmNames []string `json:"alarmNames"`
		Enable     bool     `json:"enable"`
		Rollback   bool     `json:"rollback"`
	} `json:"alarms"`
	CapacityProviderStrategies []struct {
		CapacityProvider string `json:"name"`
		Base             *int   `json:"base,omitempty"`
		Weight           *int   `json:"weight,omitempty"`
	} `json:"capacityProviderStrategies"`
	ClusterArn               string `json:"clusterArn"`
	DeploymentCircuitBreaker *struct {
		Enable   bool `json:"enable"`
		Rollback bool `json:"rollback"`
	} `json:"deploymentCircuitBreaker"`
	DeploymentController *struct {
		Type *string `json:"type,omitempty"`
	} `json:"deploymentController"`
	DeploymentMaximumPercent        *int    `json:"deploymentMaximumPercent,omitempty"`
	DeploymentMinimumHealthyPercent *int    `json:"deploymentMinimumHealthyPercent,omitempty"`
	DesiredCount                    *int    `json:"desiredCount,omitempty"`
	EnableEcsManagedTags            *bool   `json:"enableEcsManagedTags,omitempty"`
	EnableExecuteCommand            *bool   `json:"enableExecuteCommand,omitempty"`
	ForceNewDeployment              *bool   `json:"forceNewDeployment,omitempty"`
	HealthCheckGracePeriodSeconds   *int    `json:"healthCheckGracePeriodSeconds,omitempty"`
	IamRole                         *string `json:"iamRole,omitempty"`
	LaunchType                      *string `json:"launchType,omitempty"`
	LoadBalancers                   []struct {
		ContainerName  string  `json:"containerName"`
		ContainerPort  int     `json:"containerPort"`
		ElbName        *string `json:"elbName,omitempty"`
		TargetGroupArn *string `json:"targetGroupArn,omitempty"`
	} `json:"loadBalancers"`
	Name                 string `json:"name"`
	NetworkConfiguration *struct {
		AssignPublicIP *bool    `json:"assignPublicIp,omitempty"`
		SecurityGroups []string `json:"securityGroups,omitempty"`
		Subnets        []string `json:"subnets"`
	} `json:"networkConfiguration"`
	OrderedPlacementStrategies []struct {
		Field *string `json:"field,omitempty"`
		Type  string  `json:"type"`
	} `json:"orderedPlacementStrategies"`
	PlacementConstraints []struct {
		Expression *string `json:"expression,omitempty"`
		Type       string  `json:"type"`
	} `json:"placementConstraints"`
	PlatformVersion             *string `json:"platformVersion,omitempty"`
	PropagateTags               *string `json:"propagateTags,omitempty"`
	SchedulingStrategy          *string `json:"schedulingStrategy,omitempty"`
	ServiceConnectConfiguration *struct {
		Enabled          bool `json:"enabled"`
		LogConfiguration *struct {
			LogDriver     string            `json:"logDriver"`
			Options       map[string]string `json:"options,omitempty"`
			SecretOptions []struct {
				Name      string `json:"name"`
				ValueFrom string `json:"valueFrom"`
			} `json:"secretOptions"`
		} `json:"logConfiguration"`
		Namespace *string `json:"namespace,omitempty"`
		Services  []struct {
			ClientAlias []struct {
				Port    int    `json:"port"`
				DNSName string `json:"dnsName,omitempty"`
			} `json:"clientAlias"`
			DiscoveryName       *string `json:"discoveryName,omitempty"`
			IngressPortOverride *int    `json:"ingressPortOverride,omitempty"`
			PortName            string  `json:"portName"`
			Timeout             *struct {
				IdleTimeoutSeconds       *int `json:"idleTimeoutSeconds,omitempty"`
				PerRequestTimeoutSeconds *int `json:"perRequestTimeoutSeconds,omitempty"`
			} `json:"timeout"`
			TLS *struct {
				IssuerCertAuthority struct {
					AwsPcaAuthorityArn string `json:"awsPcaAuthorityArn"`
				} `json:"issuerCertAuthority"`
				KmsKey  *string `json:"kmsKey,omitempty"`
				RoleArn *string `json:"roleArn,omitempty"`
			} `json:"tls"`
		} `json:"services"`
	} `json:"serviceConnectConfiguration"`
	ServiceRegistry *struct {
		ContainerName *string `json:"containerName,omitempty"`
		ContainerPort *int    `json:"containerPort,omitempty"`
		Port          *int    `json:"port,omitempty"`
		RegistryArn   string  `json:"registryArn"`
	} `json:"serviceRegistry"`
	ServiceVolumeConfiguration *struct {
		ManagedEBSVolume struct {
			Encrypted      *bool   `json:"encrypted,omitempty"`
			FileSystemType *string `json:"fileSystemType,omitempty"`
			Iops           *int    `json:"iops,omitempty"`
			KmsKeyID       *string `json:"kmsKeyId,omitempty"`
			RoleArn        string  `json:"roleArn"`
			SizeInGB       *int    `json:"sizeInGb,omitempty"`
			SnapshotID     *string `json:"snapshotId,omitempty"`
			Throughput     *string `json:"throughput,omitempty"`
			VolumeType     *string `json:"volumeType,omitempty"`
		} `json:"managedEBSVolume"`
		Name string `json:"name"`
	} `json:"serviceVolumeConfiguration"`
	Tags               map[string]string `json:"tags"`
	TaskDefinition     *string           `json:"taskDefinition,omitempty"`
	Triggers           map[string]string `json:"triggers"`
	WaitForSteadyState *bool             `json:"waitForSteadyState,omitempty"`
}

// serviceOutput defines outputs from the AWS ECS service creation.
type serviceOutput struct {
	id pulumi.StringInput
}

// TaskDefinitionConfig defines arguments for creating an AWS ECS task definition.
type TaskDefinitionConfig struct {
	ContainerDefinitions []map[string]interface{} `json:"containerDefinitions"`
	CPU                  *string                  `json:"cpu,omitempty"`
	EphemeralStorage     *struct {
		SizeInGB int `json:"sizeInGb"`
	} `json:"ephemeralStorage"`
	ExecutionRoleArn      *string `json:"executionRoleArn,omitempty"`
	InferenceAccelerators []struct {
		DeviceName string `json:"deviceName"`
		DeviceType string `json:"deviceType"`
	} `json:"inferenceAccelerators"`
	IpcMode              *string `json:"ipcMode,omitempty"`
	Memory               *string `json:"memory,omitempty"`
	Name                 string  `json:"name"`
	NetworkMode          *string `json:"networkMode,omitempty"`
	PidMode              *string `json:"pidMode,omitempty"`
	PlacementConstraints []struct {
		Expression *string `json:"expression,omitempty"`
		Type       string  `json:"type"`
	} `json:"placementConstraints"`
	ProxyConfiguration *struct {
		ContainerName string            `json:"containerName"`
		Properties    map[string]string `json:"properties"`
		Type          *string           `json:"type,omitempty"`
	} `json:"proxyConfiguration"`
	RequiresCompatibilities []string `json:"requiresCompatibilities"`
	RuntimePlatform         *struct {
		CPUArchitecture       *string `json:"cpuArchitecture,omitempty"`
		OperatingSystemFamily *string `json:"operatingSystemFamily,omitempty"`
	} `json:"runtimePlatform"`
	SkipDestroy *bool             `json:"skipDestroy,omitempty"`
	Tags        map[string]string `json:"tags"`
	TaskRoleArn *string           `json:"taskRoleArn,omitempty"`
	TrackLatest *bool             `json:"trackLatest,omitempty"`
	Volumes     []struct {
		ConfigureAtLaunch         *bool `json:"configureAtLaunch,omitempty"`
		DockerVolumeConfiguration *struct {
			Autoprovision *bool             `json:"autoprovision,omitempty"`
			Driver        *string           `json:"driver,omitempty"`
			DriverOpts    map[string]string `json:"driverOpts"`
			Labels        map[string]string `json:"labels"`
			Scope         *string           `json:"scope,omitempty"`
		} `json:"dockerVolumeConfiguration"`
		EfsVolumeConfiguration *struct {
			AuthorizationConfig *struct {
				AccessPointID *string `json:"accessPointId,omitempty"`
				Iam           *string `json:"iam,omitempty"`
			} `json:"authorizationConfig"`
			FileSystemID          string  `json:"fileSystemId"`
			RootDirectory         *string `json:"rootDirectory,omitempty"`
			TransitEncryption     *string `json:"transitEncryption,omitempty"`
			TransitEncryptionPort *int    `json:"transitEncryptionPort,omitempty"`
		} `json:"efsVolumeConfiguration"`
		FsxWindowsFileServerVolumeConfiguration *struct {
			AuthorizationConfig struct {
				CredentialsParameter string `json:"credentialsParameter"`
				Domain               string `json:"domain"`
			}
			FileSystemID  string `json:"fileSystemId"`
			RootDirectory string `json:"rootDirectory"`
		} `json:"fsxWindowsFileServerVolumeConfiguration"`
		HostPath *string `json:"hostPath,omitempty"`
		Name     string  `json:"name"`
	} `json:"volumes"`
}

// taskDefinitionOutput defines outputs from the AWS ECS task definition creation.
type taskDefinitionOutput struct {
	arn pulumi.StringOutput
}

// TaskSetConfig defines arguments for creating an AWS ECS task set.
type TaskSetConfig struct {
	CapacityProviderStrategies []struct {
		Base             *int   `json:"base,omitempty"`
		CapacityProvider string `json:"capacityProvider"`
		Weight           int    `json:"weight"`
	} `json:"capacityProviderStrategies"`
	Cluster       string  `json:"cluster"`
	ExternalID    *string `json:"externalId,omitempty"`
	ForceDelete   *bool   `json:"forceDelete,omitempty"`
	LaunchType    *string `json:"launchType,omitempty"`
	LoadBalancers []struct {
		ContainerName    string  `json:"containerName"`
		ContainerPort    *int    `json:"containerPort,omitempty"`
		LoadBalancerName *string `json:"loadBalancerName,omitempty"`
		TargetGroupArn   *string `json:"targetGroupArn,omitempty"`
	} `json:"loadBalancers"`
	Name                 string `json:"name"`
	NetworkConfiguration *struct {
		AssignPublicIP *bool    `json:"assignPublicIp,omitempty"`
		SecurityGroups []string `json:"securityGroups"`
		Subnets        []string `json:"subnets"`
	} `json:"networkConfiguration"`
	PlatformVersion *string `json:"platformVersion,omitempty"`
	Scale           *struct {
		Unit  *string  `json:"unit,omitempty"`
		Value *float64 `json:"value,omitempty"`
	} `json:"scale"`
	Service           string `json:"service"`
	ServiceRegistries *struct {
		ContainerName *string `json:"containerName,omitempty"`
		ContainerPort *int    `json:"containerPort,omitempty"`
		Port          *int    `json:"port,omitempty"`
		RegistryArn   string  `json:"registryArn"`
	} `json:"serviceRegistries"`
	Tags                   map[string]string `json:"tags"`
	TaskDefinition         string            `json:"taskDefinition"`
	WaitUntilStable        *bool             `json:"waitUntilStable,omitempty"`
	WaitUntilStableTimeout *string           `json:"waitUntilStableTimeout,omitempty"`
}

// taskSetOutput defines outputs from the AWS ECS task set creation.
type taskSetOutput struct {
	arn pulumi.StringOutput
}

// NewAccountSettingsDefault creates new AWS ECS default account settings.
func NewAccountSettingsDefault(ctx *pulumi.Context, accountSettingsDefault []AccountSettingDefaultConfig, opts ...pulumi.ResourceOption) ([]*ecs.AccountSettingDefault, error) {
	component := &pulumi.ResourceState{}

	var accountSettingsDefaultOutputs []*ecs.AccountSettingDefault
	for i, accountSettingDefault := range accountSettingsDefault {
		err := ctx.RegisterComponentResource("aws:ecs:AccountSettingsDefault", accountSettingDefault.Name, component, opts...)
		if err != nil {
			return accountSettingsDefaultOutputs, fmt.Errorf("failed to register component resource: %v", err)
		}

		output, err := ecs.NewAccountSettingDefault(ctx, fmt.Sprintf("accountSettingDefault-%d", i+1), &ecs.AccountSettingDefaultArgs{
			Name:  pulumi.String(accountSettingDefault.Name),
			Value: pulumi.String(accountSettingDefault.Value),
		}, pulumi.Parent(component))
		if err != nil {
			return accountSettingsDefaultOutputs, fmt.Errorf("failed to create new default account setting: %v", err)
		}
		accountSettingsDefaultOutputs = append(accountSettingsDefaultOutputs, output)
	}

	return accountSettingsDefaultOutputs, nil
}

// NewCapacityProviders creates new ECS capacity providers.
func NewCapacityProviders(ctx *pulumi.Context, capacityProviders []CapacityProviderConfig, opts ...pulumi.ResourceOption) error {
	component := &pulumi.ResourceState{}

	for i, capacityProvider := range capacityProviders {
		err := ctx.RegisterComponentResource("aws:ecs:capacityProvider", capacityProvider.Name, component, opts...)
		if err != nil {
			return fmt.Errorf("failed to register component resource: %v", err)
		}

		var managedScaling *ecs.CapacityProviderAutoScalingGroupProviderManagedScalingArgs
		if capacityProvider.AutoscalingGroupProvider.ManagedScaling != nil {
			managedScaling = &ecs.CapacityProviderAutoScalingGroupProviderManagedScalingArgs{
				InstanceWarmupPeriod:   pulumi.IntPtrFromPtr(capacityProvider.AutoscalingGroupProvider.ManagedScaling.InstanceWarmupPeriod),
				MaximumScalingStepSize: pulumi.IntPtrFromPtr(capacityProvider.AutoscalingGroupProvider.ManagedScaling.MaximumScalingStepSize),
				MinimumScalingStepSize: pulumi.IntPtrFromPtr(capacityProvider.AutoscalingGroupProvider.ManagedScaling.MinimumScalingStepSize),
				Status:                 pulumi.StringPtrFromPtr(capacityProvider.AutoscalingGroupProvider.ManagedScaling.Status),
				TargetCapacity:         pulumi.IntPtrFromPtr(capacityProvider.AutoscalingGroupProvider.ManagedScaling.TargetCapacity),
			}
		}

		_, err = ecs.NewCapacityProvider(ctx, fmt.Sprintf("capacityProvider-%d", i+1), &ecs.CapacityProviderArgs{
			AutoScalingGroupProvider: &ecs.CapacityProviderAutoScalingGroupProviderArgs{
				AutoScalingGroupArn:          pulumi.String(capacityProvider.AutoscalingGroupProvider.AutoscalingGroupArn),
				ManagedDraining:              pulumi.StringPtrFromPtr(capacityProvider.AutoscalingGroupProvider.ManagedDraining),
				ManagedScaling:               managedScaling,
				ManagedTerminationProtection: pulumi.StringPtrFromPtr(capacityProvider.AutoscalingGroupProvider.ManagedTerminationProtection),
			},
			Name: pulumi.String(capacityProvider.Name),
			Tags: pulumi.ToStringMap(capacityProvider.Tags),
		}, pulumi.Parent(component))
		if err != nil {
			return fmt.Errorf("failed to create new capacity provider: %v", err)
		}
	}

	return nil
}

// NewCluster creates a new ECS cluster.
func NewCluster(ctx *pulumi.Context, config ClusterConfig, opts ...pulumi.ResourceOption) (*clusterOutput, error) {
	component := &pulumi.ResourceState{}
	err := ctx.RegisterComponentResource("aws:ecs:Cluster", config.Name, component, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to register component resource: %v", err)
	}

	var settings ecs.ClusterSettingArray
	for _, setting := range config.Settings {
		settings = append(settings, &ecs.ClusterSettingArgs{
			Name:  pulumi.String(setting.Name),
			Value: pulumi.String(setting.Value),
		})
	}

	var configuration *ecs.ClusterConfigurationArgs
	if config.Configuration != nil {
		var logConfiguration *ecs.ClusterConfigurationExecuteCommandConfigurationLogConfigurationArgs
		if config.Configuration.ExecuteCommand.LogConfiguration != nil {
			logConfiguration = &ecs.ClusterConfigurationExecuteCommandConfigurationLogConfigurationArgs{
				CloudWatchEncryptionEnabled: pulumi.BoolPtrFromPtr(config.Configuration.ExecuteCommand.LogConfiguration.CloudWatchEncryptionEnabled),
				CloudWatchLogGroupName:      pulumi.StringPtrFromPtr(config.Configuration.ExecuteCommand.LogConfiguration.CloudWatchLogGroupName),
				S3BucketEncryptionEnabled:   pulumi.BoolPtrFromPtr(config.Configuration.ExecuteCommand.LogConfiguration.S3BucketEncryptionEnabled),
				S3BucketName:                pulumi.StringPtrFromPtr(config.Configuration.ExecuteCommand.LogConfiguration.S3BucketName),
				S3KeyPrefix:                 pulumi.StringPtrFromPtr(config.Configuration.ExecuteCommand.LogConfiguration.S3KeyPrefix),
			}
		}
		configuration = &ecs.ClusterConfigurationArgs{
			ExecuteCommandConfiguration: &ecs.ClusterConfigurationExecuteCommandConfigurationArgs{
				KmsKeyId:         pulumi.StringPtrFromPtr(config.Configuration.ExecuteCommand.KmsKeyID),
				LogConfiguration: logConfiguration,
				Logging:          pulumi.StringPtrFromPtr(config.Configuration.ExecuteCommand.Logging),
			},
		}
	}

	var serviceConnectDefaults *ecs.ClusterServiceConnectDefaultsArgs
	if config.ServiceConnectDefaults != nil {
		serviceConnectDefaults = &ecs.ClusterServiceConnectDefaultsArgs{
			Namespace: pulumi.String(config.ServiceConnectDefaults.Namespace),
		}
	}

	cluster, err := ecs.NewCluster(ctx, "cluster", &ecs.ClusterArgs{
		Configuration:          configuration,
		Name:                   pulumi.String(config.Name),
		ServiceConnectDefaults: serviceConnectDefaults,
		Settings:               settings,
		Tags:                   pulumi.ToStringMap(config.Tags),
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create new cluster: %v", err)
	}

	return &clusterOutput{
		clusterArn: cluster.Arn,
		id:         cluster.ID(),
	}, nil
}

// NewClusterCapacityProvider creates a new capacity provider for an AWS ECS cluster.
func NewClusterCapacityProvider(ctx *pulumi.Context, config ClusterCapacityProviderConfig, opts ...pulumi.ResourceOption) error {
	component := &pulumi.ResourceState{}
	err := ctx.RegisterComponentResource("aws:ecs:clusterCapacityProvider", config.ClusterName, component, opts...)
	if err != nil {
		return fmt.Errorf("failed to register component resource: %v", err)
	}

	var defaultCapacityProviderStrategies ecs.ClusterCapacityProvidersDefaultCapacityProviderStrategyArray
	for _, defaultCapacityProviderStrategy := range config.DefaultCapacityProviderStrategies {
		defaultCapacityProviderStrategies = append(defaultCapacityProviderStrategies, &ecs.ClusterCapacityProvidersDefaultCapacityProviderStrategyArgs{
			Base:             pulumi.IntPtrFromPtr(defaultCapacityProviderStrategy.Base),
			CapacityProvider: pulumi.String(defaultCapacityProviderStrategy.CapacityProvider),
			Weight:           pulumi.IntPtrFromPtr(defaultCapacityProviderStrategy.Weight),
		})
	}

	_, err = ecs.NewClusterCapacityProviders(ctx, "clusterCapacityProviders", &ecs.ClusterCapacityProvidersArgs{
		CapacityProviders:                 pulumi.ToStringArray(config.CapacityProviders),
		ClusterName:                       pulumi.String(config.ClusterName),
		DefaultCapacityProviderStrategies: defaultCapacityProviderStrategies,
	}, pulumi.Parent(component))
	if err != nil {
		return fmt.Errorf("failed to create new cluster capacity provider: %v", err)
	}

	return nil
}

func createServiceConnectConfiguration(config ServiceConfig) *ecs.ServiceServiceConnectConfigurationArgs {
	var serviceConnectConfiguration *ecs.ServiceServiceConnectConfigurationArgs
	if config.ServiceConnectConfiguration != nil {
		var serviceConnectServicesLogConfiguration *ecs.ServiceServiceConnectConfigurationLogConfigurationArgs
		var serviceConnectLogsSecretOptions ecs.ServiceServiceConnectConfigurationLogConfigurationSecretOptionArray
		if config.ServiceConnectConfiguration.LogConfiguration != nil {
			for _, serviceConnectLogsSecretOption := range config.ServiceConnectConfiguration.LogConfiguration.SecretOptions {
				serviceConnectLogsSecretOptions = append(serviceConnectLogsSecretOptions, &ecs.ServiceServiceConnectConfigurationLogConfigurationSecretOptionArgs{
					Name:      pulumi.String(serviceConnectLogsSecretOption.Name),
					ValueFrom: pulumi.String(serviceConnectLogsSecretOption.ValueFrom),
				})
			}

			serviceConnectServicesLogConfiguration = &ecs.ServiceServiceConnectConfigurationLogConfigurationArgs{
				LogDriver:     pulumi.String(config.ServiceConnectConfiguration.LogConfiguration.LogDriver),
				Options:       pulumi.ToStringMap(config.ServiceConnectConfiguration.LogConfiguration.Options),
				SecretOptions: serviceConnectLogsSecretOptions,
			}
		}

		var serviceConnectServicesClientAliases ecs.ServiceServiceConnectConfigurationServiceClientAliasArray
		for _, service := range config.ServiceConnectConfiguration.Services {
			for _, clientAlias := range service.ClientAlias {
				serviceConnectServicesClientAliases = append(serviceConnectServicesClientAliases, &ecs.ServiceServiceConnectConfigurationServiceClientAliasArgs{
					DnsName: pulumi.String(clientAlias.DNSName),
					Port:    pulumi.Int(clientAlias.Port),
				})
			}
		}

		var serviceConnectServices ecs.ServiceServiceConnectConfigurationServiceArray
		for _, serviceConnectService := range config.ServiceConnectConfiguration.Services {
			var timeout *ecs.ServiceServiceConnectConfigurationServiceTimeoutArgs
			if serviceConnectService.Timeout != nil {
				timeout = &ecs.ServiceServiceConnectConfigurationServiceTimeoutArgs{
					IdleTimeoutSeconds:       pulumi.IntPtrFromPtr(serviceConnectService.Timeout.IdleTimeoutSeconds),
					PerRequestTimeoutSeconds: pulumi.IntPtrFromPtr(serviceConnectService.Timeout.PerRequestTimeoutSeconds),
				}
			}

			var tls *ecs.ServiceServiceConnectConfigurationServiceTlsArgs
			if serviceConnectService.TLS != nil {
				tls = &ecs.ServiceServiceConnectConfigurationServiceTlsArgs{
					IssuerCertAuthority: &ecs.ServiceServiceConnectConfigurationServiceTlsIssuerCertAuthorityArgs{
						AwsPcaAuthorityArn: pulumi.String(serviceConnectService.TLS.IssuerCertAuthority.AwsPcaAuthorityArn),
					},
					KmsKey:  pulumi.StringPtrFromPtr(serviceConnectService.TLS.KmsKey),
					RoleArn: pulumi.StringPtrFromPtr(serviceConnectService.TLS.RoleArn),
				}
			}

			serviceConnectServices = append(serviceConnectServices, &ecs.ServiceServiceConnectConfigurationServiceArgs{
				ClientAlias:         serviceConnectServicesClientAliases,
				DiscoveryName:       pulumi.StringPtrFromPtr(serviceConnectService.DiscoveryName),
				IngressPortOverride: pulumi.IntPtrFromPtr(serviceConnectService.IngressPortOverride),
				PortName:            pulumi.String(serviceConnectService.PortName),
				Timeout:             timeout,
				Tls:                 tls,
			})
		}

		serviceConnectConfiguration = &ecs.ServiceServiceConnectConfigurationArgs{
			Enabled:          pulumi.Bool(config.ServiceConnectConfiguration.Enabled),
			LogConfiguration: serviceConnectServicesLogConfiguration,
			Namespace:        pulumi.StringPtrFromPtr(config.ServiceConnectConfiguration.Namespace),
			Services:         serviceConnectServices,
		}
	}

	return serviceConnectConfiguration
}

// NewService creates a new AWS ECS service.
func NewService(ctx *pulumi.Context, config ServiceConfig, opts ...pulumi.ResourceOption) (*serviceOutput, error) {
	component := &pulumi.ResourceState{}
	err := ctx.RegisterComponentResource("aws:ecs:Service", config.Name, component, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to register component resource: %v", err)
	}

	var alarms *ecs.ServiceAlarmsArgs
	if config.Alarms != nil {
		alarms = &ecs.ServiceAlarmsArgs{
			AlarmNames: pulumi.ToStringArray(config.Alarms.AlarmNames),
			Enable:     pulumi.Bool(config.Alarms.Enable),
			Rollback:   pulumi.Bool(config.Alarms.Rollback),
		}
	}

	var capacityProviderStrategies ecs.ServiceCapacityProviderStrategyArray
	for _, capacityProviderStrategy := range config.CapacityProviderStrategies {
		capacityProviderStrategies = append(capacityProviderStrategies, &ecs.ServiceCapacityProviderStrategyArgs{
			CapacityProvider: pulumi.String(capacityProviderStrategy.CapacityProvider),
			Base:             pulumi.IntPtrFromPtr(capacityProviderStrategy.Base),
			Weight:           pulumi.IntPtrFromPtr(capacityProviderStrategy.Weight),
		})
	}

	var deploymentCircuitBreaker *ecs.ServiceDeploymentCircuitBreakerArgs
	if config.DeploymentCircuitBreaker != nil {
		deploymentCircuitBreaker = &ecs.ServiceDeploymentCircuitBreakerArgs{
			Enable:   pulumi.Bool(config.DeploymentCircuitBreaker.Enable),
			Rollback: pulumi.Bool(config.DeploymentCircuitBreaker.Rollback),
		}
	}

	var deploymentController *ecs.ServiceDeploymentControllerArgs
	if config.DeploymentController != nil {
		deploymentController = &ecs.ServiceDeploymentControllerArgs{
			Type: pulumi.StringPtrFromPtr(config.DeploymentController.Type),
		}
	}

	var loadBalancers ecs.ServiceLoadBalancerArray
	for _, loadBalancer := range config.LoadBalancers {
		loadBalancers = append(loadBalancers, &ecs.ServiceLoadBalancerArgs{
			ContainerName:  pulumi.String(loadBalancer.ContainerName),
			ContainerPort:  pulumi.Int(loadBalancer.ContainerPort),
			ElbName:        pulumi.StringPtrFromPtr(loadBalancer.ElbName),
			TargetGroupArn: pulumi.StringPtrFromPtr(loadBalancer.TargetGroupArn),
		})
	}

	var networkConfiguration *ecs.ServiceNetworkConfigurationArgs
	if config.NetworkConfiguration != nil {
		networkConfiguration = &ecs.ServiceNetworkConfigurationArgs{
			AssignPublicIp: pulumi.BoolPtrFromPtr(config.NetworkConfiguration.AssignPublicIP),
			SecurityGroups: pulumi.ToStringArray(config.NetworkConfiguration.SecurityGroups),
			Subnets:        pulumi.ToStringArray(config.NetworkConfiguration.Subnets),
		}
	}

	var orderedPlacementStrategies ecs.ServiceOrderedPlacementStrategyArray
	for _, orderedPlacementStrategy := range config.OrderedPlacementStrategies {
		orderedPlacementStrategies = append(orderedPlacementStrategies, &ecs.ServiceOrderedPlacementStrategyArgs{
			Field: pulumi.StringPtrFromPtr(orderedPlacementStrategy.Field),
			Type:  pulumi.String(orderedPlacementStrategy.Type),
		})
	}

	var placementConstraints ecs.ServicePlacementConstraintArray
	for _, placementConstraint := range config.PlacementConstraints {
		placementConstraints = append(placementConstraints, &ecs.ServicePlacementConstraintArgs{
			Expression: pulumi.StringPtrFromPtr(placementConstraint.Expression),
			Type:       pulumi.String(placementConstraint.Type),
		})
	}

	var serviceRegistries *ecs.ServiceServiceRegistriesArgs
	if config.ServiceRegistry != nil {
		serviceRegistries = &ecs.ServiceServiceRegistriesArgs{
			ContainerName: pulumi.StringPtrFromPtr(config.ServiceRegistry.ContainerName),
			ContainerPort: pulumi.IntPtrFromPtr(config.ServiceRegistry.ContainerPort),
			Port:          pulumi.IntPtrFromPtr(config.ServiceRegistry.Port),
			RegistryArn:   pulumi.String(config.ServiceRegistry.RegistryArn),
		}
	}

	serviceConnectConfiguration := createServiceConnectConfiguration(config)

	var serviceVolumeConfiguration *ecs.ServiceVolumeConfigurationArgs
	if config.ServiceVolumeConfiguration != nil {
		serviceVolumeConfiguration = &ecs.ServiceVolumeConfigurationArgs{
			ManagedEbsVolume: &ecs.ServiceVolumeConfigurationManagedEbsVolumeArgs{
				Encrypted:      pulumi.BoolPtrFromPtr(config.ServiceVolumeConfiguration.ManagedEBSVolume.Encrypted),
				FileSystemType: pulumi.StringPtrFromPtr(config.ServiceVolumeConfiguration.ManagedEBSVolume.FileSystemType),
				Iops:           pulumi.IntPtrFromPtr(config.ServiceVolumeConfiguration.ManagedEBSVolume.Iops),
				KmsKeyId:       pulumi.StringPtrFromPtr(config.ServiceVolumeConfiguration.ManagedEBSVolume.KmsKeyID),
				RoleArn:        pulumi.String(config.ServiceVolumeConfiguration.ManagedEBSVolume.RoleArn),
				SizeInGb:       pulumi.IntPtrFromPtr(config.ServiceVolumeConfiguration.ManagedEBSVolume.SizeInGB),
				SnapshotId:     pulumi.StringPtrFromPtr(config.ServiceVolumeConfiguration.ManagedEBSVolume.SnapshotID),
				Throughput:     pulumi.StringPtrFromPtr(config.ServiceVolumeConfiguration.ManagedEBSVolume.Throughput),
				VolumeType:     pulumi.StringPtrFromPtr(config.ServiceVolumeConfiguration.ManagedEBSVolume.VolumeType),
			},
			Name: pulumi.String(config.ServiceVolumeConfiguration.Name),
		}
	}

	service, err := ecs.NewService(ctx, "service", &ecs.ServiceArgs{
		Alarms:                          alarms,
		CapacityProviderStrategies:      capacityProviderStrategies,
		Cluster:                         pulumi.String(config.ClusterArn),
		DeploymentCircuitBreaker:        deploymentCircuitBreaker,
		DeploymentController:            deploymentController,
		DeploymentMaximumPercent:        pulumi.IntPtrFromPtr(config.DeploymentMaximumPercent),
		DeploymentMinimumHealthyPercent: pulumi.IntPtrFromPtr(config.DeploymentMinimumHealthyPercent),
		DesiredCount:                    pulumi.IntPtrFromPtr(config.DesiredCount),
		EnableEcsManagedTags:            pulumi.BoolPtrFromPtr(config.EnableEcsManagedTags),
		EnableExecuteCommand:            pulumi.BoolPtrFromPtr(config.EnableExecuteCommand),
		ForceNewDeployment:              pulumi.BoolPtrFromPtr(config.ForceNewDeployment),
		HealthCheckGracePeriodSeconds:   pulumi.IntPtrFromPtr(config.HealthCheckGracePeriodSeconds),
		IamRole:                         pulumi.StringPtrFromPtr(config.IamRole),
		LaunchType:                      pulumi.StringPtrFromPtr(config.LaunchType),
		LoadBalancers:                   loadBalancers,
		Name:                            pulumi.String(config.Name),
		NetworkConfiguration:            networkConfiguration,
		OrderedPlacementStrategies:      orderedPlacementStrategies,
		PlacementConstraints:            placementConstraints,
		PlatformVersion:                 pulumi.StringPtrFromPtr(config.PlatformVersion),
		PropagateTags:                   pulumi.StringPtrFromPtr(config.PropagateTags),
		SchedulingStrategy:              pulumi.StringPtrFromPtr(config.SchedulingStrategy),
		ServiceConnectConfiguration:     serviceConnectConfiguration,
		ServiceRegistries:               serviceRegistries,
		Tags:                            pulumi.ToStringMap(config.Tags),
		TaskDefinition:                  pulumi.StringPtrFromPtr(config.TaskDefinition),
		Triggers:                        pulumi.ToStringMap(config.Triggers),
		VolumeConfiguration:             serviceVolumeConfiguration,
		WaitForSteadyState:              pulumi.BoolPtrFromPtr(config.WaitForSteadyState),
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create new service: %v", err)
	}

	return &serviceOutput{
		id: service.ID(),
	}, nil
}

// NewTaskDefinition creates a new AWS ECS task definition.
func NewTaskDefinition(ctx *pulumi.Context, config TaskDefinitionConfig, opts ...pulumi.ResourceOption) (*taskDefinitionOutput, error) {
	component := &pulumi.ResourceState{}
	err := ctx.RegisterComponentResource("aws:ecs:TaskDefinition", config.Name, component, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to register component resource: %v", err)
	}

	containerDefinitions, err := json.Marshal(config.ContainerDefinitions)
	if err != nil {
		return nil, fmt.Errorf("could not marshal container definitions json: %v", err)
	}

	var ephemeralStorage *ecs.TaskDefinitionEphemeralStorageArgs
	if config.EphemeralStorage != nil {
		ephemeralStorage = &ecs.TaskDefinitionEphemeralStorageArgs{
			SizeInGib: pulumi.Int(config.EphemeralStorage.SizeInGB),
		}
	}

	var inferenceAccelerators ecs.TaskDefinitionInferenceAcceleratorArray
	for _, inferenceAccelerator := range config.InferenceAccelerators {
		inferenceAccelerators = append(inferenceAccelerators, &ecs.TaskDefinitionInferenceAcceleratorArgs{
			DeviceName: pulumi.String(inferenceAccelerator.DeviceName),
			DeviceType: pulumi.String(inferenceAccelerator.DeviceType),
		})
	}

	var placementConstraints ecs.TaskDefinitionPlacementConstraintArray
	for _, placementConstraint := range config.PlacementConstraints {
		placementConstraints = append(placementConstraints, &ecs.TaskDefinitionPlacementConstraintArgs{
			Expression: pulumi.StringPtrFromPtr(placementConstraint.Expression),
			Type:       pulumi.String(placementConstraint.Type),
		})
	}

	var proxyConfiguration *ecs.TaskDefinitionProxyConfigurationArgs
	if config.ProxyConfiguration != nil {
		proxyConfiguration = &ecs.TaskDefinitionProxyConfigurationArgs{
			ContainerName: pulumi.String(config.ProxyConfiguration.ContainerName),
			Properties:    pulumi.ToStringMap(config.ProxyConfiguration.Properties),
			Type:          pulumi.StringPtrFromPtr(config.ProxyConfiguration.Type),
		}
	}

	var runtimePlatform *ecs.TaskDefinitionRuntimePlatformArgs
	if config.RuntimePlatform != nil {
		runtimePlatform = &ecs.TaskDefinitionRuntimePlatformArgs{
			CpuArchitecture:       pulumi.StringPtrFromPtr(config.RuntimePlatform.CPUArchitecture),
			OperatingSystemFamily: pulumi.StringPtrFromPtr(config.RuntimePlatform.OperatingSystemFamily),
		}
	}

	var volumes ecs.TaskDefinitionVolumeArray
	for _, volume := range config.Volumes {
		var dockerVolumeConfiguration *ecs.TaskDefinitionVolumeDockerVolumeConfigurationArgs
		if volume.DockerVolumeConfiguration != nil {
			dockerVolumeConfiguration = &ecs.TaskDefinitionVolumeDockerVolumeConfigurationArgs{
				Autoprovision: pulumi.BoolPtrFromPtr(volume.DockerVolumeConfiguration.Autoprovision),
				Driver:        pulumi.StringPtrFromPtr(volume.DockerVolumeConfiguration.Driver),
				DriverOpts:    pulumi.ToStringMap(volume.DockerVolumeConfiguration.DriverOpts),
				Labels:        pulumi.ToStringMap(volume.DockerVolumeConfiguration.Labels),
				Scope:         pulumi.StringPtrFromPtr(volume.DockerVolumeConfiguration.Scope),
			}
		}

		var efsVolumeConfiguration *ecs.TaskDefinitionVolumeEfsVolumeConfigurationArgs
		if volume.EfsVolumeConfiguration != nil {
			var authorizationConfig *ecs.TaskDefinitionVolumeEfsVolumeConfigurationAuthorizationConfigArgs
			if volume.EfsVolumeConfiguration.AuthorizationConfig != nil {
				authorizationConfig = &ecs.TaskDefinitionVolumeEfsVolumeConfigurationAuthorizationConfigArgs{
					AccessPointId: pulumi.StringPtrFromPtr(volume.EfsVolumeConfiguration.AuthorizationConfig.AccessPointID),
					Iam:           pulumi.StringPtrFromPtr(volume.EfsVolumeConfiguration.AuthorizationConfig.Iam),
				}
			}

			efsVolumeConfiguration = &ecs.TaskDefinitionVolumeEfsVolumeConfigurationArgs{
				AuthorizationConfig:   authorizationConfig,
				FileSystemId:          pulumi.String(volume.EfsVolumeConfiguration.FileSystemID),
				RootDirectory:         pulumi.StringPtrFromPtr(volume.EfsVolumeConfiguration.RootDirectory),
				TransitEncryption:     pulumi.StringPtrFromPtr(volume.EfsVolumeConfiguration.TransitEncryption),
				TransitEncryptionPort: pulumi.IntPtrFromPtr(volume.EfsVolumeConfiguration.TransitEncryptionPort),
			}
		}

		var fsxVolumeConfiguration *ecs.TaskDefinitionVolumeFsxWindowsFileServerVolumeConfigurationArgs
		if volume.FsxWindowsFileServerVolumeConfiguration != nil {
			fsxVolumeConfiguration = &ecs.TaskDefinitionVolumeFsxWindowsFileServerVolumeConfigurationArgs{
				AuthorizationConfig: &ecs.TaskDefinitionVolumeFsxWindowsFileServerVolumeConfigurationAuthorizationConfigArgs{
					CredentialsParameter: pulumi.String(volume.FsxWindowsFileServerVolumeConfiguration.AuthorizationConfig.CredentialsParameter),
					Domain:               pulumi.String(volume.FsxWindowsFileServerVolumeConfiguration.AuthorizationConfig.Domain),
				},
				FileSystemId:  pulumi.String(volume.FsxWindowsFileServerVolumeConfiguration.FileSystemID),
				RootDirectory: pulumi.String(volume.FsxWindowsFileServerVolumeConfiguration.RootDirectory),
			}
		}

		volumes = append(volumes, &ecs.TaskDefinitionVolumeArgs{
			DockerVolumeConfiguration:               dockerVolumeConfiguration,
			EfsVolumeConfiguration:                  efsVolumeConfiguration,
			FsxWindowsFileServerVolumeConfiguration: fsxVolumeConfiguration,
			HostPath:                                pulumi.StringPtrFromPtr(volume.HostPath),
			Name:                                    pulumi.String(volume.Name),
		})
	}

	taskDefinition, err := ecs.NewTaskDefinition(ctx, "taskDefinition", &ecs.TaskDefinitionArgs{
		ContainerDefinitions:    pulumi.String(containerDefinitions),
		Cpu:                     pulumi.StringPtrFromPtr(config.CPU),
		EphemeralStorage:        ephemeralStorage,
		ExecutionRoleArn:        pulumi.StringPtrFromPtr(config.ExecutionRoleArn),
		Family:                  pulumi.String(config.Name),
		IpcMode:                 pulumi.StringPtrFromPtr(config.IpcMode),
		InferenceAccelerators:   inferenceAccelerators,
		Memory:                  pulumi.StringPtrFromPtr(config.Memory),
		NetworkMode:             pulumi.StringPtrFromPtr(config.NetworkMode),
		PidMode:                 pulumi.StringPtrFromPtr(config.PidMode),
		PlacementConstraints:    placementConstraints,
		ProxyConfiguration:      proxyConfiguration,
		RequiresCompatibilities: pulumi.ToStringArray(config.RequiresCompatibilities),
		RuntimePlatform:         runtimePlatform,
		SkipDestroy:             pulumi.BoolPtrFromPtr(config.SkipDestroy),
		Tags:                    pulumi.ToStringMap(config.Tags),
		TaskRoleArn:             pulumi.StringPtrFromPtr(config.TaskRoleArn),
		TrackLatest:             pulumi.BoolPtrFromPtr(config.TrackLatest),
		Volumes:                 volumes,
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("failed to create new task definition: %v", err)
	}

	return &taskDefinitionOutput{
		arn: taskDefinition.Arn,
	}, nil
}

// NewTaskSets creates new AWS ECS task sets.
func NewTaskSets(ctx *pulumi.Context, taskSets []TaskSetConfig, opts ...pulumi.ResourceOption) ([]*taskSetOutput, error) {
	component := &pulumi.ResourceState{}
	var taskSetOutputs []*taskSetOutput

	for i, taskSet := range taskSets {
		err := ctx.RegisterComponentResource("aws:ecs:TaskSet", taskSet.Name, component, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to register component resource: %v", err)
		}

		var capacityProviderStrategies ecs.TaskSetCapacityProviderStrategyArray
		for _, capacityProviderStrategy := range taskSet.CapacityProviderStrategies {
			capacityProviderStrategies = append(capacityProviderStrategies, &ecs.TaskSetCapacityProviderStrategyArgs{
				Base:             pulumi.IntPtrFromPtr(capacityProviderStrategy.Base),
				CapacityProvider: pulumi.String(capacityProviderStrategy.CapacityProvider),
				Weight:           pulumi.Int(capacityProviderStrategy.Weight),
			})
		}

		var loadBalancers ecs.TaskSetLoadBalancerArray
		for _, loadBalancer := range taskSet.LoadBalancers {
			loadBalancers = append(loadBalancers, &ecs.TaskSetLoadBalancerArgs{
				ContainerName:    pulumi.String(loadBalancer.ContainerName),
				ContainerPort:    pulumi.IntPtrFromPtr(loadBalancer.ContainerPort),
				LoadBalancerName: pulumi.StringPtrFromPtr(loadBalancer.LoadBalancerName),
				TargetGroupArn:   pulumi.StringPtrFromPtr(loadBalancer.TargetGroupArn),
			})
		}

		var networkConfiguration *ecs.TaskSetNetworkConfigurationArgs
		if taskSet.NetworkConfiguration != nil {
			networkConfiguration = &ecs.TaskSetNetworkConfigurationArgs{
				AssignPublicIp: pulumi.BoolPtrFromPtr(taskSet.NetworkConfiguration.AssignPublicIP),
				SecurityGroups: pulumi.ToStringArray(taskSet.NetworkConfiguration.SecurityGroups),
				Subnets:        pulumi.ToStringArray(taskSet.NetworkConfiguration.Subnets),
			}
		}

		var scale *ecs.TaskSetScaleArgs
		if taskSet.Scale != nil {
			scale = &ecs.TaskSetScaleArgs{
				Unit:  pulumi.StringPtrFromPtr(taskSet.Scale.Unit),
				Value: pulumi.Float64PtrFromPtr(taskSet.Scale.Value),
			}
		}

		var serviceRegistries *ecs.TaskSetServiceRegistriesArgs
		if taskSet.ServiceRegistries != nil {
			serviceRegistries = &ecs.TaskSetServiceRegistriesArgs{
				ContainerName: pulumi.StringPtrFromPtr(taskSet.ServiceRegistries.ContainerName),
				ContainerPort: pulumi.IntPtrFromPtr(taskSet.ServiceRegistries.ContainerPort),
				Port:          pulumi.IntPtrFromPtr(taskSet.ServiceRegistries.Port),
				RegistryArn:   pulumi.String(taskSet.ServiceRegistries.RegistryArn),
			}
		}

		output, err := ecs.NewTaskSet(ctx, fmt.Sprintf("taskSet-%d", i+1), &ecs.TaskSetArgs{
			CapacityProviderStrategies: capacityProviderStrategies,
			Cluster:                    pulumi.String(taskSet.Cluster),
			ExternalId:                 pulumi.StringPtrFromPtr(taskSet.ExternalID),
			ForceDelete:                pulumi.BoolPtrFromPtr(taskSet.ForceDelete),
			LaunchType:                 pulumi.StringPtrFromPtr(taskSet.LaunchType),
			LoadBalancers:              loadBalancers,
			NetworkConfiguration:       networkConfiguration,
			PlatformVersion:            pulumi.StringPtrFromPtr(taskSet.PlatformVersion),
			Scale:                      scale,
			Service:                    pulumi.String(taskSet.Service),
			ServiceRegistries:          serviceRegistries,
			Tags:                       pulumi.ToStringMap(taskSet.Tags),
			TaskDefinition:             pulumi.String(taskSet.TaskDefinition),
			WaitUntilStable:            pulumi.BoolPtrFromPtr(taskSet.WaitUntilStable),
			WaitUntilStableTimeout:     pulumi.StringPtrFromPtr(taskSet.WaitUntilStableTimeout),
		}, pulumi.Parent(component))
		if err != nil {
			return nil, fmt.Errorf("failed to create new task set: %v", err)
		}

		taskSetOutputs = append(taskSetOutputs, &taskSetOutput{
			arn: output.Arn,
		})
	}

	return taskSetOutputs, nil
}
