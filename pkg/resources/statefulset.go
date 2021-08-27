package qserv

import (
	qservv1alpha1 "github.com/lsst/qserv-operator/api/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
	"github.com/lsst/qserv-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("qserv")

func getValue(value string, defaultValue string) string {
	if value == "" {
		value = defaultValue
	}
	return value
}

// GenerateCzarStatefulSet generate statefulset specification for Qserv Czar
func GenerateCzarStatefulSet(cr *qservv1alpha1.Qserv) *appsv1.StatefulSet {
	name := cr.Name + "-" + string(constants.Czar)
	namespace := cr.Namespace
	labels := util.GetComponentLabels(constants.Czar, cr.Name)

	reqLogger := log.WithValues("Request.Namespace", cr.Namespace, "Request.Name", cr.Name)

	storageClass := getValue(cr.Spec.Czar.StorageClass, cr.Spec.StorageClass)
	storageSize := getValue(cr.Spec.Czar.StorageCapacity, cr.Spec.StorageCapacity)

	initContainer, initVolumes := getInitContainer(cr, constants.Czar)
	mariadbContainer, mariadbVolumes := getMariadbContainer(cr, constants.Czar)
	proxyContainer, proxyVolumes := getProxyContainer(cr)

	var volumes VolumeSet
	volumes.make(initVolumes, mariadbVolumes, proxyVolumes)

	ss := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: name,
			Replicas:    &cr.Spec.Czar.Replicas,
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: "RollingUpdate",
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					Affinity: &cr.Spec.Czar.Affinity,
					InitContainers: []v1.Container{
						initContainer,
					},
					Containers: []v1.Container{
						mariadbContainer,
						proxyContainer,
					},
					SecurityContext: &v1.PodSecurityContext{
						FSGroup: &constants.QservGID,
					},
					Volumes: volumes.toSlice(),
				},
			},
			VolumeClaimTemplates: []v1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: constants.DataVolumeClaimTemplateName,
					},
					Spec: v1.PersistentVolumeClaimSpec{
						AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
						StorageClassName: &storageClass,
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								"storage": resource.MustParse(storageSize),
							},
						},
					},
				},
			},
		},
	}

	addDebuggerContainer(reqLogger, ss, cr)

	ss.Spec.Template.Spec.Tolerations = cr.Spec.Tolerations

	return ss
}

// GenerateIngestDbStatefulSet generate statefulset specification for Qserv Ingest Database
func GenerateIngestDbStatefulSet(cr *qservv1alpha1.Qserv) *appsv1.StatefulSet {
	name := cr.Name + "-" + string(constants.IngestDb)
	namespace := cr.Namespace

	labels := util.GetComponentLabels(constants.IngestDb, cr.Name)

	var replicas int32 = 1
	storageClass := cr.Spec.StorageClass
	storageSize := cr.Spec.StorageCapacity

	initContainer, initVolumes := getInitContainer(cr, constants.IngestDb)
	mariadbContainer, mariadbVolumes := getMariadbContainer(cr, constants.IngestDb)

	var volumes VolumeSet
	volumes.make(initVolumes, mariadbVolumes)

	ss := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: name,
			Replicas:    &replicas,
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: "RollingUpdate",
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					Affinity: &cr.Spec.Ingest.Affinity,
					InitContainers: []v1.Container{
						initContainer,
					},
					Containers: []v1.Container{
						mariadbContainer,
					},
					SecurityContext: &v1.PodSecurityContext{
						FSGroup: &constants.QservGID,
					},
					Volumes: volumes.toSlice(),
				},
			},
			VolumeClaimTemplates: []v1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: constants.DataVolumeClaimTemplateName,
					},
					Spec: v1.PersistentVolumeClaimSpec{
						AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
						StorageClassName: &storageClass,
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								"storage": resource.MustParse(storageSize),
							},
						},
					},
				},
			},
		},
	}

	ss.Spec.Template.Spec.Tolerations = cr.Spec.Tolerations

	return ss
}

// GenerateReplicationCtlStatefulSet generate statefulset specification for Qserv Replication Controller
func GenerateReplicationCtlStatefulSet(cr *qservv1alpha1.Qserv) *appsv1.StatefulSet {
	name := cr.Name + "-" + string(constants.ReplCtlName)
	namespace := cr.Namespace

	labels := util.GetComponentLabels(constants.ReplCtl, cr.Name)

	reqLogger := log.WithValues("Request.Namespace", cr.Namespace, "Request.Name", cr.Name)

	var replicas int32 = 1

	replCtlContainer, replCtlVolumes := getReplicationCtlContainer(cr)

	var volumes VolumeSet
	volumes.make(replCtlVolumes)

	ss := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			PodManagementPolicy: "Parallel",
			ServiceName:         name,
			Replicas:            &replicas,
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: "RollingUpdate",
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					Affinity: &cr.Spec.Replication.Affinity,
					Containers: []v1.Container{
						replCtlContainer,
					},
					Volumes: volumes.toSlice(),
				},
			},
		},
	}

	addDebuggerContainer(reqLogger, ss, cr)

	ss.Spec.Template.Spec.Tolerations = cr.Spec.Tolerations

	return ss
}

