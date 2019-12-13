package qserv

import (
	"fmt"
	"path/filepath"

	qservv1alpha1 "github.com/lsst/qserv-operator/pkg/apis/qserv/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
	"github.com/lsst/qserv-operator/pkg/util"
	v1 "k8s.io/api/core/v1"
)

// VolumeSet contains a set of v1.Volume
type VolumeSet map[string]v1.Volume

// InstanceVolumeSet contains a set of v1.Volume for a given Qserv instance
type InstanceVolumeSet struct {
	volumeSet VolumeSet
	cr        *qservv1alpha1.Qserv
}

func (vs *VolumeSet) make(volumeSets ...VolumeSet) {
	*vs = VolumeSet(make(map[string]v1.Volume))
	for _, vols := range volumeSets {
		for k, v := range map[string]v1.Volume(vols) {
			(*vs)[k] = v
		}
	}
}

func (ivs *InstanceVolumeSet) make(cr *qservv1alpha1.Qserv) {
	ivs.volumeSet = VolumeSet(make(map[string]v1.Volume))
	ivs.cr = cr
}

func (vs *VolumeSet) add(vols VolumeSet) {
	for k, v := range vols {
		(*vs)[k] = v
	}
}

func (ivs *InstanceVolumeSet) addConfigMapExecVolume(container constants.ContainerName, executeMode *int32) {

	suffix := fmt.Sprintf("%s-start", container)
	configmapName := util.PrefixConfigmap(ivs.cr, suffix)
	volumeName := util.GetConfigVolumeName(suffix)

	volume := v1.Volume{
		Name: volumeName,
		VolumeSource: v1.VolumeSource{
			ConfigMap: &v1.ConfigMapVolumeSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: configmapName,
				},
				DefaultMode: executeMode,
			}}}
	ivs.volumeSet[volumeName] = volume
}

func (ivs *InstanceVolumeSet) addConfigMapVolume(suffix string) {

	configmapName := util.PrefixConfigmap(ivs.cr, suffix)
	volumeName := util.GetConfigVolumeName(suffix)

	volume := v1.Volume{
		Name: volumeName,
		VolumeSource: v1.VolumeSource{
			ConfigMap: &v1.ConfigMapVolumeSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: configmapName,
				},
			}}}
	ivs.volumeSet[volumeName] = volume
}

func (ivs *InstanceVolumeSet) addEmptyDirVolume(name string) {
	volume := v1.Volume{
		Name: name,
		VolumeSource: v1.VolumeSource{
			EmptyDir: &v1.EmptyDirVolumeSource{},
		},
	}
	ivs.volumeSet[name] = volume
}

func (ivs *InstanceVolumeSet) addSecretVolume(containerName constants.ContainerName) {
	secretName := util.GetSecretName(ivs.cr, containerName)
	volume := v1.Volume{
		Name: util.GetSecretVolumeName(containerName),
		VolumeSource: v1.VolumeSource{
			Secret: &v1.SecretVolumeSource{
				SecretName: secretName,
			},
		}}
	ivs.volumeSet[secretName] = volume
}

func (vs VolumeSet) toSlice() []v1.Volume {
	var volumes []v1.Volume
	for _, v := range vs {
		volumes = append(volumes, v)
	}
	return volumes
}

func (ivs *InstanceVolumeSet) addEtcVolume(containerName constants.ContainerName) {
	suffix := fmt.Sprintf("%s-etc", containerName)
	ivs.addConfigMapVolume(suffix)
}

func (ivs *InstanceVolumeSet) addStartVolume(containerName constants.ContainerName) {
	mode := int32(0555)
	ivs.addConfigMapExecVolume(containerName, &mode)
}

func (ivs *InstanceVolumeSet) addEtcStartVolumes(containerName constants.ContainerName) {
	ivs.addEtcVolume(containerName)
	ivs.addStartVolume(containerName)
}

func getDataVolumeMount() v1.VolumeMount {
	return v1.VolumeMount{
		MountPath: filepath.Join("/", "qserv", "data"),
		Name:      GetVolumeClaimTemplateName(),
		ReadOnly:  false,
	}
}

func getAdminPathMount() v1.VolumeMount {
	return v1.VolumeMount{
		MountPath: filepath.Join("/", "tmp", "xrd"),
		Name:      constants.XrootdAdminPathVolumeName,
		ReadOnly:  false,
	}
}

func getEtcVolumeMount(microservice constants.ContainerName) v1.VolumeMount {
	volumeName := fmt.Sprintf("config-%s-etc", microservice)
	return v1.VolumeMount{Name: volumeName, MountPath: "/config-etc"}
}

func getSecretVolumeMount(containerName constants.ContainerName) v1.VolumeMount {
	secretVolumeName := util.GetSecretVolumeName(containerName)
	return v1.VolumeMount{
		MountPath: filepath.Join("/", secretVolumeName),
		Name:      secretVolumeName,
		ReadOnly:  true}
}

func getStartVolumeMount(microservice constants.ContainerName) v1.VolumeMount {
	volumeName := fmt.Sprintf("config-%s-start", microservice)
	return v1.VolumeMount{Name: volumeName, MountPath: "/config-start"}
}

func getTmpVolumeMount() v1.VolumeMount {
	return v1.VolumeMount{
		MountPath: "/qserv/run/tmp",
		Name:      "tmp-volume",
		ReadOnly:  false,
	}
}

func getXrootdVolumeMounts(component constants.ComponentName) []v1.VolumeMount {
	volumeMounts := []v1.VolumeMount{
		getAdminPathMount(),
		getEtcVolumeMount(constants.XrootdName),
		getStartVolumeMount(constants.XrootdName),
	}

	// xrootd/cmsd workers only
	if component == constants.WorkerName {
		volumeMounts = append(volumeMounts, getDataVolumeMount())
	}
	return volumeMounts
}
