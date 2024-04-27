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
	Name                     string `json:"name"`
	AutoscalingGroupProvider struct {
		AutoscalingGroupArn string `json:"autoscalingGroupArn"`
		ManagedDraining     string `json:"managedDraining"`
		ManagedScaling      struct {
			InstanceWarmupPeriod   int    `json:"instanceWarmupPeriod"`
			MaximumScalingStepSize int    `json:"maximumScalingStepSize"`
			MinimumScalingStepSize int    `json:"minimumScalingStepSize"`
			Status                 string `json:"status"`
			TargetCapacity         int    `json:"targetCapacity"`
		} `json:"managedScaling"`
		ManagedTerminationProtection string `json:"managedTerminationProtection"`
	} `json:"autoscalingGroupProvider"`
	Tags map[string]string `json:"tags"`
}

// ClusterConfig defines arguments for creating an AWS ECS cluster.
type ClusterConfig struct {
	Name          string `json:"name"`
	Configuration struct {
		ExecuteCommand struct {
			KmsKeyID string `json:"kmsKeyId"`
			Logging  string `json:"logging"`
		} `json:"executeCommand"`
		LogConfiguration struct {
			CloudWatchEncryptionEnabled bool   `json:"cloudWatchEncryptionEnabled"`
			CloudWatchLogGroupName      string `json:"cloudWatchLogGroupName"`
			S3BucketEncryptionEnabled   bool   `json:"s3BucketEncryptionEnabled"`
			S3BucketName                string `json:"s3BucketName"`
			S3KeyPrefix                 string `json:"s3KeyPrefix"`
		} `json:"logConfiguration"`
	} `json:"configuration"`
	Settings []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"settings"`
	ServiceConnectNamespace string            `json:"serviceConnectNamespace"`
	Tags                    map[string]string `json:"tags"`
}

// clusterOutput defines outputs from the AWS ECS cluster creation.
type clusterOutput struct {
	clusterArn pulumi.StringInput
	id         pulumi.StringInput
}

// ClusterCapacityProviderConfig defines arguments for setting up capacity providers for an AWS ECS cluster.
type ClusterCapacityProviderConfig struct {
	ClusterName                       string   `json:"clusterName"`
	CapacityProviders                 []string `json:"capacityProviders"`
	DefaultCapacityProviderStrategies []struct {
		Base             int    `json:"base"`
		CapacityProvider string `json:"capacityProvider"`
		Weight           int    `json:"weight"`
	} `json:"defaultCapacityProviderStrategies"`
}

// ServiceConfig defines arguments for creating an AWS ECS service.
type ServiceConfig struct {
	Name           string `json:"name"`
	ClusterArn     string `json:"clusterArn"`
	TaskDefinition string `json:"taskDefinition"`
	DesiredCount   int    `json:"desiredCount"`
	Alarms         struct {
		AlarmNames []string `json:"alarmNames"`
		Enable     bool     `json:"enable"`
		Rollback   bool     `json:"rollback"`
	} `json:"alarms"`
	CapacityProviderStrategies []struct {
		Name   string `json:"name"`
		Base   int    `json:"base"`
		Weight int    `json:"weight"`
	} `json:"capacityProviderStrategies"`
	IamRole                  string `json:"iamRole"`
	DeploymentCircuitBreaker struct {
		Enable   bool `json:"enable"`
		Rollback bool `json:"rollback"`
	} `json:"deploymentCircuitBreaker"`
	DeploymentController struct {
		Type string `json:"type"`
	} `json:"deploymentController"`
	DeploymentMaximumPercent        int    `json:"deploymentMaximumPercent"`
	DeploymentMinimumHealthyPercent int    `json:"deploymentMinimumHealthyPercent"`
	EnableEcsManagedTags            bool   `json:"enableEcsManagedTags"`
	EnableExecuteCommand            bool   `json:"enableExecuteCommand"`
	ForceNewDeployment              bool   `json:"forceNewDeployment"`
	HealthCheckGracePeriodSeconds   int    `json:"healthCheckGracePeriodSeconds"`
	LaunchType                      string `json:"launchType"`
	LoadBalancers                   []struct {
		ContainerName  string `json:"containerName"`
		ContainerPort  int    `json:"containerPort"`
		ElbName        string `json:"elbName"`
		TargetGroupArn string `json:"targetGroupArn"`
	} `json:"loadBalancers"`
	NetworkConfiguration struct {
		Subnets        []string `json:"subnets"`
		AssignPublicIP bool     `json:"assignPublicIp"`
		SecurityGroups []string `json:"securityGroups"`
	} `json:"networkConfiguration"`
	OrderedPlacementStrategies []struct {
		Type  string `json:"type"`
		Field string `json:"field"`
	} `json:"orderedPlacementStrategies"`
	PlacementConstraints []struct {
		Type       string `json:"type"`
		Expression string `json:"expression"`
	} `json:"placementConstraints"`
	PlatformVersion             string `json:"platformVersion"`
	PropagateTags               string `json:"propagateTags"`
	SchedulingStrategy          string `json:"schedulingStrategy"`
	ServiceConnectConfiguration struct {
		Enabled          bool `json:"enabled"`
		LogConfiguration struct {
			LogDriver     string            `json:"logDriver"`
			Options       map[string]string `json:"options"`
			SecretOptions []struct {
				Name      string `json:"name"`
				ValueFrom string `json:"valueFrom"`
			} `json:"secretOptions"`
		} `json:"logConfiguration"`
		Namespace string `json:"namespace"`
		Services  []struct {
			PortName    string `json:"portName"`
			ClientAlias []struct {
				Port    int    `json:"port"`
				DNSName string `json:"dnsName"`
			} `json:"clientAlias"`
			DiscoveryName       string `json:"discoveryName"`
			IngressPortOverride int    `json:"ingressPortOverride"`
			Timeout             struct {
				IdleTimeoutSeconds       int `json:"idleTimeoutSeconds"`
				PerRequestTimeoutSeconds int `json:"perRequestTimeoutSeconds"`
			} `json:"timeout"`
			TLS struct {
				IssuerCertAuthority struct {
					AwsPcaAuthorityArn string `json:"awsPcaAuthorityArn"`
				} `json:"issuerCertAuthority"`
				KmsKey  string `json:"kmsKey"`
				RoleArn string `json:"roleArn"`
			} `json:"tls"`
		} `json:"services"`
	} `json:"serviceConnectConfiguration"`
	ServiceRegistry struct {
		ContainerName string `json:"containerName"`
		ContainerPort int    `json:"containerPort"`
		Port          int    `json:"port"`
		RegistryArn   string `json:"registryArn"`
	} `json:"serviceRegistry"`
	Tags               map[string]string `json:"tags"`
	Triggers           map[string]string `json:"triggers"`
	WaitForSteadyState bool              `json:"waitForSteadyState"`
}

// serviceOutput defines outputs from the AWS ECS service creation.
type serviceOutput struct {
	id pulumi.StringInput
}

// TaskDefinitionConfig defines arguments for creating an AWS ECS task definition.
type TaskDefinitionConfig struct {
	Name                  string `json:"name"`
	ContainerDefinitions  string `json:"containerDefinitions"`
	NetworkMode           string `json:"networkMode"`
	CPU                   string `json:"cpu"`
	EphemeralStorage      int    `json:"ephemeralStorage"`
	ExecutionRoleArn      string `json:"executionRoleArn"`
	IpcMode               string `json:"ipcMode"`
	InferenceAccelerators []struct {
		DeviceName string `json:"deviceName"`
		DeviceType string `json:"deviceType"`
	} `json:"inferenceAccelerators"`
	Memory               string `json:"memory"`
	PidMode              string `json:"pidMode"`
	PlacementConstraints []struct {
		Expression string `json:"expression"`
		Type       string `json:"type"`
	} `json:"placementConstraints"`
	ProxyConfiguration struct {
		ContainerName string            `json:"containerName"`
		Properties    map[string]string `json:"properties"`
		Type          string            `json:"type"`
	} `json:"proxyConfiguration"`
	RequiresCompatibilities []string `json:"requiresCompatibilities"`
	RuntimePlatform         struct {
		CPUArchitecture       string `json:"cpuArchitecture"`
		OperatingSystemFamily string `json:"operatingSystemFamily"`
	} `json:"runtimePlatform"`
	SkipDestroy bool              `json:"skipDestroy"`
	Tags        map[string]string `json:"tags"`
	TaskRoleArn string            `json:"taskRoleArn"`
	TrackLatest bool              `json:"trackLatest"`
	Volumes     []struct {
		Name                      string `json:"name"`
		DockerVolumeConfiguration struct {
			Autoprovision bool              `json:"autoprovision"`
			Driver        string            `json:"driver"`
			DriverOpts    map[string]string `json:"driverOpts"`
			Labels        map[string]string `json:"labels"`
			Scope         string            `json:"scope"`
		} `json:"dockerVolumeConfiguration"`
		EfsVolumeConfiguration struct {
			FileSystemID        string `json:"fileSystemId"`
			AuthorizationConfig struct {
				AccessPointID string `json:"accessPointId"`
				Iam           string `json:"iam"`
			} `json:"authorizationConfig"`
			RootDirectory         string `json:"rootDirectory"`
			TransitEncryption     string `json:"transitEncryption"`
			TransitEncryptionPort int    `json:"transitEncryptionPort"`
		} `json:"efsVolumeConfiguration"`
		FsxWindowsFileServerVolumeConfiguration struct {
			AuthorizationConfig struct {
				CredentialsParameter string `json:"credentialsParameter"`
				Domain               string `json:"domain"`
			}
			FileSystemID  string `json:"fileSystemId"`
			RootDirectory string `json:"rootDirectory"`
		} `json:"fsxWindowsFileServerVolumeConfiguration"`
		HostPath string `json:"hostPath"`
	} `json:"volumes"`
}

// taskDefinitionOutput defines outputs from the AWS ECS task definition creation.
type taskDefinitionOutput struct {
	arn pulumi.StringOutput
}

// TaskSetConfig defines arguments for creating an AWS ECS task set.
type TaskSetConfig struct {
	Name                       string `json:"name"`
	Service                    string `json:"service"`
	Cluster                    string `json:"cluster"`
	TaskDefinition             string `json:"taskDefinition"`
	CapacityProviderStrategies []struct {
		Base             int    `json:"base"`
		CapacityProvider string `json:"capacityProvider"`
		Weight           int    `json:"weight"`
	} `json:"capacityProviderStrategies"`
	ExternalID    string `json:"externalId"`
	ForceDelete   bool   `json:"forceDelete"`
	LaunchType    string `json:"launchType"`
	LoadBalancers []struct {
		ContainerName    string `json:"containerName"`
		ContainerPort    int    `json:"containerPort"`
		LoadBalancerName string `json:"loadBalancerName"`
		TargetGroupArn   string `json:"targetGroupArn"`
	} `json:"loadBalancers"`
	NetworkConfiguration struct {
		Subnets        []string `json:"subnets"`
		AssignPublicIP bool     `json:"assignPublicIp"`
		SecurityGroups []string `json:"securityGroups"`
	} `json:"networkConfiguration"`
	PlatformVersion string `json:"platformVersion"`
	Scale           struct {
		Unit  string `json:"unit"`
		Value int    `json:"value"`
	} `json:"scale"`
	ServiceRegistries struct {
		RegistryArn   string `json:"registryArn"`
		ContainerName string `json:"containerName"`
		ContainerPort int    `json:"containerPort"`
		Port          int    `json:"port"`
	} `json:"serviceRegistries"`
	Tags                   map[string]string `json:"tags"`
	WaitUntilStable        bool              `json:"waitUntilStable"`
	WaitUntilStableTimeout string            `json:"waitUntilStableTimeout"`
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

		_, err = ecs.NewCapacityProvider(ctx, fmt.Sprintf("capacityProvider-%d", i+1), &ecs.CapacityProviderArgs{
			Name: pulumi.String(capacityProvider.Name),
			AutoScalingGroupProvider: &ecs.CapacityProviderAutoScalingGroupProviderArgs{
				AutoScalingGroupArn: pulumi.String(capacityProvider.AutoscalingGroupProvider.AutoscalingGroupArn),
				ManagedDraining:     pulumi.String(capacityProvider.AutoscalingGroupProvider.ManagedDraining),
				ManagedScaling: &ecs.CapacityProviderAutoScalingGroupProviderManagedScalingArgs{
					InstanceWarmupPeriod:   pulumi.Int(capacityProvider.AutoscalingGroupProvider.ManagedScaling.InstanceWarmupPeriod),
					MaximumScalingStepSize: pulumi.Int(capacityProvider.AutoscalingGroupProvider.ManagedScaling.MaximumScalingStepSize),
					MinimumScalingStepSize: pulumi.Int(capacityProvider.AutoscalingGroupProvider.ManagedScaling.MinimumScalingStepSize),
					Status:                 pulumi.String(capacityProvider.AutoscalingGroupProvider.ManagedScaling.Status),
					TargetCapacity:         pulumi.Int(capacityProvider.AutoscalingGroupProvider.ManagedScaling.TargetCapacity),
				},
				ManagedTerminationProtection: pulumi.String(capacityProvider.AutoscalingGroupProvider.ManagedTerminationProtection),
			},
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

	settings := ecs.ClusterSettingArray{}
	for _, setting := range config.Settings {
		settings = append(settings, &ecs.ClusterSettingArgs{
			Name:  pulumi.String(setting.Name),
			Value: pulumi.String(setting.Value),
		})
	}

	cluster, err := ecs.NewCluster(ctx, "cluster", &ecs.ClusterArgs{
		Name: pulumi.String(config.Name),
		Configuration: &ecs.ClusterConfigurationArgs{
			ExecuteCommandConfiguration: &ecs.ClusterConfigurationExecuteCommandConfigurationArgs{
				KmsKeyId: pulumi.String(config.Configuration.ExecuteCommand.KmsKeyID),
				Logging:  pulumi.String(config.Configuration.ExecuteCommand.Logging),
				LogConfiguration: &ecs.ClusterConfigurationExecuteCommandConfigurationLogConfigurationArgs{
					CloudWatchEncryptionEnabled: pulumi.Bool(config.Configuration.LogConfiguration.CloudWatchEncryptionEnabled),
					CloudWatchLogGroupName:      pulumi.String(config.Configuration.LogConfiguration.CloudWatchLogGroupName),
					S3BucketEncryptionEnabled:   pulumi.Bool(config.Configuration.LogConfiguration.S3BucketEncryptionEnabled),
					S3BucketName:                pulumi.String(config.Configuration.LogConfiguration.S3BucketName),
					S3KeyPrefix:                 pulumi.String(config.Configuration.LogConfiguration.S3KeyPrefix),
				},
			},
		},
		Settings: settings,
		ServiceConnectDefaults: &ecs.ClusterServiceConnectDefaultsArgs{
			Namespace: pulumi.String(config.ServiceConnectNamespace),
		},
		Tags: pulumi.ToStringMap(config.Tags),
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
			Base:             pulumi.Int(defaultCapacityProviderStrategy.Base),
			CapacityProvider: pulumi.String(defaultCapacityProviderStrategy.CapacityProvider),
			Weight:           pulumi.Int(defaultCapacityProviderStrategy.Weight),
		})
	}

	_, err = ecs.NewClusterCapacityProviders(ctx, "clusterCapacityProviders", &ecs.ClusterCapacityProvidersArgs{
		ClusterName:                       pulumi.String(config.ClusterName),
		CapacityProviders:                 pulumi.ToStringArray(config.CapacityProviders),
		DefaultCapacityProviderStrategies: defaultCapacityProviderStrategies,
	}, pulumi.Parent(component))
	if err != nil {
		return fmt.Errorf("failed to create new cluster capacity provider: %v", err)
	}

	return nil
}

// NewService creates a new AWS ECS service.
func NewService(ctx *pulumi.Context, config ServiceConfig, opts ...pulumi.ResourceOption) (*serviceOutput, error) {
	component := &pulumi.ResourceState{}
	err := ctx.RegisterComponentResource("aws:ecs:Service", config.Name, component, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to register component resource: %v", err)
	}

	var capacityProviderStrategies ecs.ServiceCapacityProviderStrategyArray
	for _, capacityProviderStrategy := range config.CapacityProviderStrategies {
		capacityProviderStrategies = append(capacityProviderStrategies, &ecs.ServiceCapacityProviderStrategyArgs{
			CapacityProvider: pulumi.String(capacityProviderStrategy.Name),
			Base:             pulumi.Int(capacityProviderStrategy.Base),
			Weight:           pulumi.Int(capacityProviderStrategy.Weight),
		})
	}

	var loadBalancers ecs.ServiceLoadBalancerArray
	for _, loadBalancer := range config.LoadBalancers {
		loadBalancers = append(loadBalancers, &ecs.ServiceLoadBalancerArgs{
			ContainerName:  pulumi.String(loadBalancer.ContainerName),
			ContainerPort:  pulumi.Int(loadBalancer.ContainerPort),
			ElbName:        pulumi.String(loadBalancer.ElbName),
			TargetGroupArn: pulumi.String(loadBalancer.TargetGroupArn),
		})
	}

	var orderedPlacementStrategies ecs.ServiceOrderedPlacementStrategyArray
	for _, orderedPlacementStrategy := range config.OrderedPlacementStrategies {
		orderedPlacementStrategies = append(orderedPlacementStrategies, &ecs.ServiceOrderedPlacementStrategyArgs{
			Type:  pulumi.String(orderedPlacementStrategy.Type),
			Field: pulumi.String(orderedPlacementStrategy.Field),
		})
	}

	var placementConstraints ecs.ServicePlacementConstraintArray
	for _, placementConstraint := range config.PlacementConstraints {
		placementConstraints = append(placementConstraints, &ecs.ServicePlacementConstraintArgs{
			Type:       pulumi.String(placementConstraint.Type),
			Expression: pulumi.String(placementConstraint.Expression),
		})
	}

	var serviceConnectLogsSecretOptions ecs.ServiceServiceConnectConfigurationLogConfigurationSecretOptionArray
	for _, serviceConnectLogsSecretOption := range config.ServiceConnectConfiguration.LogConfiguration.SecretOptions {
		serviceConnectLogsSecretOptions = append(serviceConnectLogsSecretOptions, &ecs.ServiceServiceConnectConfigurationLogConfigurationSecretOptionArgs{
			Name:      pulumi.String(serviceConnectLogsSecretOption.Name),
			ValueFrom: pulumi.String(serviceConnectLogsSecretOption.ValueFrom),
		})
	}

	var serviceConnectServicesClientAliases ecs.ServiceServiceConnectConfigurationServiceClientAliasArray
	for _, service := range config.ServiceConnectConfiguration.Services {
		for _, clientAlias := range service.ClientAlias {
			serviceConnectServicesClientAliases = append(serviceConnectServicesClientAliases, &ecs.ServiceServiceConnectConfigurationServiceClientAliasArgs{
				Port:    pulumi.Int(clientAlias.Port),
				DnsName: pulumi.String(clientAlias.DNSName),
			})
		}
	}

	var serviceConnectServices ecs.ServiceServiceConnectConfigurationServiceArray
	for _, serviceConnectService := range config.ServiceConnectConfiguration.Services {
		serviceConnectServices = append(serviceConnectServices, &ecs.ServiceServiceConnectConfigurationServiceArgs{
			PortName:            pulumi.String(serviceConnectService.PortName),
			ClientAlias:         serviceConnectServicesClientAliases,
			DiscoveryName:       pulumi.String(serviceConnectService.DiscoveryName),
			IngressPortOverride: pulumi.Int(serviceConnectService.IngressPortOverride),
			Timeout: &ecs.ServiceServiceConnectConfigurationServiceTimeoutArgs{
				IdleTimeoutSeconds:       pulumi.Int(serviceConnectService.Timeout.IdleTimeoutSeconds),
				PerRequestTimeoutSeconds: pulumi.Int(serviceConnectService.Timeout.PerRequestTimeoutSeconds),
			},
			Tls: &ecs.ServiceServiceConnectConfigurationServiceTlsArgs{
				IssuerCertAuthority: &ecs.ServiceServiceConnectConfigurationServiceTlsIssuerCertAuthorityArgs{
					AwsPcaAuthorityArn: pulumi.String(serviceConnectService.TLS.IssuerCertAuthority.AwsPcaAuthorityArn),
				},
				KmsKey:  pulumi.String(serviceConnectService.TLS.KmsKey),
				RoleArn: pulumi.String(serviceConnectService.TLS.RoleArn),
			},
		})
	}

	service, err := ecs.NewService(ctx, "service", &ecs.ServiceArgs{
		Name:           pulumi.String(config.Name),
		Cluster:        pulumi.String(config.ClusterArn),
		TaskDefinition: pulumi.String(config.TaskDefinition),
		DesiredCount:   pulumi.Int(config.DesiredCount),
		Alarms: &ecs.ServiceAlarmsArgs{
			AlarmNames: pulumi.ToStringArray(config.Alarms.AlarmNames),
			Enable:     pulumi.Bool(config.Alarms.Enable),
			Rollback:   pulumi.Bool(config.Alarms.Rollback),
		},
		CapacityProviderStrategies: capacityProviderStrategies,
		DeploymentCircuitBreaker: &ecs.ServiceDeploymentCircuitBreakerArgs{
			Enable:   pulumi.Bool(config.DeploymentCircuitBreaker.Enable),
			Rollback: pulumi.Bool(config.DeploymentCircuitBreaker.Rollback),
		},
		DeploymentController: &ecs.ServiceDeploymentControllerArgs{
			Type: pulumi.String(config.DeploymentController.Type),
		},
		DeploymentMaximumPercent:        pulumi.Int(config.DeploymentMaximumPercent),
		DeploymentMinimumHealthyPercent: pulumi.Int(config.DeploymentMinimumHealthyPercent),
		EnableEcsManagedTags:            pulumi.Bool(config.EnableEcsManagedTags),
		EnableExecuteCommand:            pulumi.Bool(config.EnableExecuteCommand),
		ForceNewDeployment:              pulumi.Bool(config.ForceNewDeployment),
		HealthCheckGracePeriodSeconds:   pulumi.Int(config.HealthCheckGracePeriodSeconds),
		IamRole:                         pulumi.String(config.IamRole),
		LaunchType:                      pulumi.String(config.LaunchType),
		LoadBalancers:                   loadBalancers,
		NetworkConfiguration: &ecs.ServiceNetworkConfigurationArgs{
			Subnets:        pulumi.ToStringArray(config.NetworkConfiguration.Subnets),
			AssignPublicIp: pulumi.Bool(config.NetworkConfiguration.AssignPublicIP),
			SecurityGroups: pulumi.ToStringArray(config.NetworkConfiguration.SecurityGroups),
		},
		OrderedPlacementStrategies: orderedPlacementStrategies,
		PlacementConstraints:       placementConstraints,
		PlatformVersion:            pulumi.String(config.PlatformVersion),
		PropagateTags:              pulumi.String(config.PropagateTags),
		SchedulingStrategy:         pulumi.String(config.SchedulingStrategy),
		ServiceConnectConfiguration: &ecs.ServiceServiceConnectConfigurationArgs{
			Enabled: pulumi.Bool(config.ServiceConnectConfiguration.Enabled),
			LogConfiguration: &ecs.ServiceServiceConnectConfigurationLogConfigurationArgs{
				LogDriver:     pulumi.String(config.ServiceConnectConfiguration.LogConfiguration.LogDriver),
				Options:       pulumi.ToStringMap(config.ServiceConnectConfiguration.LogConfiguration.Options),
				SecretOptions: serviceConnectLogsSecretOptions,
			},
			Namespace: pulumi.String(config.ServiceConnectConfiguration.Namespace),
			Services:  serviceConnectServices,
		},
		ServiceRegistries: &ecs.ServiceServiceRegistriesArgs{
			ContainerName: pulumi.String(config.ServiceRegistry.ContainerName),
			ContainerPort: pulumi.Int(config.ServiceRegistry.ContainerPort),
			Port:          pulumi.Int(config.ServiceRegistry.Port),
			RegistryArn:   pulumi.String(config.ServiceRegistry.RegistryArn),
		},
		Tags:               pulumi.ToStringMap(config.Tags),
		Triggers:           pulumi.ToStringMap(config.Triggers),
		WaitForSteadyState: pulumi.Bool(config.WaitForSteadyState),
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
			Expression: pulumi.String(placementConstraint.Expression),
			Type:       pulumi.String(placementConstraint.Type),
		})
	}

	var volumes ecs.TaskDefinitionVolumeArray
	for _, volume := range config.Volumes {
		volumes = append(volumes, &ecs.TaskDefinitionVolumeArgs{
			Name: pulumi.String(volume.Name),
			DockerVolumeConfiguration: &ecs.TaskDefinitionVolumeDockerVolumeConfigurationArgs{
				Autoprovision: pulumi.Bool(volume.DockerVolumeConfiguration.Autoprovision),
				Driver:        pulumi.String(volume.DockerVolumeConfiguration.Driver),
				DriverOpts:    pulumi.ToStringMap(volume.DockerVolumeConfiguration.DriverOpts),
				Labels:        pulumi.ToStringMap(volume.DockerVolumeConfiguration.Labels),
				Scope:         pulumi.String(volume.DockerVolumeConfiguration.Scope),
			},
			EfsVolumeConfiguration: &ecs.TaskDefinitionVolumeEfsVolumeConfigurationArgs{
				FileSystemId: pulumi.String(volume.EfsVolumeConfiguration.FileSystemID),
				AuthorizationConfig: &ecs.TaskDefinitionVolumeEfsVolumeConfigurationAuthorizationConfigArgs{
					AccessPointId: pulumi.String(volume.EfsVolumeConfiguration.AuthorizationConfig.AccessPointID),
					Iam:           pulumi.String(volume.EfsVolumeConfiguration.AuthorizationConfig.Iam),
				},
				RootDirectory:         pulumi.String(volume.EfsVolumeConfiguration.RootDirectory),
				TransitEncryption:     pulumi.String(volume.EfsVolumeConfiguration.TransitEncryption),
				TransitEncryptionPort: pulumi.Int(volume.EfsVolumeConfiguration.TransitEncryptionPort),
			},
			FsxWindowsFileServerVolumeConfiguration: &ecs.TaskDefinitionVolumeFsxWindowsFileServerVolumeConfigurationArgs{
				AuthorizationConfig: &ecs.TaskDefinitionVolumeFsxWindowsFileServerVolumeConfigurationAuthorizationConfigArgs{
					CredentialsParameter: pulumi.String(volume.FsxWindowsFileServerVolumeConfiguration.AuthorizationConfig.CredentialsParameter),
					Domain:               pulumi.String(volume.FsxWindowsFileServerVolumeConfiguration.AuthorizationConfig.Domain),
				},
				FileSystemId:  pulumi.String(volume.FsxWindowsFileServerVolumeConfiguration.FileSystemID),
				RootDirectory: pulumi.String(volume.FsxWindowsFileServerVolumeConfiguration.RootDirectory),
			},
			HostPath: pulumi.String(volume.HostPath),
		})
	}

	taskDefinition, err := ecs.NewTaskDefinition(ctx, "taskDefinition", &ecs.TaskDefinitionArgs{
		Family:               pulumi.String(config.Name),
		ContainerDefinitions: pulumi.String(config.ContainerDefinitions),
		NetworkMode:          pulumi.String(config.NetworkMode),
		Cpu:                  pulumi.String(config.CPU),
		EphemeralStorage: &ecs.TaskDefinitionEphemeralStorageArgs{
			SizeInGib: pulumi.Int(config.EphemeralStorage),
		},
		ExecutionRoleArn:      pulumi.String(config.ExecutionRoleArn),
		IpcMode:               pulumi.String(config.IpcMode),
		InferenceAccelerators: inferenceAccelerators,
		Memory:                pulumi.String(config.Memory),
		PidMode:               pulumi.String(config.PidMode),
		PlacementConstraints:  placementConstraints,
		ProxyConfiguration: &ecs.TaskDefinitionProxyConfigurationArgs{
			ContainerName: pulumi.String(config.ProxyConfiguration.ContainerName),
			Properties:    pulumi.ToStringMap(config.ProxyConfiguration.Properties),
			Type:          pulumi.String(config.ProxyConfiguration.Type),
		},
		RequiresCompatibilities: pulumi.ToStringArray(config.RequiresCompatibilities),
		RuntimePlatform: &ecs.TaskDefinitionRuntimePlatformArgs{
			CpuArchitecture:       pulumi.String(config.RuntimePlatform.CPUArchitecture),
			OperatingSystemFamily: pulumi.String(config.RuntimePlatform.OperatingSystemFamily),
		},
		SkipDestroy: pulumi.Bool(config.SkipDestroy),
		Tags:        pulumi.ToStringMap(config.Tags),
		TaskRoleArn: pulumi.String(config.TaskRoleArn),
		TrackLatest: pulumi.Bool(config.TrackLatest),
		Volumes:     volumes,
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

		var loadBalancers ecs.TaskSetLoadBalancerArray
		for _, loadBalancer := range taskSet.LoadBalancers {
			loadBalancers = append(loadBalancers, &ecs.TaskSetLoadBalancerArgs{
				ContainerName:    pulumi.String(loadBalancer.ContainerName),
				ContainerPort:    pulumi.Int(loadBalancer.ContainerPort),
				LoadBalancerName: pulumi.String(loadBalancer.LoadBalancerName),
				TargetGroupArn:   pulumi.String(loadBalancer.TargetGroupArn),
			})
		}

		var capacityProviderStrategies ecs.TaskSetCapacityProviderStrategyArray
		for _, capacityProviderStrategy := range taskSet.CapacityProviderStrategies {
			capacityProviderStrategies = append(capacityProviderStrategies, &ecs.TaskSetCapacityProviderStrategyArgs{
				Base:             pulumi.Int(capacityProviderStrategy.Base),
				CapacityProvider: pulumi.String(capacityProviderStrategy.CapacityProvider),
				Weight:           pulumi.Int(capacityProviderStrategy.Weight),
			})
		}

		output, err := ecs.NewTaskSet(ctx, fmt.Sprintf("taskSet-%d", i+1), &ecs.TaskSetArgs{
			Service:                    pulumi.String(taskSet.Service),
			Cluster:                    pulumi.String(taskSet.Cluster),
			TaskDefinition:             pulumi.String(taskSet.TaskDefinition),
			CapacityProviderStrategies: capacityProviderStrategies,
			ExternalId:                 pulumi.String(taskSet.ExternalID),
			ForceDelete:                pulumi.Bool(taskSet.ForceDelete),
			LaunchType:                 pulumi.String(taskSet.LaunchType),
			LoadBalancers:              loadBalancers,
			NetworkConfiguration: &ecs.TaskSetNetworkConfigurationArgs{
				Subnets:        pulumi.ToStringArray(taskSet.NetworkConfiguration.Subnets),
				AssignPublicIp: pulumi.Bool(taskSet.NetworkConfiguration.AssignPublicIP),
				SecurityGroups: pulumi.ToStringArray(taskSet.NetworkConfiguration.SecurityGroups),
			},
			PlatformVersion: pulumi.String(taskSet.PlatformVersion),
			Scale: &ecs.TaskSetScaleArgs{
				Unit:  pulumi.String(taskSet.Scale.Unit),
				Value: pulumi.Float64(taskSet.Scale.Value),
			},

			ServiceRegistries: &ecs.TaskSetServiceRegistriesArgs{
				RegistryArn:   pulumi.String(taskSet.ServiceRegistries.RegistryArn),
				ContainerName: pulumi.String(taskSet.ServiceRegistries.ContainerName),
				ContainerPort: pulumi.Int(taskSet.ServiceRegistries.ContainerPort),
				Port:          pulumi.Int(taskSet.ServiceRegistries.Port),
			},
			Tags:                   pulumi.ToStringMap(taskSet.Tags),
			WaitUntilStable:        pulumi.Bool(taskSet.WaitUntilStable),
			WaitUntilStableTimeout: pulumi.String(taskSet.WaitUntilStableTimeout),
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
