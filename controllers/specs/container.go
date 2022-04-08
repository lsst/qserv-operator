package specs

import (
	"fmt"
	"path/filepath"

	"github.com/go-logr/logr"
	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/controllers/constants"
	"github.com/lsst/qserv-operator/controllers/util"
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
		Image:           cr.Spec.DbImage,
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
		Image:           cr.Spec.DbImage,
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

func getProxyContainer(cr *qservv1beta1.Qserv) (v1.Container, VolumeSet) {
	spec := cr.Spec

	volumeMounts := []v1.VolumeMount{
		// Used for mysql socket access
		// TODO move mysql socket in emptyDir?
		getDataVolumeMount(),
		getStartVolumeMount(constants.ProxyName),
	}

	// Volumes
	var volumes InstanceVolumeSet
	volumes.make(cr)

	setCorePath(spec.Devel.CorePath, &volumeMounts, &volumes)

	container := v1.Container{
		Command:         constants.Command,
		Image:           spec.Image,
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

	volumes.addStartVolume(constants.ProxyName)

	return container, volumes.volumeSet
}

func getReplicationCtlContainer(cr *qservv1beta1.Qserv) (v1.Container, VolumeSet) {
	spec := cr.Spec

	var probeTimeoutSeconds int32 = 3
	volumeMounts := []v1.VolumeMount{
		getStartVolumeMount(constants.ReplCtlName),
		getSecretVolumeMount(constants.ReplDbName),
		getSecretVolumeMount(constants.MariadbName),
	}

	var volumes InstanceVolumeSet
	volumes.make(cr)

	setCorePath(spec.Devel.CorePath, &volumeMounts, &volumes)

	container := v1.Container{
		Command:         constants.Command,
		LivenessProbe:   getHTTPProbe(constants.HTTPPortName, 10, probeTimeoutSeconds, "meta/version"),
		ReadinessProbe:  getHTTPProbe(constants.HTTPPortName, 5, probeTimeoutSeconds, "meta/version"),
		Name:            string(constants.ReplCtlName),
		Image:           spec.Image,
		ImagePullPolicy: spec.ImagePullPolicy,
		Ports: []v1.ContainerPort{
			{
				Name:          constants.HTTPPortName,
				ContainerPort: constants.HTTPPort,
				Protocol:      v1.ProtocolTCP,
			},
		},
		VolumeMounts: volumeMounts,
	}

	reqLogger := log.WithValues("Request.Namespace", cr.Namespace, "Request.Name", cr.Name)

	reqLogger.Info(fmt.Sprintf("Debug level for replication controller: %s", spec.Replication.Debug))

	volumes.addStartVolume(constants.ReplCtlName)
	volumes.addSecretVolume(constants.ReplDbName)
	volumes.addSecretVolume(constants.MariadbName)

	setDebug(spec.Replication.Debug, constants.ReplCtlName, &container)
	return container, volumes.volumeSet
}

func getReplicationRegistryContainer(cr *qservv1beta1.Qserv) (v1.Container, VolumeSet) {
	spec := cr.Spec

	var probeTimeoutSeconds int32 = 3
	volumeMounts := []v1.VolumeMount{
		getStartVolumeMount(constants.ReplRegistryName),
		getSecretVolumeMount(constants.ReplDbName),
		getSecretVolumeMount(constants.MariadbName),
	}

	var volumes InstanceVolumeSet
	volumes.make(cr)

	setCorePath(spec.Devel.CorePath, &volumeMounts, &volumes)

	container := v1.Container{
		Command:         constants.Command,
		LivenessProbe:   getHTTPProbe(constants.HTTPPortName, 10, probeTimeoutSeconds, "workers"),
		ReadinessProbe:  getHTTPProbe(constants.HTTPPortName, 5, probeTimeoutSeconds, "workers"),
		Name:            string(constants.ReplRegistryName),
		Image:           spec.Image,
		ImagePullPolicy: spec.ImagePullPolicy,
		Ports: []v1.ContainerPort{
			{
				Name:          constants.HTTPPortName,
				ContainerPort: constants.HTTPPort,
				Protocol:      v1.ProtocolTCP,
			},
		},
		VolumeMounts: volumeMounts,
	}

	reqLogger := log.WithValues("Request.Namespace", cr.Namespace, "Request.Name", cr.Name)

	reqLogger.Info(fmt.Sprintf("Debug level for replication controller: %s", spec.Replication.Debug))

	volumes.addStartVolume(constants.ReplRegistryName)
	volumes.addSecretVolume(constants.ReplDbName)
	volumes.addSecretVolume(constants.MariadbName)

	setDebug(spec.Replication.Debug, constants.ReplCtlName, &container)
	return container, volumes.volumeSet
}

func getReplicationWrkContainer(cr *qservv1beta1.Qserv) (v1.Container, VolumeSet) {
	spec := cr.Spec

	volumeMounts := []v1.VolumeMount{
		getDataVolumeMount(),
		getStartVolumeMount(constants.ReplWrkName),
		getSecretVolumeMount(constants.MariadbName),
		getSecretVolumeMount(constants.ReplDbName),
	}

	var volumes InstanceVolumeSet
	volumes.make(cr)

	setCorePath(spec.Devel.CorePath, &volumeMounts, &volumes)

	container := v1.Container{
		Name:            string(constants.ReplWrkName),
		Image:           spec.Image,
		ImagePullPolicy: spec.ImagePullPolicy,
		Resources:       cr.Spec.Worker.ReplicationResources,
		Command:         constants.Command,
		// TODO add ports
		SecurityContext: &v1.SecurityContext{
			RunAsUser: &constants.QservUID,
		},
		VolumeMounts: volumeMounts,
	}

	volumes.addStartVolume(constants.ReplWrkName)
	volumes.addSecretVolume(constants.MariadbName)
	volumes.addSecretVolume(constants.ReplDbName)

	setDebug(spec.Replication.Debug, constants.ReplWrkName, &container)
	return container, volumes.volumeSet
}

func getXrootdContainers(cr *qservv1beta1.Qserv, component constants.PodClass) ([]v1.Container, VolumeSet) {

	spec := cr.Spec

	var cmsdVolumeMounts []v1.VolumeMount
	var xrootdVolumeMounts []v1.VolumeMount
	var cmsdContainerName constants.ContainerName
	var xrootdContainerName constants.ContainerName

	// Volumes
	var volumes InstanceVolumeSet
	volumes.make(cr)

	if component == constants.XrootdRedirector {
		cmsdVolumeMounts = getXrootdVolumeMounts(constants.CmsdRedirectorName)
		xrootdVolumeMounts = getXrootdVolumeMounts(constants.XrootdRedirectorName)
		cmsdContainerName = constants.CmsdRedirectorName
		xrootdContainerName = constants.XrootdRedirectorName
	} else if component == constants.Worker {
		cmsdVolumeMounts = getXrootdVolumeMounts(constants.CmsdServerName)
		xrootdVolumeMounts = getXrootdVolumeMounts(constants.XrootdServerName)
		cmsdContainerName = constants.CmsdServerName
		xrootdContainerName = constants.XrootdServerName
	} else {
		var noVolume map[string]v1.Volume
		return []v1.Container{}, noVolume
	}

	setCorePath(spec.Devel.CorePath, &cmsdVolumeMounts, &volumes)
	setCorePath(spec.Devel.CorePath, &xrootdVolumeMounts, &volumes)

	containers := []v1.Container{
		{
			Name:            string(cmsdContainerName),
			Image:           spec.Image,
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
			VolumeMounts: cmsdVolumeMounts,
		},
		{
			Name:            string(xrootdContainerName),
			Image:           spec.Image,
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
			VolumeMounts: xrootdVolumeMounts,
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

	volumes.addStartVolume(cmsdContainerName)
	volumes.addStartVolume(xrootdContainerName)
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
		InitialDelaySeconds: constants.ProbeInitialDelaySeconds,
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

// updateContainersImages update container image field with "image" in "containers" list, only for containers whose name are in containersNames
// return true is image has been updated, else false
func updateContainersImages(qserv *qservv1beta1.Qserv, containers []v1.Container) bool {
	hasUpdate := false
	for i := range containers {
		if util.HasValue(containers[i].Name, constants.WithQservImage) && containers[i].Image != qserv.Spec.Image {
			containers[i].Image = qserv.Spec.Image
			hasUpdate = true
		} else if util.HasValue(containers[i].Name, constants.WithMariadbImage) && containers[i].Image != qserv.Spec.DbImage {
			containers[i].Image = qserv.Spec.DbImage
			hasUpdate = true
		}
	}
	return hasUpdate
}
