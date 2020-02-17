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

func getInitContainer(cr *qservv1alpha1.Qserv, component constants.ComponentName) (v1.Container, VolumeSet) {
	componentName := string(component)

	sqlConfigSuffix := fmt.Sprintf("sql-%s", component)

	var dbName constants.ContainerName
	if component == constants.ReplName {
		dbName = constants.ReplDbName
	} else {
		dbName = constants.MariadbName
	}

	container := v1.Container{
		Name:  string(constants.InitDbName),
		Image: getMariadbImage(cr, component),
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
			getEtcVolumeMount(dbName),
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
	volumes.addEtcVolume(dbName)
	volumes.addStartVolume(constants.InitDbName)
	volumes.addSecretVolume(constants.MariadbName)

	if dbName == constants.ReplDbName {
		container.VolumeMounts = append(container.VolumeMounts, getSecretVolumeMount(dbName))
		volumes.addSecretVolume(dbName)
	}

	return container, volumes.volumeSet
}

func getMariadbContainer(cr *qservv1alpha1.Qserv, component constants.ComponentName) (v1.Container, VolumeSet) {

	var uservice constants.ContainerName
	if component == constants.ReplName {
		uservice = constants.ReplDbName
	} else {
		uservice = constants.MariadbName
	}

	mariadbPortName := string(constants.MariadbName)

	container := v1.Container{
		Name:  string(constants.MariadbName),
		Image: getMariadbImage(cr, component),
		Ports: []v1.ContainerPort{
			{
				Name:          mariadbPortName,
				ContainerPort: constants.MariadbPort,
				Protocol:      v1.ProtocolTCP,
			},
		},
		Command:        constants.Command,
		LivenessProbe:  getProbe(constants.MariadbPortName, 10, tcpAction),
		ReadinessProbe: getProbe(constants.MariadbPortName, 5, tcpAction),
		VolumeMounts: []v1.VolumeMount{
			getDataVolumeMount(),
			getEtcVolumeMount(uservice),
			getStartVolumeMount(uservice),
			getTmpVolumeMount(),
		},
	}

	// Volumes
	var volumes InstanceVolumeSet
	volumes.make(cr)

	volumes.addEmptyDirVolume("tmp-volume")
	volumes.addEtcStartVolumes(uservice)

	return container, volumes.volumeSet
}

func getMariadbImage(cr *qservv1alpha1.Qserv, component constants.ComponentName) string {
	spec := cr.Spec
	var image string
	if component == constants.ReplName {
		image = spec.Replication.DbImage
	} else if component == constants.WorkerName {
		image = spec.Worker.Image
	} else if component == constants.CzarName {
		image = spec.Czar.Image
	}
	return image
}

func getProxyContainer(cr *qservv1alpha1.Qserv) (v1.Container, VolumeSet) {
	spec := cr.Spec
	container := v1.Container{
		Name:  string(constants.ProxyName),
		Image: spec.Czar.Image,
		Ports: []v1.ContainerPort{
			{
				Name:          string(constants.ProxyName),
				ContainerPort: constants.ProxyPort,
				Protocol:      v1.ProtocolTCP,
			},
		},
		LivenessProbe:  getProbe(constants.ProxyPortName, 10, tcpAction),
		ReadinessProbe: getProbe(constants.ProxyPortName, 5, tcpAction),
		Command:        constants.Command,
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

	container := v1.Container{
		Name:    string(constants.ReplCtlName),
		Image:   spec.Replication.Image,
		Command: constants.Command,
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
		VolumeMounts: []v1.VolumeMount{
			getEtcVolumeMount(constants.ReplCtlName),
			getStartVolumeMount(constants.ReplCtlName),
			getSecretVolumeMount(constants.ReplDbName),
			getSecretVolumeMount(constants.MariadbName),
		},
	}

	var volumes InstanceVolumeSet
	volumes.make(cr)

	volumes.addEtcStartVolumes(constants.ReplCtlName)
	volumes.addSecretVolume(constants.ReplDbName)
	volumes.addSecretVolume(constants.MariadbName)

	return container, volumes.volumeSet
}

func getReplicationWrkContainer(cr *qservv1alpha1.Qserv) (v1.Container, VolumeSet) {
	spec := cr.Spec

	var runAsUser int64 = 1000

	container := v1.Container{
		Name:    string(constants.ReplWrkName),
		Image:   spec.Replication.Image,
		Command: constants.Command,
		Env: []v1.EnvVar{
			{
				Name:  "REPL_DB_DN",
				Value: util.GetName(cr, string(constants.ReplDbName)),
			},
		},
		SecurityContext: &v1.SecurityContext{
			RunAsUser: &runAsUser,
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
		Name:  string(constants.WmgrName),
		Image: cr.Spec.Worker.Image,
		Ports: []v1.ContainerPort{
			{
				Name:          constants.WmgrPortName,
				ContainerPort: constants.WmgrPort,
				Protocol:      v1.ProtocolTCP,
			},
		},
		Command:        constants.Command,
		LivenessProbe:  getProbe(constants.WmgrPortName, 10, tcpAction),
		ReadinessProbe: getProbe(constants.WmgrPortName, 5, tcpAction),
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

func getXrootdContainers(cr *qservv1alpha1.Qserv, component constants.ComponentName) ([]v1.Container, VolumeSet) {

	const (
		CMSD = iota
		XROOTD
	)

	spec := cr.Spec

	volumeMounts := getXrootdVolumeMounts(component)
	xrootdPortName := string(constants.XrootdName)

	containers := []v1.Container{
		{
			Name:    string(constants.CmsdName),
			Image:   spec.Worker.Image,
			Command: constants.Command,
			Args:    []string{"-S", "cmsd"},
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
			Name:  string(constants.XrootdName),
			Image: spec.Worker.Image,
			Ports: []v1.ContainerPort{
				{
					Name:          xrootdPortName,
					ContainerPort: 1094,
					Protocol:      v1.ProtocolTCP,
				},
			},
			Command:        constants.Command,
			LivenessProbe:  getProbe(constants.XrootdPortName, 10, tcpAction),
			ReadinessProbe: getProbe(constants.XrootdPortName, 5, tcpAction),
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
				Name:          string(constants.CmsdName),
				ContainerPort: 2131,
				Protocol:      v1.ProtocolTCP,
			},
		}
		containers[0].LivenessProbe = getProbe(constants.CmsdPortName, 10, tcpAction)
		containers[0].ReadinessProbe = getProbe(constants.CmsdPortName, 5, tcpAction)
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

func getProbe(portName string, periodSeconds int32, action NetworkAction) *v1.Probe {
	var handler *v1.Handler
	if action == httpAction {
		handler = &v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Port: intstr.FromString(portName),
			},
		}
	} else {
		handler = &v1.Handler{
			TCPSocket: &v1.TCPSocketAction{
				Port: intstr.FromString(portName),
			},
		}
	}
	return &v1.Probe{
		Handler:             *handler,
		InitialDelaySeconds: 10,
		PeriodSeconds:       periodSeconds,
	}
}
