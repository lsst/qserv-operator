// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CzarSettings) DeepCopyInto(out *CzarSettings) {
	*out = *in
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
func (in *IngestSettings) DeepCopyInto(out *IngestSettings) {
	*out = *in
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
	out.Status = in.Status
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
	out.Czar = in.Czar
	out.Ingest = in.Ingest
	if in.Redis != nil {
		in, out := &in.Redis, &out.Redis
		*out = new(RedisSettings)
		**out = **in
	}
	out.Replication = in.Replication
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	out.Worker = in.Worker
	out.Xrootd = in.Xrootd
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
func (in *RedisSettings) DeepCopyInto(out *RedisSettings) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RedisSettings.
func (in *RedisSettings) DeepCopy() *RedisSettings {
	if in == nil {
		return nil
	}
	out := new(RedisSettings)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReplicationSettings) DeepCopyInto(out *ReplicationSettings) {
	*out = *in
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
func (in *XrootdSettings) DeepCopyInto(out *XrootdSettings) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new XrootdSettings.
func (in *XrootdSettings) DeepCopy() *XrootdSettings {
	if in == nil {
		return nil
	}
	out := new(XrootdSettings)
	in.DeepCopyInto(out)
	return out
}
