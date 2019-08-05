package qserv

import (
	"fmt"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
)

type Volumes map[v1.Volume]struct{}

func (vs *Volumes) make(volumesList ...Volumes) {
	*vs = Volumes(make(map[v1.Volume]struct{}))
	for _, vols := range volumesList {
		for k := range map[v1.Volume]struct{}(vols) {
			(*vs)[k] = struct{}{}
		}
	}
}

func (vs *Volumes) add(vols Volumes) {
	for k := range vols {
		(*vs)[k] = struct{}{}
	}
}

func (vs *Volumes) addConfigMapExecVolume(name string, executeMode *int32) {
	volume := v1.Volume{
		Name: name,
		VolumeSource: v1.VolumeSource{
			ConfigMap: &v1.ConfigMapVolumeSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: name,
				},
				DefaultMode: executeMode,
			}}}
	(*vs)[volume] = struct{}{}
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
	(*vs)[volume] = struct{}{}
}

func (vs *Volumes) addEmptyDirVolume(name string) {
	volume := v1.Volume{
		Name: name,
		VolumeSource: v1.VolumeSource{
			EmptyDir: &v1.EmptyDirVolumeSource{},
		},
	}
	(*vs)[volume] = struct{}{}
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

func (vs *Volumes) addEtcStartVolumes(component string) {

	configName := fmt.Sprintf("config-%s-etc", component)
	(*vs).addConfigMapVolume(configName)

	configName = fmt.Sprintf("config-%s-start", component)
	mode := int32(0555)
	(*vs).addConfigMapExecVolume(configName, &mode)
}

func getDataVolumeMount() v1.VolumeMount {
	return v1.VolumeMount{
		MountPath: filepath.Join("/", "qserv", "data"),
		Name:      "qserv-data",
		ReadOnly:  false,
	}
}

func getEtcVolumeMount(microservice string) v1.VolumeMount {
	volumeName := fmt.Sprintf("config-%s-etc", microservice)
	return v1.VolumeMount{Name: volumeName, MountPath: "/config-etc"}
}

func getStartVolumeMount(microservice string) v1.VolumeMount {
	volumeName := fmt.Sprintf("config-%s-start", microservice)
	return v1.VolumeMount{Name: volumeName, MountPath: "/config-start"}
}
