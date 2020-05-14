// +kubebuilder:validation:Required
package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// QservSpec defines the desired state of Qserv
// +k8s:openapi-gen=true
// +kubebuilder:validation:Required
type QservSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	// Common settings
	StorageClass    string `json:"storageclass,omitempty"`
	StorageCapacity string `json:"storagecapacity,omitempty"`

	// Czar defines the settings for czar cluster
	Czar CzarSettings `json:"czar,omitempty"`

	// Redis defines the settings for redis cluster
	// +kubebuilder:validation:Optional
	Redis *RedisSettings `json:"redis,omitempty"`

	// Replication defines the settings for the replication framework
	Replication ReplicationSettings `json:"replication,omitempty"`

	// Tolerations defines the settings for adding custom tolerations to all pods
	// +kubebuilder:validation:Optional
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`

	// Worker defines the settings for worker cluster
	Worker WorkerSettings `json:"worker,omitempty"`

	// Xrootd defines the settings for worker cluster
	Xrootd XrootdSettings `json:"xrootd,omitempty"`

	// SecureNetwork secures the cluster network using Network Policies.
	// Ensure the Kubernetes cluster has enabled Network plugin.
	// +kubebuilder:default:=false
	SecureNetwork bool `json:"securenetwork,omitempty"`
}

// CzarSettings defines the specification of the czar cluster
type CzarSettings struct {
	Image string `json:"image,omitempty"`
	// + kubebuilder:default:=1
	Replicas int32 `json:"replicas,omitempty"`
}

// WorkerSettings defines the specification of the worker cluster
type WorkerSettings struct {
	Image    string `json:"image,omitempty"`
	Replicas int32  `json:"replicas,omitempty"`
}

// RedisSettings defines the specification of the Redis database for secondary index
type RedisSettings struct {
	// + kubebuilder:default:=5.0.3
	Version string `json:"version,omitempty"`
	// + kubebuilder:default:=3
	Master int32 `json:"master,omitempty"`
	// + kubebuilder:default:=1
	Replicas int32 `json:"replicas,omitempty"`
}

// ReplicationSettings defines the specification of the replication framework
type ReplicationSettings struct {
	DbImage string `json:"dbimage,omitempty"`
	Image   string `json:"image,omitempty"`
}

// XrootdSettings defines the specification of the xrootd redirectors cluster
type XrootdSettings struct {
	Image    string `json:"image,omitempty"`
	Replicas int32  `json:"replicas,omitempty"`
}

// QservStatus defines the observed state of Qserv
// +k8s:openapi-gen=true
type QservStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Qserv is the Schema for the qservs API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Qserv struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QservSpec   `json:"spec,omitempty"`
	Status QservStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// QservList contains a list of Qserv
type QservList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Qserv `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Qserv{}, &QservList{})
}
