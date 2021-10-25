package qserv

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/go-logr/logr"
	qservv1beta1 "github.com/lsst/qserv-operator/api/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
	"github.com/lsst/qserv-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// addDebuggerContainer perform an in-place update of ss to add a debugger container inside it
func addDebuggerContainer(reqLogger logr.Logger, ss *appsv1.StatefulSet, cr *qservv1beta1.Qserv) {
	if cr.Spec.Devel.EnableDebugger {
		reqLogger.Info("Debugger enabled")
		ss.Spec.Template.Spec.ShareProcessNamespace = &cr.Spec.Devel.EnableDebugger
		ss.Spec.Template.Spec.Containers = append(ss.Spec.Template.Spec.Containers, getDebuggerContainer(reqLogger, cr))
	}
}

func getDebuggerContainer(reqLogger logr.Logger, cr *qservv1beta1.Qserv) v1.Container {

	debuggerContainerName := constants.DebuggerName

	// Container

	reqLogger.Info("Debugger image:", "image", cr.Spec.Devel.DebuggerImage)

	container := v1.Container{
		Image:           cr.Spec.Devel.DebuggerImage,
		ImagePullPolicy: cr.Spec.ImagePullPolicy,
		Name:            string(debuggerContainerName),
		SecurityContext: &v1.SecurityContext{
			Capabilities: &v1.Capabilities{
				Add: []v1.Capability{
					v1.Capability("SYS_PTRACE"),
				},
			},
		},
		Stdin: true,
		TTY:   true,
	}

	return container
}

func getInitContainer(cr *qservv1beta1.Qserv, component constants.PodClass) (v1.Container, VolumeSet) {
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
			getMysqlCnfVolumeMount(dbContainerName),
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

func getMariadbContainer(cr *qservv1beta1.Qserv, pod constants.PodClass) (v1.Container, VolumeSet) {

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
			getMysqlCnfVolumeMount(dbContainerName),
			getStartVolumeMount(dbContainerName),
			getTmpVolumeMount(),
		},
	}

	return container, volumes.volumeSet
}

func getMariadbImage(cr *qservv1beta1.Qserv, component constants.PodClass) string {
	spec := cr.Spec
	var image string
	if component == constants.ReplDb {
		image = spec.Replication.DbImage
	} else if component == constants.IngestDb {
		image = spec.Ingest.DbImage
	} else if component == constants.Worker {
		image = spec.Worker.DbImage
	} else if component == constants.Czar {
		image = spec.Czar.DbImage
	}
	return image
}

func getProxyContainer(cr *qservv1beta1.Qserv) (v1.Container, VolumeSet) {
	spec := cr.Spec

	volumeMounts := []v1.VolumeMount{
		// Used for mysql socket access
		// TODO move mysql socket in emptyDir?
		getDataVolumeMount(),
		getEtcVolumeMount(constants.ProxyName),
		getStartVolumeMount(constants.ProxyName),
	}

	// Volumes
	var volumes InstanceVolumeSet
	volumes.make(cr)

	setCorePath(spec.Devel.CorePath, &volumeMounts, &volumes)

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
		Resources:      spec.Czar.ProxyResources,
		LivenessProbe:  getTCPProbe(constants.ProxyPortName, 10),
		ReadinessProbe: getTCPProbe(constants.ProxyPortName, 5),
		VolumeMounts:   volumeMounts,
	}

	volumes.addEtcStartVolumes(constants.ProxyName)

	return container, volumes.volumeSet
}

func getReplicationCtlContainer(cr *qservv1beta1.Qserv) (v1.Container, VolumeSet) {
	spec := cr.Spec

	var probeTimeoutSeconds int32 = 3
	volumeMounts := []v1.VolumeMount{
		getEtcVolumeMount(constants.ReplCtlName),
		getStartVolumeMount(constants.ReplCtlName),
		getSecretVolumeMount(constants.ReplDbName),
		getSecretVolumeMount(constants.MariadbName),
	}

	var volumes InstanceVolumeSet
	volumes.make(cr)

	setCorePath(spec.Devel.CorePath, &volumeMounts, &volumes)

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
		VolumeMounts: volumeMounts,
	}

	reqLogger := log.WithValues("Request.Namespace", cr.Namespace, "Request.Name", cr.Name)

	reqLogger.Info(fmt.Sprintf("Debug level for replication controller: %s", spec.Replication.Debug))

	volumes.addEtcStartVolumes(constants.ReplCtlName)
	volumes.addSecretVolume(constants.ReplDbName)
	volumes.addSecretVolume(constants.MariadbName)

	setDebug(spec.Replication.Debug, constants.ReplCtlName, &container)
	return container, volumes.volumeSet
}

