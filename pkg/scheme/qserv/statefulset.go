package qserv

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	qservv1alpha1 "github.com/lsst/qserv-operator/pkg/apis/qserv/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
	"github.com/lsst/qserv-operator/pkg/util"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("qserv")

// GenerateCzarStatefulSet generate statefulset specification for Qserv Czar
func GenerateCzarStatefulSet(cr *qservv1alpha1.Qserv, labels map[string]string) *appsv1.StatefulSet {
	name := cr.Name + "-" + string(constants.CzarName)
	namespace := cr.Namespace
	labels = util.MergeLabels(labels, util.GetLabels(constants.CzarName, cr.Name))

	var replicas int32 = 1
	storageClass := cr.Spec.StorageClass
	storageSize := cr.Spec.StorageCapacity

	initContainer, initVolumes := getInitContainer(cr, constants.CzarName)
	mariadbContainer, mariadbVolumes := getMariadbContainer(cr, constants.CzarName)
	proxyContainer, proxyVolumes := getProxyContainer(cr)
	wmgrContainer, wmgrVolumes := getWmgrContainer(cr)

	var volumes VolumeSet
	volumes.make(initVolumes, mariadbVolumes, proxyVolumes, wmgrVolumes)

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
					InitContainers: []v1.Container{
						initContainer,
					},
					Containers: []v1.Container{
						mariadbContainer,
						proxyContainer,
						wmgrContainer,
					},
					Volumes: volumes.toSlice(),
				},
			},
			VolumeClaimTemplates: []v1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: GetVolumeClaimTemplateName(),
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

func GenerateReplicationCtlStatefulSet(cr *qservv1alpha1.Qserv, labels map[string]string) *appsv1.StatefulSet {
	name := cr.Name + "-" + string(constants.ReplCtlName)
	namespace := cr.Namespace

	labels = util.MergeLabels(labels, util.GetLabels(constants.ReplName, cr.Name))

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
					Containers: []v1.Container{
						replCtlContainer,
					},
					Volumes: volumes.toSlice(),
				},
			},
		},
	}

	ss.Spec.Template.Spec.Tolerations = cr.Spec.Tolerations

	return ss
}

func GenerateReplicationDbStatefulSet(cr *qservv1alpha1.Qserv, labels map[string]string) *appsv1.StatefulSet {
	name := cr.Name + "-" + string(constants.ReplDbName)
	namespace := cr.Namespace

	labels = util.MergeLabels(labels, util.GetLabels(constants.ReplName, cr.Name))

	var replicas int32 = 1
	storageClass := cr.Spec.StorageClass
	storageSize := cr.Spec.StorageCapacity

	initContainer, initVolumes := getInitContainer(cr, constants.ReplName)
	mariadbContainer, mariadbVolumes := getMariadbContainer(cr, constants.ReplName)

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
					InitContainers: []v1.Container{
						initContainer,
					},
					Containers: []v1.Container{
						mariadbContainer,
					},
					Volumes: volumes.toSlice(),
				},
			},
			VolumeClaimTemplates: []v1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: GetVolumeClaimTemplateName(),
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

func GetVolumeClaimTemplateName() string {
	return constants.QservName + "-data"
}

func GenerateWorkerStatefulSet(cr *qservv1alpha1.Qserv, labels map[string]string) *appsv1.StatefulSet {
	name := cr.Name + "-" + string(constants.WorkerName)
	namespace := cr.Namespace

	const (
		MARIADB = iota
		WMGR
	)

	const INIT = 0

	labels = util.MergeLabels(labels, util.GetLabels(constants.WorkerName, cr.Name))

	replicas := cr.Spec.Worker.Replicas
	storageClass := cr.Spec.StorageClass
	storageSize := cr.Spec.StorageCapacity

	initContainer, initVolumes := getInitContainer(cr, constants.WorkerName)
	mariadbContainer, mariadbVolumes := getMariadbContainer(cr, constants.WorkerName)
	xrootdContainers, xrootdVolumes := getXrootdContainers(cr, constants.WorkerName)
	wmgrContainer, wmgrVolumes := getWmgrContainer(cr)
	replicationWrkContainer, replicationWrkVolumes := getReplicationWrkContainer(cr)

	// Volumes
	var volumes VolumeSet
	volumes.make(initVolumes, mariadbVolumes, replicationWrkVolumes, wmgrVolumes, xrootdVolumes)

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
					InitContainers: []v1.Container{
						initContainer,
					},
					Containers: []v1.Container{
						mariadbContainer,
						replicationWrkContainer,
						wmgrContainer,
						xrootdContainers[0],
						xrootdContainers[1],
					},
					Volumes: volumes.toSlice(),
				},
			},
			VolumeClaimTemplates: []v1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: GetVolumeClaimTemplateName(),
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

	ss.Spec.Template.Spec.Tolerations = cr.Spec.Tolerations

	return ss
}

func GenerateXrootdStatefulSet(cr *qservv1alpha1.Qserv, labels map[string]string) *appsv1.StatefulSet {
	namespace := cr.Namespace
	name := util.GetName(cr, string(constants.XrootdRedirectorName))

	labels = util.MergeLabels(labels, util.GetLabels(constants.XrootdRedirectorName, cr.Name))

	var replicas int32 = 2

	containers, volumes := getXrootdContainers(cr, constants.XrootdRedirectorName)

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
					Containers: containers,
					Volumes:    volumes.toSlice(),
				},
			},
		},
	}

	ss.Spec.Template.Spec.Tolerations = cr.Spec.Tolerations

	return ss
}
