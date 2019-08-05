package qserv

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
)

type Volumes map[v1.Volume]struct{}

func (vs *Volumes) make() {
	*vs = Volumes(make(map[v1.Volume]struct{}))
}

func (vs *Volumes) addConfigMapVolume(name string) {
	volume := v1.Volume{
		Name: name,
		VolumeSource: v1.VolumeSource{
			ConfigMap: &v1.ConfigMapVolumeSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: name,
				},
			}}}
	var s struct{}
	(*vs)[volume] = s
}

func (vs *Volumes) addSecretVolume(name string) {
	volume := v1.Volume{
		Name: name,
		VolumeSource: v1.VolumeSource{
			Secret: &v1.SecretVolumeSource{
				SecretName: name,
			},
		}}

	var s struct{}
	(*vs)[volume] = s
}

func (vs Volumes) toSlice() []v1.Volume {
	var volumes []v1.Volume
	for k := range vs {
		volumes = append(volumes, k)
	}
	return volumes
}

func getConfigVolumes(service string) []v1.Volume {
	var volumes []v1.Volume

	var configName string
	executeMode := int32(0555)
	configName = fmt.Sprintf("config-%s-etc", service)
	volumes = append(volumes, v1.Volume{Name: configName, VolumeSource: v1.VolumeSource{ConfigMap: &v1.ConfigMapVolumeSource{
		LocalObjectReference: v1.LocalObjectReference{
			Name: configName,
		},
	}}})
	configName = fmt.Sprintf("config-%s-start", service)
	volumes = append(volumes, v1.Volume{Name: configName, VolumeSource: v1.VolumeSource{ConfigMap: &v1.ConfigMapVolumeSource{
		LocalObjectReference: v1.LocalObjectReference{
			Name: configName,
		},
		DefaultMode: &executeMode,
	}}})
	return volumes
}

func mountConfigVolumes(container *v1.Container, service string) {
	container.VolumeMounts = append(container.VolumeMounts, getConfigVolumeMounts(service)...)
}

func getConfigVolumeMounts(service string) []v1.VolumeMount {
	var volumeMounts []v1.VolumeMount
	volumeName := fmt.Sprintf("config-%s-etc", service)
	volumeMounts = append(volumeMounts, v1.VolumeMount{Name: volumeName, MountPath: "/config-etc"})
	volumeName = fmt.Sprintf("config-%s-start", service)
	volumeMounts = append(volumeMounts, v1.VolumeMount{Name: volumeName, MountPath: "/config-start"})
	return volumeMounts
}
