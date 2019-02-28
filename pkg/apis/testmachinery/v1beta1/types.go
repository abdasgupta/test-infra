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

package v1beta1

import (
	argov1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

////////////////////////////////////////////////////
//                    TESTRUN                     //
////////////////////////////////////////////////////

// LocationType is the type of a TestDef location.
type LocationType string

// Testdefinition location types
const (
	LocationTypeGit   LocationType = "git"
	LocationTypeLocal LocationType = "local"
)

// ConditionType is the type of a testflow step indicating when it should be executed.
type ConditionType string

// Testflow step conditions
const (
	ConditionTypeError   ConditionType = "error"
	ConditionTypeSuccess ConditionType = "success"
	ConditionTypeAlways  ConditionType = "always"
)

// ConfigType is the type of a ConfigElement.
type ConfigType string

// ConfigElement types
const (
	ConfigTypeEnv  = "env"
	ConfigTypeFile = "file"
)

// Testrun statuses
const (
	PhaseStatusInit    argov1.NodePhase = "init"
	PhaseStatusPending                  = argov1.NodePending
	PhaseStatusRunning                  = argov1.NodeRunning
	PhaseStatusSuccess                  = argov1.NodeSucceeded
	PhaseStatusFailed                   = argov1.NodeFailed
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Testrun is the description of the testflow that should be executed.
// +k8s:openapi-gen=true
type Testrun struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TestrunSpec   `json:"spec"`
	Status TestrunStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TestrunList contains a list of Testruns
type TestrunList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Testrun `json:"items"`
}

// TestrunSpec is the specification of a Testrun.
type TestrunSpec struct {
	Creator string `json:"creator,omitempty"`

	TTLSecondsAfterFinished *int32 `json:"ttlSecondsAfterFinished,omitempty"`

	// TestLocation define repositories to look for TestDefinitions that then executed in a workkflow as specified in testflow.
	TestLocations []TestLocation `json:"testLocations,omitempty"`

	// Base64 encoded kubeconfigs that are mounted to every testflow step.
	// They are available at $TM_KUBECONFIG_PATH/xxx.config, where xxx is either (gardener, seed or shoot).
	// +optional
	Kubeconfigs TestrunKubeconfigs `json:"kubeconfigs,omitempty"`

	// Global config which is available to all tests in the testflow and onExit flow.
	// +optional
	Config []ConfigElement `json:"config,omitempty"`

	TestFlow TestFlow `json:"testFlow,omitempty"`

	// OnExit flow is called when the test is completed.
	// +optional
	OnExit TestFlow `json:"onExit,omitempty"`
}

// TestrunStatus is the status of the Testrun.
type TestrunStatus struct {
	// Phase is the summary of all executed steps.
	Phase argov1.NodePhase `json:"phase,omitempty"`

	// State is a string that represents the actual state and/or process of the testrun.
	State string `json:"state,omitempty"`

	// StartTime is the time when the argo workflow starts executing the steps.
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// CompletionTime is the time when the argo workflow is completed.
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// Duration represents the overall duration of the argo workflow.
	// This value is calculated by (CompletionTime - StartTime)
	Duration int64 `json:"duration,omitempty"`

	// Workflow is the name of the generated argo workflow
	Workflow string `json:"workflow,omitempty"`

	// Steps is the detailed summary of every step.
	// It also shows all specific executed tests.
	Steps [][]*TestflowStepStatus `json:"steps,omitempty"`
}

// TestflowStepStatus is the status of Testflow step
type TestflowStepStatus struct {
	TestDefinition    TestflowStepStatusTestDefinition `json:"testdefinition,omitempty"`
	Phase             argov1.NodePhase                 `json:"phase,omitempty"`
	StartTime         *metav1.Time                     `json:"startTime,omitempty"`
	CompletionTime    *metav1.Time                     `json:"completionTime,omitempty"`
	Duration          int64                            `json:"duration,omitempty"`
	ExportArtifactKey string                           `json:"exportArtifactKey"`
}

// TestflowStepStatusTestDefinition holds information about the used testdefinition and its location.
type TestflowStepStatusTestDefinition struct {
	Name                string            `json:"name,omitempty"`
	Location            TestLocation      `json:"location,omitempty"`
	Owner               string            `json:"owner,omitempty"`
	RecipientsOnFailure []string          `json:"recipientsOnFailure"`
	Position            map[string]string `json:"position"`
}

// TestLocation describes a location to search for TestDefinitions.
type TestLocation struct {
	Type LocationType `json:"type"`
	// Only for LocationType git
	// +optional
	Repo string `json:"repo,omitempty"`
	// +optional
	Revision string `json:"revision,omitempty"`
	// The absolute host on the minikube VM.
	// Only for local
	// +optional
	HostPath string `json:"hostPath,omitempty"`
}

// TestrunKubeconfigs are parameters where Shoot, Seed or a Gardener kubeconfig for the Testrun can be specified.
type TestrunKubeconfigs struct {
	Gardener string `json:"gardener,omitempty"`
	Seed     string `json:"seed,omitempty"`
	Shoot    string `json:"shoot,omitempty"`
}

// ConfigElement is a parameter of a certain type which is passed to TestDefinitions.
type ConfigElement struct {
	// Type of the config value. For now only environament varibales are supported.
	Type ConfigType `json:"type"`

	// Name of the environment variable. Must be a C_IDENTIFIER.
	Name string `json:"name"`

	// value of the environament variable.
	// +opional
	Value string `json:"value"`

	// Fetches the value from a secret or configmap on the testmachinery cluster.
	// +optional
	ValueFrom *ConfigSource `json:"valueFrom,omitempty"`

	// Only for type=file. Path where the file should be mounted.
	// +optional
	Path string `json:"path"`
}

// ConfigSource represents a source for the value of a config element.
type ConfigSource struct {
	// Selects a key of a ConfigMap.
	// +optional
	ConfigMapKeyRef *corev1.ConfigMapKeySelector `json:"configMapKeyRef,omitempty"`
	// Selects a key of a secret in the pod's namespace
	// +optional
	SecretKeyRef *corev1.SecretKeySelector `json:"secretKeyRef,omitempty"`
}

// TestFlow is a 2 dimensional array of testflowsteps which define the execution order of TestDefinitions.
type TestFlow [][]TestflowStep

// TestflowStep is a reference to one or more TestDefinitions to execute in a series of steps.TestflowStep
// TestDefinitions can be either defined by a Name or a Label.
type TestflowStep struct {
	Name      string          `json:"name,omitempty"`
	Label     string          `json:"label,omitempty"`
	Condition ConditionType   `json:"condition,omitempty"`
	Config    []ConfigElement `json:"config,omitempty"`
}

////////////////////////////////////////////////////
//                TESTDEFINITION                  //
////////////////////////////////////////////////////

// TestDefinitionName is the kind identifier of a testdefinition.
const TestDefinitionName = "TestDefinition"

// TestDefinition describes the execution of a test.
type TestDefinition struct {
	Kind     string          `json:"kind"`
	Metadata TestDefMetadata `json:"metadata,omitempty"`

	Spec TestDefSpec `json:"spec"`
}

// TestDefMetadata holds the metadata of a testrun.
type TestDefMetadata struct {
	Name string `json:"name"`
}

// TestDefSpec is the actual description of the test.
type TestDefSpec struct {
	Owner                 string          `json:"owner"`
	RecipientsOnFailure   []string        `json:"recipientsOnFailure"`
	Description           string          `json:"description"`
	Labels                []string        `json:"labels"`
	Behavior              []string        `json:"behavior"`
	ActiveDeadlineSeconds *int64          `json:"activeDeadlineSeconds"`
	Command               []string        `json:"command"`
	Args                  []string        `json:"args"`
	Image                 string          `json:"image"`
	Config                []ConfigElement `json:"config,omitempty"`
}
