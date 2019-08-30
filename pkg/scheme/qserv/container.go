package qserv

import (
	"fmt"
	"path/filepath"

	qservv1alpha1 "github.com/lsst/qserv-operator/pkg/apis/qserv/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
	"github.com/lsst/qserv-operator/pkg/util"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func getMariadbImage(cr *qservv1alpha1.Qserv, component string) string {
	spec := cr.Spec
	var image string
	if component == constants.ReplName {
		image = spec.Replication.DbImage
	} else {
		image = spec.Worker.Image
	}
	return image
}

func getInitContainer(cr *qservv1alpha1.Qserv, component string) (v1.Container, VolumeSet) {
	sqlConfigMap := fmt.Sprintf("config-sql-%s", component)

	container := v1.Container{
		Name:  constants.InitDbName,
		Image: getMariadbImage(cr, component),
		Command: []string{
			"/config-start/mariadb-configure.sh",
		},
		Env: []v1.EnvVar{
			{
				Name:  "COMPONENT_NAME",
				Value: component,
			},
		},
		VolumeMounts: []v1.VolumeMount{
			getDataVolumeMount(),
			getEtcVolumeMount("mariadb"),
			getStartVolumeMount("mariadb"),
			{
				MountPath: filepath.Join("/", "secret-mariadb"),
				Name:      "secret-mariadb",
				ReadOnly:  false,
			},
			{
				MountPath: filepath.Join("/", "config-sql", component),
				Name:      sqlConfigMap,
				ReadOnly:  false,
			},
		},
	}

	var volumes VolumeSet
	volumes.make(nil)

	volumes.addConfigMapVolume(sqlConfigMap)
	volumes.addEtcStartVolumes(constants.MariadbName)
	volumes.addSecretVolume("secret-mariadb")

	return container, volumes
}

func getMariadbContainer(cr *qservv1alpha1.Qserv, component string) (v1.Container, VolumeSet) {

	var uservice string
	if component == constants.ReplName {
		uservice = constants.ReplDbName
	} else {
		uservice = constants.MariadbName
	}

	container := v1.Container{
		Name:  constants.MariadbName,
		Image: getMariadbImage(cr, component),
		Ports: []v1.ContainerPort{
			{
				Name:          constants.MariadbName,
				ContainerPort: constants.MariadbPort,
				Protocol:      v1.ProtocolTCP,
			},
		},
		Command:        constants.Command,
		LivenessProbe:  getLivenessProbe(constants.MariadbName),
		ReadinessProbe: getReadinessProbe(constants.MariadbName),
		VolumeMounts: []v1.VolumeMount{
			getDataVolumeMount(),
			getEtcVolumeMount(uservice),
			getStartVolumeMount(uservice),
			{
				MountPath: "/qserv/run/tmp",
				Name:      "tmp-volume",
				ReadOnly:  false,
			},
		},
	}

	// Volumes
	var volumes VolumeSet
	volumes.make(nil)

	volumes.addEmptyDirVolume("tmp-volume")
	volumes.addEtcStartVolumes(uservice)

	return container, volumes
}

func getProxyContainer(cr *qservv1alpha1.Qserv) (v1.Container, VolumeSet) {
	spec := cr.Spec
	container := v1.Container{
		Name:  constants.ProxyName,
		Image: spec.Worker.Image,
		Ports: []v1.ContainerPort{
			{
				Name:          constants.ProxyName,
				ContainerPort: constants.ProxyPort,
				Protocol:      v1.ProtocolTCP,
			},
		},
		LivenessProbe:  getLivenessProbe(constants.ProxyName),
		ReadinessProbe: getReadinessProbe(constants.ProxyName),
		Command:        constants.Command,
		Env: []v1.EnvVar{
			{
				Name:  "XROOTD_RDR_DN",
				Value: util.GetXrootdRedirectorName(cr),
			},
		},
		VolumeMounts: []v1.VolumeMount{
			// Used for mysql socket access
			// TODO move mysql socket in emptyDir?
			getDataVolumeMount(),
			getEtcVolumeMount(constants.ProxyName),
			getStartVolumeMount(constants.ProxyName),
		},
	}

	// Volumes
	var volumes VolumeSet
	volumes.make(nil)

	volumes.addEtcStartVolumes(constants.ProxyName)

	return container, volumes
}

func getWmgrContainer(cr *qservv1alpha1.Qserv) (v1.Container, VolumeSet) {
	spec := cr.Spec

	container := v1.Container{
		Name:  constants.WmgrName,
		Image: spec.Worker.Image,
		Ports: []v1.ContainerPort{
			{
				Name:          constants.WmgrName,
				ContainerPort: constants.WmgrPort,
				Protocol:      v1.ProtocolTCP,
			},
		},
		Command: constants.Command,
		Env: []v1.EnvVar{
			{
				Name:  "CZAR_DN",
				Value: util.GetCzarName(cr),
			},
		},
		LivenessProbe:  getLivenessProbe(constants.WmgrName),
		ReadinessProbe: getReadinessProbe(constants.WmgrName),
		VolumeMounts: []v1.VolumeMount{
			{
				MountPath: filepath.Join("/", "config-dot-qserv"),
				Name:      "config-dot-qserv",
				ReadOnly:  true,
			},
			{
				MountPath: "/qserv/run/tmp",
				Name:      "tmp-volume",
				ReadOnly:  false,
			},
			{
				MountPath: "/secret-mariadb",
				Name:      "secret-mariadb",
				ReadOnly:  true,
			},
			{
				MountPath: "/secret-wmgr",
				Name:      "secret-wmgr",
				ReadOnly:  true,
			},
			getDataVolumeMount(),
			getEtcVolumeMount(constants.WmgrName),
			getStartVolumeMount(constants.WmgrName),
		},
	}

	// Volumes
	var volumes VolumeSet
	volumes.make(nil)

	volumes.addConfigMapVolume("config-dot-qserv")
	volumes.addSecretVolume("secret-wmgr")
	volumes.addEmptyDirVolume("tmp-volume")
	volumes.addEtcStartVolumes("wmgr")

	// TODO Add volumes
	return container, volumes
}

func getXrootdContainers(cr *qservv1alpha1.Qserv, component string) ([]v1.Container, VolumeSet) {

	const (
		CMSD = iota
		XROOTD
	)

	spec := cr.Spec

	envRedirector := v1.EnvVar{
		Name:  "XROOTD_RDR_DN",
		Value: util.GetXrootdRedirectorName(cr),
	}

	volumeMounts := getXrootdVolumeMounts(component)

	containers := []v1.Container{
		{
			Name:    constants.CmsdName,
			Image:   spec.Worker.Image,
			Command: constants.Command,
			Args:    []string{"-S", "cmsd"},
			Env: []v1.EnvVar{
				envRedirector,
			},
			SecurityContext: &v1.SecurityContext{
				Capabilities: &v1.Capabilities{
					Add: []v1.Capability{
						v1.Capability("IPC_LOCK"),
					},
				},
			},
			VolumeMounts: volumeMounts,
		},
		{
			Name:  constants.XrootdName,
			Image: spec.Worker.Image,
			Ports: []v1.ContainerPort{
				{
					Name:          constants.XrootdName,
					ContainerPort: 1094,
					Protocol:      v1.ProtocolTCP,
				},
			},
			Command: constants.Command,
			Env: []v1.EnvVar{
				envRedirector,
			},
			LivenessProbe:  getLivenessProbe(constants.XrootdName),
			ReadinessProbe: getReadinessProbe(constants.XrootdName),
			SecurityContext: &v1.SecurityContext{
				Capabilities: &v1.Capabilities{
					Add: []v1.Capability{
						v1.Capability("IPC_LOCK"),
						v1.Capability("SYS_RESOURCE"),
					},
				},
			},
			VolumeMounts: volumeMounts,
		},
	}

	// Cmsd port is only open on redirectors, not on workers
	if component == constants.XrootdRedirectorName {
		containers[0].Ports = []v1.ContainerPort{
			{
				Name:          constants.CmsdName,
				ContainerPort: 2131,
				Protocol:      v1.ProtocolTCP,
			},
		}
		containers[0].LivenessProbe = getLivenessProbe(constants.CmsdName)
		containers[0].ReadinessProbe = getReadinessProbe(constants.CmsdName)
	}

	// Volumes
	var volumes VolumeSet
	volumes.make(nil)

	volumes.addEtcStartVolumes(constants.XrootdName)
	volumes.addEmptyDirVolume(constants.XrootdAdminPathVolumeName)

	return containers, volumes
}

func getLivenessProbe(portName string) *v1.Probe {
	return &v1.Probe{
		Handler: v1.Handler{
			TCPSocket: &v1.TCPSocketAction{
				Port: intstr.FromString(portName),
			},
		},
		InitialDelaySeconds: 10,
		PeriodSeconds:       10,
	}
}

func getReadinessProbe(portName string) *v1.Probe {
	return &v1.Probe{
		Handler: v1.Handler{
			TCPSocket: &v1.TCPSocketAction{
				Port: intstr.FromString(portName),
			},
		},
		InitialDelaySeconds: 10,
		PeriodSeconds:       5,
	}
}
