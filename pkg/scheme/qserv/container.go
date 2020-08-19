package qserv

import (
	"fmt"
	"path/filepath"
	"strconv"

	qservv1alpha1 "github.com/lsst/qserv-operator/pkg/apis/qserv/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
	"github.com/lsst/qserv-operator/pkg/util"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func getInitContainer(cr *qservv1alpha1.Qserv, component constants.PodClass) (v1.Container, VolumeSet) {
	componentName := string(component)

	sqlConfigSuffix := fmt.Sprintf("sql-%s", component)

	dbContainerName := constants.GetDbContainerName(component)

	container := v1.Container{
		Name:            string(constants.InitDbName),
		Image:           getMariadbImage(cr, component),
		ImagePullPolicy: cr.Spec.ImagePullPolicy,
		Command: []string{
			"/config-start/initdb.sh",
		},
		Env: []v1.EnvVar{
			{
				Name:  "COMPONENT_NAME",
				Value: componentName,
			},
		},
		VolumeMounts: []v1.VolumeMount{
			getDataVolumeMount(),
			getEtcVolumeMount(dbContainerName),
			// db startup script and root passwords are shared
			getStartVolumeMount(constants.InitDbName),
			getSecretVolumeMount(constants.MariadbName),
			{
				MountPath: filepath.Join("/", "config-sql", componentName),
				Name:      util.GetConfigVolumeName(sqlConfigSuffix),
				ReadOnly:  true,
			},
		},
	}

	var volumes InstanceVolumeSet
	volumes.make(cr)

	volumes.addConfigMapVolume(sqlConfigSuffix)
	volumes.addEtcVolume(dbContainerName)
	volumes.addStartVolume(constants.InitDbName)
	volumes.addSecretVolume(constants.MariadbName)

	if dbContainerName == constants.ReplDbName || dbContainerName == constants.IngestDbName {
		container.VolumeMounts = append(container.VolumeMounts, getSecretVolumeMount(dbContainerName))
		volumes.addSecretVolume(dbContainerName)
	}

	return container, volumes.volumeSet
}

func getMariadbContainer(cr *qservv1alpha1.Qserv, pod constants.PodClass) (v1.Container, VolumeSet) {

	dbContainerName := constants.GetDbContainerName(pod)

	mariadbPortName := string(constants.MariadbName)

	// Volumes
	var volumes InstanceVolumeSet
	volumes.make(cr)

	volumes.addEmptyDirVolume("tmp-volume")
	volumes.addEtcStartVolumes(dbContainerName)

	// Container
	container := v1.Container{
		Command:         constants.Command,
		Image:           getMariadbImage(cr, pod),
		ImagePullPolicy: cr.Spec.ImagePullPolicy,
		Name:            string(dbContainerName),
		LivenessProbe:   getTCPProbe(constants.MariadbPortName, 10),
		Ports: []v1.ContainerPort{
			{
				Name:          mariadbPortName,
				ContainerPort: constants.MariadbPort,
				Protocol:      v1.ProtocolTCP,
			},
		},
		ReadinessProbe: getTCPProbe(constants.MariadbPortName, 5),
		VolumeMounts: []v1.VolumeMount{
			getDataVolumeMount(),
			getEtcVolumeMount(dbContainerName),
			getStartVolumeMount(dbContainerName),
			getTmpVolumeMount(),
		},
	}

	return container, volumes.volumeSet
}

func getMariadbImage(cr *qservv1alpha1.Qserv, component constants.PodClass) string {
	spec := cr.Spec
	var image string
	if component == constants.ReplDb {
		image = spec.Replication.DbImage
	} else if component == constants.IngestDb {
		image = spec.Ingest.DbImage
	} else if component == constants.Worker {
		image = spec.Worker.Image
	} else if component == constants.Czar {
		image = spec.Czar.Image
	}
	return image
}

func getProxyContainer(cr *qservv1alpha1.Qserv) (v1.Container, VolumeSet) {
	spec := cr.Spec
	container := v1.Container{
		Command:         constants.Command,
		Image:           spec.Czar.Image,
		ImagePullPolicy: spec.ImagePullPolicy,
		Name:            string(constants.ProxyName),
		Ports: []v1.ContainerPort{
			{
				Name:          string(constants.ProxyName),
				ContainerPort: constants.ProxyPort,
				Protocol:      v1.ProtocolTCP,
			},
		},
		LivenessProbe:  getTCPProbe(constants.ProxyPortName, 10),
		ReadinessProbe: getTCPProbe(constants.ProxyPortName, 5),
		VolumeMounts: []v1.VolumeMount{
			// Used for mysql socket access
			// TODO move mysql socket in emptyDir?
			getDataVolumeMount(),
			getEtcVolumeMount(constants.ProxyName),
			getStartVolumeMount(constants.ProxyName),
		},
	}

	// Volumes
	var volumes InstanceVolumeSet
	volumes.make(cr)

	volumes.addEtcStartVolumes(constants.ProxyName)

	return container, volumes.volumeSet
}

func getReplicationCtlContainer(cr *qservv1alpha1.Qserv) (v1.Container, VolumeSet) {
	spec := cr.Spec

	var probeTimeoutSeconds int32 = 3

	container := v1.Container{
		Command:         constants.Command,
		LivenessProbe:   getHTTPProbe(constants.ReplicationControllerPortName, 10, probeTimeoutSeconds, "meta/version"),
		ReadinessProbe:  getHTTPProbe(constants.ReplicationControllerPortName, 5, probeTimeoutSeconds, "meta/version"),
		Name:            string(constants.ReplCtlName),
		Image:           spec.Replication.Image,
		ImagePullPolicy: spec.ImagePullPolicy,
		Env: []v1.EnvVar{
			{
				Name:  "WORKER_COUNT",
				Value: strconv.FormatInt(int64(spec.Worker.Replicas), 10),
			},
			{
				Name:  "REPL_DB_DN",
				Value: util.GetName(cr, string(constants.ReplDbName)),
			},
		},
		Ports: []v1.ContainerPort{
			{
				Name:          constants.ReplicationControllerPortName,
				ContainerPort: constants.ReplicationControllerPort,
				Protocol:      v1.ProtocolTCP,
			},
		},
		VolumeMounts: []v1.VolumeMount{
			v1.VolumeMount{
				MountPath: filepath.Join("/", "qserv", "data"),
				Name:      "data",
				ReadOnly:  false,
			},
			getEtcVolumeMount(constants.ReplCtlName),
			getStartVolumeMount(constants.ReplCtlName),
			getSecretVolumeMount(constants.ReplDbName),
			getSecretVolumeMount(constants.MariadbName),
		},
	}

	var volumes InstanceVolumeSet
	volumes.make(cr)

	volumes.addEtcStartVolumes(constants.ReplCtlName)
	volumes.addDataVolume(cr)
	volumes.addSecretVolume(constants.ReplDbName)
	volumes.addSecretVolume(constants.MariadbName)

	return container, volumes.volumeSet
}

func getReplicationWrkContainer(cr *qservv1alpha1.Qserv) (v1.Container, VolumeSet) {
	spec := cr.Spec

	container := v1.Container{
		Name:            string(constants.ReplWrkName),
		Image:           spec.Replication.Image,
		ImagePullPolicy: spec.ImagePullPolicy,
		Command:         constants.Command,
		Env: []v1.EnvVar{
			{
				Name:  "WORKER_COUNT",
				Value: strconv.FormatInt(int64(spec.Worker.Replicas), 10),
			},
			{
				Name:  "REPL_DB_DN",
				Value: util.GetName(cr, string(constants.ReplDbName)),
			},
		},
		SecurityContext: &v1.SecurityContext{
			RunAsUser: &constants.QservUID,
		},
		VolumeMounts: []v1.VolumeMount{
			getDataVolumeMount(),
			getEtcVolumeMount(constants.ReplWrkName),
			getStartVolumeMount(constants.ReplWrkName),
			getSecretVolumeMount(constants.MariadbName),
			getSecretVolumeMount(constants.ReplDbName),
		},
	}

	var volumes InstanceVolumeSet
	volumes.make(cr)

	volumes.addEtcStartVolumes(constants.ReplWrkName)
	volumes.addSecretVolume(constants.MariadbName)
	volumes.addSecretVolume(constants.ReplDbName)

	return container, volumes.volumeSet
}

func getWmgrContainer(cr *qservv1alpha1.Qserv) (v1.Container, VolumeSet) {
	dotQserv := "dot-qserv"
	dotQservConfigVolume := util.GetConfigVolumeName(dotQserv)
	container := v1.Container{
		Name:            string(constants.WmgrName),
		Image:           cr.Spec.Worker.Image,
		ImagePullPolicy: cr.Spec.ImagePullPolicy,
		Ports: []v1.ContainerPort{
			{
				Name:          constants.WmgrPortName,
				ContainerPort: constants.WmgrPort,
				Protocol:      v1.ProtocolTCP,
			},
		},
		Command:        constants.Command,
		LivenessProbe:  getTCPProbe(constants.WmgrPortName, 10),
		ReadinessProbe: getTCPProbe(constants.WmgrPortName, 5),
		VolumeMounts: []v1.VolumeMount{
			{
				MountPath: filepath.Join("/", dotQservConfigVolume),
				Name:      dotQservConfigVolume,
				ReadOnly:  true,
			},
			getTmpVolumeMount(),
			getSecretVolumeMount(constants.MariadbName),
			getSecretVolumeMount(constants.WmgrName),
			getDataVolumeMount(),
			getEtcVolumeMount(constants.WmgrName),
			getStartVolumeMount(constants.WmgrName),
		},
	}

	// Volumes
	var volumes InstanceVolumeSet
	volumes.make(cr)

	volumes.addConfigMapVolume(dotQserv)
	volumes.addSecretVolume(constants.MariadbName)
	volumes.addSecretVolume(constants.WmgrName)
	volumes.addEmptyDirVolume("tmp-volume")
	volumes.addEtcStartVolumes(constants.WmgrName)

	return container, volumes.volumeSet
}

func getXrootdContainers(cr *qservv1alpha1.Qserv, component constants.PodClass) ([]v1.Container, VolumeSet) {

	const (
		CMSD = iota
		XROOTD
	)

	spec := cr.Spec

	volumeMounts := getXrootdVolumeMounts(component)

	containers := []v1.Container{
		{
			Name:            string(constants.CmsdName),
			Image:           spec.Worker.Image,
			ImagePullPolicy: cr.Spec.ImagePullPolicy,
			Command:         constants.Command,
			Args:            []string{"-S", "cmsd"},
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
			Name:            string(constants.XrootdName),
			Image:           spec.Worker.Image,
			ImagePullPolicy: cr.Spec.ImagePullPolicy,
			Ports: []v1.ContainerPort{
				{
					Name:          constants.XrootdPortName,
					ContainerPort: constants.XrootdPort,
					Protocol:      v1.ProtocolTCP,
				},
			},
			Command:        constants.Command,
			LivenessProbe:  getTCPProbe(constants.XrootdPortName, 10),
			ReadinessProbe: getTCPProbe(constants.XrootdPortName, 5),
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
	if component == constants.XrootdRedirector {
		containers[0].Ports = []v1.ContainerPort{
			{
				Name:          constants.CmsdPortName,
				ContainerPort: constants.CmsdPort,
				Protocol:      v1.ProtocolTCP,
			},
		}
		containers[0].LivenessProbe = getTCPProbe(constants.CmsdPortName, 10)
		containers[0].ReadinessProbe = getTCPProbe(constants.CmsdPortName, 5)
	}

	// Volumes
	var volumes InstanceVolumeSet
	volumes.make(cr)

	volumes.addEtcStartVolumes(constants.XrootdName)
	volumes.addEmptyDirVolume(constants.XrootdAdminPathVolumeName)

	return containers, volumes.volumeSet
}

type NetworkAction string

const (
	httpAction NetworkAction = "http"
	tcpAction  NetworkAction = "tcp"
)

func getHTTPProbe(portName string, periodSeconds int32, timeoutSeconds int32, path string) *v1.Probe {
	var handler *v1.Handler
	handler = &v1.Handler{
		HTTPGet: &v1.HTTPGetAction{
			Path: path,
			Port: intstr.FromString(portName),
		},
	}
	return &v1.Probe{
		Handler:             *handler,
		InitialDelaySeconds: 10,
		PeriodSeconds:       periodSeconds,
		TimeoutSeconds:      timeoutSeconds,
	}
}

func getTCPProbe(portName string, periodSeconds int32) *v1.Probe {
	var handler *v1.Handler
	handler = &v1.Handler{
		TCPSocket: &v1.TCPSocketAction{
			Port: intstr.FromString(portName),
		},
	}
	return &v1.Probe{
		Handler:             *handler,
		InitialDelaySeconds: 10,
		PeriodSeconds:       periodSeconds,
	}
}
