//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2021.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CzarSettings) DeepCopyInto(out *CzarSettings) {
	*out = *in
	in.Affinity.DeepCopyInto(&out.Affinity)
	in.ProxyResources.DeepCopyInto(&out.ProxyResources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CzarSettings.
func (in *CzarSettings) DeepCopy() *CzarSettings {
	if in == nil {
		return nil
	}
	out := new(CzarSettings)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DevelSettings) DeepCopyInto(out *DevelSettings) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DevelSettings.
func (in *DevelSettings) DeepCopy() *DevelSettings {
	if in == nil {
		return nil
	}
	out := new(DevelSettings)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngestSettings) DeepCopyInto(out *IngestSettings) {
	*out = *in
	in.Affinity.DeepCopyInto(&out.Affinity)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngestSettings.
func (in *IngestSettings) DeepCopy() *IngestSettings {
	if in == nil {
		return nil
	}
	out := new(IngestSettings)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Qserv) DeepCopyInto(out *Qserv) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Qserv.
func (in *Qserv) DeepCopy() *Qserv {
	if in == nil {
		return nil
	}
	out := new(Qserv)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Qserv) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *QservList) DeepCopyInto(out *QservList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Qserv, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new QservList.
func (in *QservList) DeepCopy() *QservList {
	if in == nil {
		return nil
	}
	out := new(QservList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *QservList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *QservSpec) DeepCopyInto(out *QservSpec) {
	*out = *in
	in.Czar.DeepCopyInto(&out.Czar)
	in.Ingest.DeepCopyInto(&out.Ingest)
	out.Devel = in.Devel
	in.QueryService.DeepCopyInto(&out.QueryService)
	in.Replication.DeepCopyInto(&out.Replication)
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Worker.DeepCopyInto(&out.Worker)
	in.Xrootd.DeepCopyInto(&out.Xrootd)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new QservSpec.
func (in *QservSpec) DeepCopy() *QservSpec {
	if in == nil {
		return nil
	}
	out := new(QservSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *QservStatus) DeepCopyInto(out *QservStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]metav1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new QservStatus.
func (in *QservStatus) DeepCopy() *QservStatus {
	if in == nil {
		return nil
	}
	out := new(QservStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *QueryServiceSettings) DeepCopyInto(out *QueryServiceSettings) {
	*out = *in
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new QueryServiceSettings.
func (in *QueryServiceSettings) DeepCopy() *QueryServiceSettings {
	if in == nil {
		return nil
	}
	out := new(QueryServiceSettings)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReplicationSettings) DeepCopyInto(out *ReplicationSettings) {
	*out = *in
	in.Affinity.DeepCopyInto(&out.Affinity)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReplicationSettings.
func (in *ReplicationSettings) DeepCopy() *ReplicationSettings {
	if in == nil {
		return nil
	}
	out := new(ReplicationSettings)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkerSettings) DeepCopyInto(out *WorkerSettings) {
	*out = *in
	in.Affinity.DeepCopyInto(&out.Affinity)
	in.ReplicationResources.DeepCopyInto(&out.ReplicationResources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkerSettings.
func (in *WorkerSettings) DeepCopy() *WorkerSettings {
	if in == nil {
		return nil
	}
	out := new(WorkerSettings)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *XrootdRedirectorSettings) DeepCopyInto(out *XrootdRedirectorSettings) {
	*out = *in
	in.Affinity.DeepCopyInto(&out.Affinity)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new XrootdRedirectorSettings.
func (in *XrootdRedirectorSettings) DeepCopy() *XrootdRedirectorSettings {
	if in == nil {
		return nil
	}
	out := new(XrootdRedirectorSettings)
	in.DeepCopyInto(out)
	return out
}
