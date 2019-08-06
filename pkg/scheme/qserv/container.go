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

func getInitContainer(cr *qservv1alpha1.Qserv, component string) (v1.Container, VolumeSet) {
	spec := cr.Spec
	sqlConfigMap := fmt.Sprintf("config-sql-%s", component)

	container := v1.Container{
		Name:  "initdb",
		Image: spec.Worker.Image,
		Command: []string{
			"/config-start/mariadb-configure.sh",
		},
		Env: []v1.EnvVar{
			{
				Name:  "INSTANCE_NAME",
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
	volumes.addEtcStartVolumes("mariadb")
	volumes.addSecretVolume("secret-mariadb")

	return container, volumes
}

func getMariadbContainer(cr *qservv1alpha1.Qserv) (v1.Container, VolumeSet) {
	spec := cr.Spec
	container := v1.Container{
		Name:  constants.MariadbName,
		Image: spec.Worker.Image,
		Ports: []v1.ContainerPort{
			{
				Name:          constants.MariadbPortName,
				ContainerPort: constants.MariadbPort,
				Protocol:      v1.ProtocolTCP,
			},
		},
		Command: constants.Command,
		VolumeMounts: []v1.VolumeMount{
			getDataVolumeMount(),
			getEtcVolumeMount("mariadb"),
			getStartVolumeMount("mariadb"),
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
	volumes.addEtcStartVolumes("mariadb")

	return container, volumes
}

func getProxyContainer(cr *qservv1alpha1.Qserv) (v1.Container, VolumeSet) {
	spec := cr.Spec
	container := v1.Container{
		Name:  constants.MysqlProxyName,
		Image: spec.Worker.Image,
		Ports: []v1.ContainerPort{
			{
				Name:          constants.MysqlProxyPortName,
				ContainerPort: constants.MysqlProxyPort,
				Protocol:      v1.ProtocolTCP,
			},
		},
		Command: constants.Command,
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
			getEtcVolumeMount(constants.MysqlProxyName),
			getStartVolumeMount(constants.MysqlProxyName),
		},
	}

	// Volumes
	var volumes VolumeSet
	volumes.make(nil)

	volumes.addEtcStartVolumes(constants.MysqlProxyName)

	return container, volumes
}

func getWmgrContainer(cr *qservv1alpha1.Qserv) (v1.Container, VolumeSet) {
	spec := cr.Spec

	container := v1.Container{
		Name:  constants.WmgrName,
		Image: spec.Worker.Image,
		Ports: []v1.ContainerPort{
			{
				Name:          constants.WmgrPortName,
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

func getXrootdContainers(cr *qservv1alpha1.Qserv) ([]v1.Container, VolumeSet) {

	const (
		CMSD = iota
		XROOTD
	)

	spec := cr.Spec

	envRedirector := v1.EnvVar{
		Name:  "XROOTD_RDR_DN",
		Value: util.GetXrootdRedirectorName(cr),
	}

	containers := []v1.Container{
		{
			Name:    "cmsd",
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
			VolumeMounts: []v1.VolumeMount{
				getEtcVolumeMount(constants.XrootdName),
				getStartVolumeMount(constants.XrootdName),
			},
		},
		{
			Name:  "xrootd",
			Image: spec.Worker.Image,
			Ports: []v1.ContainerPort{
				{
					Name:          "xrootd",
					ContainerPort: 1094,
					Protocol:      v1.ProtocolTCP,
				},
			},
			Command: constants.Command,
			Env: []v1.EnvVar{
				envRedirector,
			},
			SecurityContext: &v1.SecurityContext{
				Capabilities: &v1.Capabilities{
					Add: []v1.Capability{
						v1.Capability("IPC_LOCK"),
						v1.Capability("SYS_RESOURCE"),
					},
				},
			},
			LivenessProbe: &v1.Probe{
				Handler: v1.Handler{
					TCPSocket: &v1.TCPSocketAction{
						Port: intstr.FromString("xrootd"),
					},
				},
				InitialDelaySeconds: 10,
				PeriodSeconds:       10,
			},
			ReadinessProbe: &v1.Probe{
				Handler: v1.Handler{
					TCPSocket: &v1.TCPSocketAction{
						Port: intstr.FromString("xrootd"),
					},
				},
				InitialDelaySeconds: 10,
				PeriodSeconds:       5,
			},
			VolumeMounts: []v1.VolumeMount{
				getEtcVolumeMount(constants.XrootdName),
				getStartVolumeMount(constants.XrootdName),
			},
		},
	}

	// Volumes
	var volumes VolumeSet
	volumes.make(nil)

	volumes.addEtcStartVolumes(constants.XrootdName)
	volumes.addEmptyDirVolume(constants.XrootdAdminPathVolumeName)

	return containers, volumes
}
