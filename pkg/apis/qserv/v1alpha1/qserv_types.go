package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// QservSpec defines the desired state of Qserv
// +k8s:openapi-gen=true
type QservSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	// Czar defines the settings for czar cluster
	Czar CzarSettings `json:"czar,omitempty"`

	// Worker defines the settings for worker cluster
	Worker WorkerSettings `json:"worker,omitempty"`
}


// CzarSettings defines the specification of the redis cluster
type CzarSettings struct {
	Image             string                     `json:"image,omitempty"`
	ImagePullPolicy   corev1.PullPolicy          `json:"imagePullPolicy,omitempty"`
	Replicas          int32                      `json:"replicas,omitempty"`
	Resources         RedisResources             `json:"resources,omitempty"`
	CustomConfig      []string                   `json:"customConfig,omitempty"`
	Command           []string                   `json:"command,omitempty"`
	ShutdownConfigMap string                     `json:"shutdownConfigMap,omitempty"`
	Storage           RedisStorage               `json:"storage,omitempty"`
	Exporter          RedisExporter              `json:"exporter,omitempty"`
}

// WorkerSettings defines the specification of the sentinel cluster
type WorkerSettings struct {
	Image           string                     `json:"image,omitempty"`
	Replicas        int32                      `json:"replicas,omitempty"`
	CustomConfig    []string                   `json:"customConfig,omitempty"`
	Command         []string                   `json:"command,omitempty"`
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
