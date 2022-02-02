/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package v1alpha1 contains Qserv custom resource definition generator
// +k8s:deepcopy-gen=package
// +kubebuilder:validation:Optional
package v1beta1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// QservSpec defines the desired state of Qserv
type QservSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Common settings
	// +kubebuilder:default:="standard"
	StorageClass string `json:"storageClassName,omitempty"`
	// +kubebuilder:default:="10Gi"
	StorageCapacity string `json:"storage,omitempty"`

	// Czar defines the settings for czar cluster
	Czar CzarSettings `json:"czar,omitempty"`

	// IngestSettings defines the settings for ingest workflow
	Ingest IngestSettings `json:"ingest,omitempty"`

	// +kubebuilder:validation:Required
	DbImage string `json:"dbImage,omitempty"`

	// +kubebuilder:validation:Required
	Image string `json:"image,omitempty"`

	// ImagePullPolicy for all containers
	// +kubebuilder:default:="Always"
	ImagePullPolicy v1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// Devel defines the settings for development environment
	Devel DevelSettings `json:"devel,omitempty"`

	// QueryService defines the settings for the service which expose Qserv SQL interface (czar/proxy)
	QueryService QueryServiceSettings `json:"queryService,omitempty"`

	// Replication defines the settings for the replication framework
	Replication ReplicationSettings `json:"replication,omitempty"`

	// Tolerations defines the settings for adding custom tolerations to all pods
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`

	// Worker defines the settings for worker cluster
	Worker WorkerSettings `json:"worker,omitempty"`

	// Xrootd defines the settings for worker cluster
	Xrootd XrootdSettings `json:"xrootd,omitempty"`

	// NetworkPolicies secures the cluster network using Network Policies.
	// Ensure the Kubernetes cluster has enabled Network plugin.
	// +kubebuilder:default:=false
	NetworkPolicies bool `json:"networkPolicies,omitempty"`
}

// CzarSettings defines the specification of the czar cluster
type CzarSettings struct {
	Affinity v1.Affinity `json:"affinity,omitempty"`

	// +kubebuilder:default:=1
	Replicas       int32                   `json:"replicas,omitempty"`
	ProxyResources v1.ResourceRequirements `json:"proxyResources,omitempty"`

	StorageClass    string `json:"storageClassName,omitempty"`
	StorageCapacity string `json:"storage,omitempty"`
}

// DevelSettings defines the specification for development/debug environment
type DevelSettings struct {
	// +kubebuilder:default:="/tmp/coredump"
	CorePath string `json:"corePath,omitempty"`
	// +kubebuilder:validation:Required
	DebuggerImage string `json:"debuggerImage"`

	// EnableDebugger allows to share process namespace between containers in a Pod
	// and adds a debug container to the  pod
	// See https://kubernetes.io/docs/tasks/configure-pod-container/share-process-namespace/
	// +kubebuilder:default:=false
	EnableDebugger bool `json:"enableDebugger,omitempty"`
}

// IngestSettings defines the specification of the ingest workflow
type IngestSettings struct {
	Affinity v1.Affinity `json:"affinity,omitempty"`
}

// QueryServiceSettings defines the specification of the service which
// expose Qserv czar/proxy port
type QueryServiceSettings struct {
	Annotations    map[string]string `json:"annotations,omitempty"`
	LoadBalancerIP string            `json:"loadBalancerIP,omitempty"`
	NodePort       int32             `json:"nodePort,omitempty"`
	ServiceType    v1.ServiceType    `json:"type,omitempty"`
}

// ReplicationSettings defines the specification of the replication framework
type ReplicationSettings struct {
	Affinity v1.Affinity `json:"affinity,omitempty"`
	Debug    string      `json:"debug,omitempty"`
}

// WorkerSettings defines the specification of the worker cluster
type WorkerSettings struct {
	Affinity v1.Affinity `json:"affinity,omitempty"`
	// +kubebuilder:default:=2
	Replicas             int32                   `json:"replicas,omitempty"`
	ReplicationResources v1.ResourceRequirements `json:"replicationResources,omitempty"`

	StorageClass string `json:"storageClassName,omitempty"`
	// +kubebuilder:validation:Optional
	StorageCapacity string `json:"storage,omitempty"`
}

// XrootdSettings defines the specification of the xrootd redirectors cluster
type XrootdSettings struct {
	Affinity v1.Affinity `json:"affinity,omitempty"`
	// +kubebuilder:default:=2
	Replicas int32 `json:"replicas,omitempty"`
}

// QservStatus defines the observed state of Qserv
type QservStatus struct {
	// Important: Run "make" to regenerate code after modifying this file

	// Known .status.conditions.type are: "Available", "Progressing", and "Degraded"
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions                         []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
	CzarReadyFraction                  string             `json:"czarreadyfraction,omitempty"`
	IngestDatabaseReadyFraction        string             `json:"ingestdatabasereadyfraction,omitempty"`
	ReplicationControllerReadyFraction string             `json:"replicationcontrollerreadyfraction,omitempty"`
	ReplicationDatabaseReadyFraction   string             `json:"replicationdatabasereadyfraction,omitempty"`
	WorkerReadyFraction                string             `json:"workerreadyfraction,omitempty"`
	XrootdReadyFraction                string             `json:"xrootdreadyfraction,omitempty"`
}

// Qserv is the Schema for the qservs API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Czar",type=string,JSONPath=`.status.czarreadyfraction`
// +kubebuilder:printcolumn:name="Ingest-db",type=string,JSONPath=`.status.ingestdatabasereadyfraction`
// +kubebuilder:printcolumn:name="Repl-ctl",type=string,JSONPath=`.status.replicationcontrollerreadyfraction`
// +kubebuilder:printcolumn:name="Repl-db",type=string,JSONPath=`.status.replicationdatabasereadyfraction`
// +kubebuilder:printcolumn:name="Worker",type=string,JSONPath=`.status.workerreadyfraction`
// +kubebuilder:printcolumn:name="Xrootd",type=string,JSONPath=`.status.xrootdreadyfraction`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Qserv struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QservSpec   `json:"spec,omitempty"`
	Status QservStatus `json:"status,omitempty"`
}

// QservList contains a list of Qserv
// +kubebuilder:object:root=true
type QservList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Qserv `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Qserv{}, &QservList{})
}
