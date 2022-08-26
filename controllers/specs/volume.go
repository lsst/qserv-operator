package specs

import (
	"fmt"
	"path/filepath"

	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/controllers/constants"
	"github.com/lsst/qserv-operator/controllers/util"
	"golang.org/x/exp/maps"
	v1 "k8s.io/api/core/v1"
)

// VolumeSet contains a set of v1.Volume
type VolumeSet map[string]v1.Volume

// InstanceVolumeSet contains a set of v1.Volume for a given Qserv instance
type InstanceVolumeSet struct {
	volumeSet VolumeSet
	cr        *qservv1beta1.Qserv
}

func (vs *VolumeSet) make(volumeSets ...VolumeSet) {
	*vs = VolumeSet(make(map[string]v1.Volume))
	for _, vols := range volumeSets {
		for k, v := range map[string]v1.Volume(vols) {
			(*vs)[k] = v
		}
	}
}

func (ivs *InstanceVolumeSet) make(cr *qservv1beta1.Qserv) {
	ivs.volumeSet = VolumeSet(make(map[string]v1.Volume))
	ivs.cr = cr
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

func (ivs *InstanceVolumeSet) addCorePathVolume(path string) {
	name := constants.CorePathVolumeName
	volume := v1.Volume{
		Name: name,
		VolumeSource: v1.VolumeSource{
			HostPath: &v1.HostPathVolumeSource{
				Path: path,
			},
		},
	}
	ivs.volumeSet[name] = volume
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

func (vs VolumeSet) getNames() []string {
	names := maps.Keys(vs)
	return names
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
	if util.HasValue(string(containerName), constants.WithEtcStartConfigmaps) {
		ivs.addStartVolume(containerName)
		ivs.addEtcVolume(containerName)
	} else if util.HasValue(string(containerName), constants.WithStartConfigmap) {
		ivs.addStartVolume(containerName)
	}
}

func getAdminPathMount() v1.VolumeMount {
	return v1.VolumeMount{
		MountPath: constants.XrootdAdminPath,
		Name:      constants.XrootdAdminPathVolumeName,
		ReadOnly:  false,
	}
}

func getCorePathVolumeMount(mountPath string) v1.VolumeMount {
	volumeName := constants.CorePathVolumeName
	return v1.VolumeMount{Name: volumeName, MountPath: mountPath}
}

func getDataVolumeMount() v1.VolumeMount {
	return v1.VolumeMount{
		MountPath: filepath.Join("/", "qserv", "data"),
		Name:      constants.DataVolumeClaimTemplateName,
		ReadOnly:  false,
	}
}

func getMysqlEtcVolumeMount(container constants.ContainerName) v1.VolumeMount {
	volumeName := fmt.Sprintf("config-%s-etc", container)
	return v1.VolumeMount{
		Name:      volumeName,
		MountPath: "/etc/mysql/my.cnf",
		SubPath:   "my.cnf"}
}

func getSecretVolumeMount(containerName constants.ContainerName) v1.VolumeMount {
	secretVolumeName := util.GetSecretVolumeName(containerName)
	return v1.VolumeMount{
		MountPath: filepath.Join("/", secretVolumeName),
		Name:      secretVolumeName,
		ReadOnly:  true}
}

func getEtcVolumeMount(microservice constants.ContainerName) v1.VolumeMount {
	volumeName := fmt.Sprintf("config-%s-etc", microservice)
	return v1.VolumeMount{Name: volumeName, MountPath: constants.ConfigmapPathEtc}
}

func getStartVolumeMount(microservice constants.ContainerName) v1.VolumeMount {
	volumeName := fmt.Sprintf("config-%s-start", microservice)
	return v1.VolumeMount{Name: volumeName, MountPath: constants.ConfigmapPathStart}
}

func getTmpVolumeMount() v1.VolumeMount {
	return v1.VolumeMount{
		MountPath: "/qserv/run/tmp",
		Name:      "tmp-volume",
		ReadOnly:  false,
	}
}
func getXrootdVolumeMounts(containerName constants.ContainerName) []v1.VolumeMount {
	volumeMounts := []v1.VolumeMount{
		getAdminPathMount(),
	}

	volumeMounts = appendEtcStartVolumeMounts(volumeMounts, containerName)

	// xrootd/cmsd workers only
	if containerName == constants.CmsdServerName || containerName == constants.XrootdServerName {
		volumeMounts = append(volumeMounts, getDataVolumeMount())
	}
	return volumeMounts
}

func appendEtcStartVolumeMounts(volumeMounts []v1.VolumeMount, containerName constants.ContainerName) []v1.VolumeMount {
	if util.HasValue(string(containerName), constants.WithEtcStartConfigmaps) {
		volumeMounts = append(volumeMounts, getStartVolumeMount(containerName))
		volumeMounts = append(volumeMounts, getEtcVolumeMount(containerName))
	} else if util.HasValue(string(containerName), constants.WithStartConfigmap) {
		volumeMounts = append(volumeMounts, getStartVolumeMount(containerName))
	}
	return volumeMounts
}

func getNames(volumeMounts []v1.VolumeMount) []string {
	names := []string{}
	for _, v := range volumeMounts {
		names = append(names, v.Name)
	}
	return names
}