func getReplicationWrkContainer(cr *qservv1beta1.Qserv) (v1.Container, VolumeSet) {
	spec := cr.Spec

	volumeMounts := []v1.VolumeMount{
		getDataVolumeMount(),
		getEtcVolumeMount(constants.ReplWrkName),
		getStartVolumeMount(constants.ReplWrkName),
		getSecretVolumeMount(constants.MariadbName),
		getSecretVolumeMount(constants.ReplDbName),
	}

	var volumes InstanceVolumeSet
	volumes.make(cr)

	setCorePath(spec.Devel.CorePath, &volumeMounts, &volumes)

	container := v1.Container{
		Name:            string(constants.ReplWrkName),
		Image:           spec.Replication.Image,
		ImagePullPolicy: spec.ImagePullPolicy,
		Resources:       cr.Spec.Worker.ReplicationResources,
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
		// TODO add ports
		SecurityContext: &v1.SecurityContext{
			RunAsUser: &constants.QservUID,
		},
		VolumeMounts: volumeMounts,
	}

	volumes.addEtcStartVolumes(constants.ReplWrkName)
	volumes.addSecretVolume(constants.MariadbName)
	volumes.addSecretVolume(constants.ReplDbName)

	setDebug(spec.Replication.Debug, constants.ReplWrkName, &container)
	return container, volumes.volumeSet
}

func getDashboardContainer(cr *qservv1beta1.Qserv) (v1.Container, VolumeSet) {
	container := v1.Container{
		Name:            string(constants.DashboardName),
		Image:           cr.Spec.Dashboard.Image,
		ImagePullPolicy: cr.Spec.ImagePullPolicy,
		Ports: []v1.ContainerPort{
			{
				Name:          constants.DashboardPortName,
				ContainerPort: constants.DashboardPort,
				Protocol:      v1.ProtocolTCP,
			},
		},
		Command:        constants.Command,
		LivenessProbe:  getTCPProbe(constants.DashboardPortName, 10),
		ReadinessProbe: getTCPProbe(constants.DashboardPortName, 5),
		VolumeMounts: []v1.VolumeMount{
			getEtcVolumeMount(constants.DashboardName),
			getStartVolumeMount(constants.DashboardName),
		},
	}

	// Volumes
	var volumes InstanceVolumeSet
	volumes.make(cr)

	volumes.addEtcStartVolumes(constants.DashboardName)

	return container, volumes.volumeSet
}

func getXrootdContainers(cr *qservv1beta1.Qserv, component constants.PodClass) ([]v1.Container, VolumeSet) {

	spec := cr.Spec

	volumeMounts := getXrootdVolumeMounts(component)
	// Volumes
	var volumes InstanceVolumeSet
	volumes.make(cr)

	setCorePath(spec.Devel.CorePath, &volumeMounts, &volumes)

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

	volumes.addEtcStartVolumes(constants.XrootdName)
	volumes.addEmptyDirVolume(constants.XrootdAdminPathVolumeName)

	return containers, volumes.volumeSet
}

func getHTTPProbe(portName string, periodSeconds int32, timeoutSeconds int32, path string) *v1.Probe {
	handler := &v1.Handler{
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
	handler := &v1.Handler{
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

func setCorePath(corepath string, volumeMounts *[]v1.VolumeMount, volumes *InstanceVolumeSet) {
	if len(corepath) != 0 {
		*volumeMounts = append(*volumeMounts, getCorePathVolumeMount(corepath))
		volumes.addCorePathVolume(corepath)
	}
}

// setDebug allow to start a container in debug mode.
// - change container command to 'sleep infinity'
// - add capability SYS_PTRACE
// - remove probes
func setDebug(debug string, name constants.ContainerName, container *v1.Container) {
	switch debug {
	case string(name), "all":
		container.Command = constants.CommandDebug
		container.LivenessProbe = nil
		container.ReadinessProbe = nil
	}
}
