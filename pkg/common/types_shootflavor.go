// Copyright 2019 Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
)

// ExtendedShootFlavors contains a list of extended shoot flavors
type ExtendedShootFlavors struct {
	Flavors []*ExtendedShootFlavor `json:"flavors"`
}

// ExtendedShoot is one instance that is generated from a extended shoot flavor
type ExtendedShoot struct {
	Shoot
	ExtendedShootConfiguration
}

// ExtendedShootFlavor is the shoot flavor with extended configuration
type ExtendedShootFlavor struct {
	ShootFlavor
	ExtendedConfiguration
}

// Shoot is the internal representation of one instance that is generated from a shoot flavor
type Shoot struct {
	// Short description of the flavor
	// +optional
	Description string

	// Cloudprovider of the shoot
	Provider CloudProvider

	// Kubernetes versions to test
	KubernetesVersion gardencorev1beta1.ExpirableVersion

	// AllowPrivilegedContainers defines whether privileged containers will be allowed in the given shoot or not
	AllowPrivilegedContainers *bool

	// Worker pools to test
	Workers []gardencorev1beta1.Worker
}

// ShootFlavor describes the shoot flavors that should be tested.
type ShootFlavor struct {
	// Short description of the flavor
	// +optional
	Description string `json:"description"`

	// Cloudprovider of the shoot
	Provider CloudProvider `json:"provider"`

	// Kubernetes versions to test
	KubernetesVersions ShootKubernetesVersionFlavor `json:"kubernetes"`

	// AllowPrivilegedContainers defines whether privileged containers will be allowed in the given shoot or not
	// +optional
	AllowPrivilegedContainers *bool `json:"allowPrivilegedContainers"`

	// Worker pools to test
	Workers []ShootWorkerFlavor `json:"workers"`
}

type ShootKubernetesVersionFlavor struct {
	// Regex to select versions from the cloudprofile
	// +optional
	Pattern *string `json:"pattern"`

	// FilterPatchVersions will only keep the latest patch version of all minor versions
	// +optional
	FilterPatchVersions *bool `json:"filterPatchVersions"`

	// List of versions to test
	// +optional
	Versions *[]gardencorev1beta1.ExpirableVersion `json:"versions"`
}

// ShootWorkerFlavor defines the worker pools that should be tested
type ShootWorkerFlavor struct {
	WorkerPools []gardencorev1beta1.Worker `json:"workerPools"`
}

// ExtendedConfiguration specifies extended configuration for shoot flavors that are deployed into a preexisting landscape
type ExtendedConfiguration struct {
	ProjectName      string `json:"projectName"`
	CloudprofileName string `json:"cloudprofile"`
	SecretBinding    string `json:"secretBinding"`
	Region           string `json:"region"`
	Zone             string `json:"zone"`

	FloatingPoolName     string `json:"floatingPoolName"`
	LoadbalancerProvider string `json:"loadbalancerProvider"`

	// ControlPlaneConfig contains the provider-specific control plane config blob.
	// Overwrites the controlplane config generated by the specific provider if defined.
	ControlPlaneConfig *gardencorev1beta1.ProviderConfig `json:"controlPlaneConfig,omitempty"`
	// InfrastructureConfig contains the provider-specific infrastructure config blob.
	// Overwrites the infrastructure config generated by the specific provider if defined.
	InfrastructureConfig *gardencorev1beta1.ProviderConfig `json:"infrastructureConfig,omitempty"`
	// NetworkingConfig contains the provider-specific infrastructure config blob.
	// Overwrites the networking config generated by the specific provider if defined.
	NetworkingConfig *gardencorev1beta1.ProviderConfig `json:"networkingConfig,omitempty"`

	// ChartPath use the specific chartPath to render the flavor.
	// This will overwrite the default shoot chart path.
	ChartPath *string `json:"chartPath,omitempty"`
}

// ExtendedShootConfiguration specifies extended configuration for shoots that are deployed into a preexisting landscape
type ExtendedShootConfiguration struct {
	Name         string                         `json:"name"`
	Namespace    string                         `json:"namespace"`
	Cloudprofile gardencorev1beta1.CloudProfile `json:"-"`
	ExtendedConfiguration
}