// GenerateReplicationDbStatefulSet generate statefulset specification for Qserv Replication Database
func GenerateReplicationDbStatefulSet(cr *qservv1alpha1.Qserv) *appsv1.StatefulSet {
	name := cr.Name + "-" + string(constants.ReplDbName)
	namespace := cr.Namespace

	labels := util.GetComponentLabels(constants.ReplDb, cr.Name)

	var replicas int32 = 1
	storageClass := cr.Spec.StorageClass
	storageSize := cr.Spec.StorageCapacity

	initContainer, initVolumes := getInitContainer(cr, constants.ReplDb)
	mariadbContainer, mariadbVolumes := getMariadbContainer(cr, constants.ReplDb)

	var volumes VolumeSet
	volumes.make(initVolumes, mariadbVolumes)

	ss := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: name,
			Replicas:    &replicas,
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: "RollingUpdate",
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					Affinity: &cr.Spec.Replication.Affinity,
					InitContainers: []v1.Container{
						initContainer,
					},
					Containers: []v1.Container{
						mariadbContainer,
					},
					SecurityContext: &v1.PodSecurityContext{
						FSGroup: &constants.QservGID,
					},
					Volumes: volumes.toSlice(),
				},
			},
			VolumeClaimTemplates: []v1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: constants.DataVolumeClaimTemplateName,
					},
					Spec: v1.PersistentVolumeClaimSpec{
						AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
						StorageClassName: &storageClass,
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								"storage": resource.MustParse(storageSize),
							},
						},
					},
				},
			},
		},
	}

	ss.Spec.Template.Spec.Tolerations = cr.Spec.Tolerations

	return ss
}

// GenerateWorkerStatefulSet generate statefulset specification for Qserv Workers
func GenerateWorkerStatefulSet(cr *qservv1alpha1.Qserv) *appsv1.StatefulSet {
	name := cr.Name + "-" + string(constants.Worker)
	namespace := cr.Namespace

	labels := util.GetComponentLabels(constants.Worker, cr.Name)

	reqLogger := log.WithValues("Request.Namespace", cr.Namespace, "Request.Name", cr.Name)

	replicas := cr.Spec.Worker.Replicas

	storageClass := getValue(cr.Spec.Worker.StorageClass, cr.Spec.StorageClass)
	storageSize := getValue(cr.Spec.Worker.StorageCapacity, cr.Spec.StorageCapacity)

	initContainer, initVolumes := getInitContainer(cr, constants.Worker)
	mariadbContainer, mariadbVolumes := getMariadbContainer(cr, constants.Worker)
	xrootdContainers, xrootdVolumes := getXrootdContainers(cr, constants.Worker)
	replicationWrkContainer, replicationWrkVolumes := getReplicationWrkContainer(cr)

	// Volumes
	var volumes VolumeSet
	volumes.make(initVolumes, mariadbVolumes, replicationWrkVolumes, xrootdVolumes)

	ss := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			PodManagementPolicy: "Parallel",
			ServiceName:         name,
			Replicas:            &replicas,
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: "RollingUpdate",
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					Affinity: &cr.Spec.Worker.Affinity,
					InitContainers: []v1.Container{
						initContainer,
					},
					Containers: []v1.Container{
						mariadbContainer,
						replicationWrkContainer,
						xrootdContainers[0],
						xrootdContainers[1],
					},
					SecurityContext: &v1.PodSecurityContext{
						FSGroup: &constants.QservGID,
					},
					Volumes: volumes.toSlice(),
				},
			},
			VolumeClaimTemplates: []v1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: constants.DataVolumeClaimTemplateName,
					},
					Spec: v1.PersistentVolumeClaimSpec{
						AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
						StorageClassName: &storageClass,
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								"storage": resource.MustParse(storageSize),
							},
						},
					},
				},
			},
		}}

	addDebuggerContainer(reqLogger, ss, cr)

	ss.Spec.Template.Spec.Tolerations = cr.Spec.Tolerations

	return ss
}

// GenerateXrootdStatefulSet generate statefulset specification for xrootd redirectors
func GenerateXrootdStatefulSet(cr *qservv1alpha1.Qserv) *appsv1.StatefulSet {
	namespace := cr.Namespace
	name := util.GetName(cr, string(constants.XrootdRedirector))

	reqLogger := log.WithValues("Request.Namespace", cr.Namespace, "Request.Name", cr.Name)

	labels := util.GetComponentLabels(constants.XrootdRedirector, cr.Name)

	var replicas int32 = cr.Spec.Xrootd.Replicas

	containers, volumes := getXrootdContainers(cr, constants.XrootdRedirector)

	ss := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			PodManagementPolicy: "Parallel",
			ServiceName:         name,
			Replicas:            &replicas,
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: "RollingUpdate",
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					Affinity:   &cr.Spec.Xrootd.Affinity,
					Containers: containers,
					Volumes:    volumes.toSlice(),
				},
			},
		},
	}

	addDebuggerContainer(reqLogger, ss, cr)

	ss.Spec.Template.Spec.Tolerations = cr.Spec.Tolerations

	return ss
}
