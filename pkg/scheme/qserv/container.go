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

func getInitContainer(cr *qservv1alpha1.Qserv, component string) (v1.Container, []v1.Volume) {
	spec := cr.Spec
	trueVal := false
	sqlConfigMap := fmt.Sprintf("config-sql-%s", component)

	container := v1.Container{
		Name:  "initdb",
		Image: spec.Worker.Image,
		Command: []string{
			"/config-start/mariadb-configure.sh",
		},
		Env: []v1.EnvVar{
			{
				Name: "CZAR",
				ValueFrom: &v1.EnvVarSource{
					ConfigMapKeyRef: &v1.ConfigMapKeySelector{
						LocalObjectReference: v1.LocalObjectReference{Name: "config-domainnames"},
						Key:                  "CZAR",
						Optional:             &trueVal,
					},
				},
			},
		},
		VolumeMounts: []v1.VolumeMount{
			{
				MountPath: filepath.Join("/", "qserv", "data"),
				Name:      "qserv-data",
				ReadOnly:  false,
			},
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
	var volumes Volumes
	volumes.make()

	volumes.addConfigMapVolume("config-domainnames")
	volumes.addConfigMapVolume(sqlConfigMap)

	volumes.addSecretVolume("secret-mariadb")

	return container, volumes.toSlice()
}

func getWmgrContainer(cr *qservv1alpha1.Qserv) (v1.Container, []v1.Volume) {
	spec := cr.Spec
	trueVal := false

	container := v1.Container{
		Name:  "wmgr",
		Image: spec.Worker.Image,
		Ports: []v1.ContainerPort{
			{
				Name:          "wmgr",
				ContainerPort: 5012,
				Protocol:      v1.ProtocolTCP,
			},
		},
		Command: constants.Command,
		Env: []v1.EnvVar{
			{
				Name: "CZAR_DN",
				ValueFrom: &v1.EnvVarSource{
					ConfigMapKeyRef: &v1.ConfigMapKeySelector{
						LocalObjectReference: v1.LocalObjectReference{Name: "config-domainnames"},
						Key:                  "CZAR_DN",
						Optional:             &trueVal,
					},
				},
			},
		},
		VolumeMounts: []v1.VolumeMount{
			{
				MountPath: "/qserv/data",
				Name:      "qserv-data",
				ReadOnly:  false,
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
		},
	}

	mountConfigVolumes(&container, "wmgr")

	// Volumes
	var volumes []v1.Volume

	secretName := "secret-wmgr"
	volumes = append(volumes, v1.Volume{
		Name: secretName,
		VolumeSource: v1.VolumeSource{
			Secret: &v1.SecretVolumeSource{
				SecretName: secretName,
			},
		}})

	volumes = append(volumes, getConfigVolumes("wmgr")...)

	// TODO Add volumes
	return container, volumes
}

func getXrootdContainers(cr *qservv1alpha1.Qserv) ([]v1.Container, []v1.Volume) {

	const (
		CMSD = iota
		XROOTD
	)

	spec := cr.Spec
	redirectorName := util.GetXrootdRedirectorName(cr)

	envRedirector := v1.EnvVar{

		Name:  "XROOTD_RDR_DN",
		Value: redirectorName,
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
		},
	}
	mountConfigVolumes(&containers[CMSD], constants.XrootdName)
	mountConfigVolumes(&containers[XROOTD], constants.XrootdName)

	// Volumes
	var volumes []v1.Volume

	volumes = append(volumes, getConfigVolumes(constants.XrootdName)...)

	volumes = append(volumes, v1.Volume{Name: "xrootd-adminpath", VolumeSource: v1.VolumeSource{EmptyDir: &v1.EmptyDirVolumeSource{}}})

	return containers, volumes
}
