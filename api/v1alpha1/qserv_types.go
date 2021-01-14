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

package v1alpha1

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
	StorageClass    string `json:"storageClassName,omitempty"`
	StorageCapacity string `json:"storage,omitempty"`

	// Czar defines the settings for czar cluster
	Czar CzarSettings `json:"czar,omitempty"`

	// IngestSettings defines the settings for ingest workflow
	Ingest IngestSettings `json:"ingest,omitempty"`

	// ImagePullPolicy for all containers
	// + kubebuilder:default:=Always
	ImagePullPolicy v1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// Devel defines the settings for development environment
	Devel DevelSettings `json:"devel,omitempty"`

	// QueryService defines the settings for the service which expose Qserv SQL interface (czar/proxy)
	QueryService QueryServiceSettings `json:"queryService,omitempty"`

	// Replication defines the settings for the replication framework
	Replication ReplicationSettings `json:"replication,omitempty"`

	// Tolerations defines the settings for adding custom tolerations to all pods
	// +kubebuilder:validation:Optional
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`

	// Worker defines the settings for worker cluster
	Worker WorkerSettings `json:"worker,omitempty"`

	// Xrootd defines the settings for worker cluster
	Xrootd XrootdSettings `json:"xrootd,omitempty"`

	// NetworkPolicies secures the cluster network using Network Policies.
	// Ensure the Kubernetes cluster has enabled Network plugin.
	NetworkPolicies bool `json:"networkPolicies,omitempty"`
}

// CzarSettings defines the specification of the czar cluster
type CzarSettings struct {
	Affinity v1.Affinity `json:"affinity,omitempty"`
	Image    string      `json:"image,omitempty"`
	// + kubebuilder:default:=1
	Replicas       int32                   `json:"replicas,omitempty"`
	ProxyResources v1.ResourceRequirements `json:"proxyResources,omitempty"`

	// +kubebuilder:validation:Optional
	StorageClass string `json:"storageClassName,omitempty"`
	// +kubebuilder:validation:Optional
	StorageCapacity string `json:"storage,omitempty"`
}

// DevelSettings defines the specification for development/debug environment
type DevelSettings struct {
	CorePath string `json:"corePath,omitempty"`
}

// IngestSettings defines the specification of the ingest workflow
type IngestSettings struct {
	DbImage string `json:"dbImage,omitempty"`
}

// QueryServiceSettings defines the specification of the service which
// expose Qserv czar/proxy port
type QueryServiceSettings struct {
	ServiceType v1.ServiceType    `json:"type,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// ReplicationSettings defines the specification of the replication framework
type ReplicationSettings struct {
	Debug   string `json:"debug,omitempty"`
	DbImage string `json:"dbImage,omitempty"`
	Image   string `json:"image,omitempty"`
}

// WorkerSettings defines the specification of the worker cluster
type WorkerSettings struct {
	Affinity             v1.Affinity             `json:"affinity,omitempty"`
	Image                string                  `json:"image,omitempty"`
	Replicas             int32                   `json:"replicas,omitempty"`
	ReplicationResources v1.ResourceRequirements `json:"replicationResources,omitempty"`
	// +kubebuilder:validation:Optional
	StorageClass string `json:"storageClassName,omitempty"`
	// +kubebuilder:validation:Optional
	StorageCapacity string `json:"storage,omitempty"`
}

// XrootdSettings defines the specification of the xrootd redirectors cluster
type XrootdSettings struct {
	Image    string `json:"image,omitempty"`
	Replicas int32  `json:"replicas,omitempty"`
}

// QservStatus defines the observed state of Qserv
type QservStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Qserv is the Schema for the qservs API
type Qserv struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QservSpec   `json:"spec,omitempty"`
	Status QservStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// QservList contains a list of Qserv
type QservList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Qserv `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Qserv{}, &QservList{})
}
