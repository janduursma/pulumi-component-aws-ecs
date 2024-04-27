# Pulumi Component AWS ECS

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

## Overview

This Go package is designed to simplify the creation of AWS ECS (Elastic Container Service) resources using Pulumi. It provides an intuitive API to define ECS services, task definitions, clusters, and associated resources, abstracting away the complexity of the underlying infrastructure setup.

## Features

- **Abstraction Layer**: Easily create and manage AWS ECS resources without directly dealing with Pulumi details.
- **Concise API**: Intuitive API design for defining ECS services, task definitions, and clusters with minimal boilerplate.
- **Documentation**: Well-documented codebase with inline comments and comprehensive API documentation generated with godoc.

## Installation

To use this package in your Go project, simply import it as a dependency:

```bash
go get github.com/janduursma/pulumi-component-aws-ecs
```

## Usage

### Example

Here's a basic example of how to use this package to create an AWS ECS service:

```go
package main

import (
	"fmt"
	"github.com/janduursma/pulumi-component-aws-ecs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	// Define your ECS service
	serviceArgs := ecs.ServiceConfig{
		Name:           "test-service",
		ClusterArn:     "<cluster-arn>",
		TaskDefinition: "<task-definition>",
		DesiredCount:   2,
		// Provide other arguments here
	}

	// Create the ECS service
	var ctx *pulumi.Context
	_, err := ecs.NewService(ctx, serviceArgs)
	if err != nil {
		fmt.Println(err)
	}
}
```

## Additional Resources

- [Pulumi AWS ECS Documentation](https://www.pulumi.com/docs/reference/pkg/aws/ecs/)
- [Pulumi AWS SDK Documentation](https://github.com/pulumi/pulumi-aws)

## License

This project is licensed under the [MIT License](LICENSE).