package qserv

import (
	"fmt"
	"path/filepath"

	"github.com/lsst/qserv-operator/pkg/constants"
	v1 "k8s.io/api/core/v1"
)

// VolumeSet contains a set of v1.Volume
type VolumeSet map[string]v1.Volume

func (vs *VolumeSet) make(volumeSets ...VolumeSet) {
	*vs = VolumeSet(make(map[string]v1.Volume))
	for _, vols := range volumeSets {
		for k, v := range map[string]v1.Volume(vols) {
			(*vs)[k] = v
		}
	}
}

func (vs *VolumeSet) add(vols VolumeSet) {
	for k, v := range vols {
		(*vs)[k] = v
	}
}

func (vs *VolumeSet) addConfigMapExecVolume(name string, executeMode *int32) {
	volume := v1.Volume{
		Name: name,
		VolumeSource: v1.VolumeSource{
			ConfigMap: &v1.ConfigMapVolumeSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: name,
				},
				DefaultMode: executeMode,
			}}}
	(*vs)[name] = volume
}

func (vs *VolumeSet) addConfigMapVolume(name string) {
	volume := v1.Volume{
		Name: name,
		VolumeSource: v1.VolumeSource{
			ConfigMap: &v1.ConfigMapVolumeSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: name,
				},
			}}}
	(*vs)[name] = volume
}

func (vs *VolumeSet) addEmptyDirVolume(name string) {
	volume := v1.Volume{
		Name: name,
		VolumeSource: v1.VolumeSource{
			EmptyDir: &v1.EmptyDirVolumeSource{},
		},
	}
	(*vs)[name] = volume
}

func (vs *VolumeSet) addSecretVolume(name string) {
	volume := v1.Volume{
		Name: name,
		VolumeSource: v1.VolumeSource{
			Secret: &v1.SecretVolumeSource{
				SecretName: name,
			},
		}}
	(*vs)[name] = volume
}

func (vs VolumeSet) toSlice() []v1.Volume {
	var volumes []v1.Volume
	for _, v := range vs {
		volumes = append(volumes, v)
	}
	return volumes
}

func (vs *VolumeSet) addEtcStartVolumes(microservice string) {

	configName := fmt.Sprintf("config-%s-etc", microservice)
	(*vs).addConfigMapVolume(configName)

	configName = fmt.Sprintf("config-%s-start", microservice)
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

func getAdminPathMount() v1.VolumeMount {
	return v1.VolumeMount{
		MountPath: filepath.Join("/", "qserv", "run", "tmp", "xrd"),
		Name:      constants.XrootdAdminPathVolumeName,
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
