package qserv

import (
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	qservv1alpha1 "github.com/lsst/qserv-operator/pkg/apis/qserv/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
	"github.com/lsst/qserv-operator/pkg/util"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("qserv")

func GenerateReplicationDbStatefulSet(cr *qservv1alpha1.Qserv, labels map[string]string) *appsv1beta2.StatefulSet {
	name := cr.Name + "-" + constants.ReplDbName
	namespace := cr.Namespace

	labels = util.MergeLabels(labels, util.GetLabels(constants.ReplDbName, cr.Name))

	var replicas int32 = 1
	storageClass := "standard"
	storageSize := "1G"

	initContainer, initVolumes := getInitContainer(cr, constants.ReplName)
	mariadbContainer, mariadbVolumes := getMariadbContainer(cr, constants.ReplName)

	var volumes VolumeSet
	volumes.make(initVolumes, mariadbVolumes)

	ss := &appsv1beta2.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1beta2.StatefulSetSpec{
			ServiceName: name,
			Replicas:    &replicas,
			UpdateStrategy: appsv1beta2.StatefulSetUpdateStrategy{
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
						Name: getPersistentVolumeClaimName(constants.QservName),
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

	return ss
}

func getPersistentVolumeClaimName(componentName string) string {
	return componentName + "-data"
}

func GenerateCzarStatefulSet(cr *qservv1alpha1.Qserv, labels map[string]string) *appsv1beta2.StatefulSet {
	name := cr.Name + "-" + constants.CzarName
	namespace := cr.Namespace

	labels = util.MergeLabels(labels, util.GetLabels(constants.CzarName, cr.Name))

	var replicas int32 = 1
	storageClass := "standard"
	storageSize := "1G"

	initContainer, initVolumes := getInitContainer(cr, constants.CzarName)
	mariadbContainer, mariadbVolumes := getMariadbContainer(cr, constants.CzarName)
	proxyContainer, proxyVolumes := getProxyContainer(cr)
	wmgrContainer, wmgrVolumes := getWmgrContainer(cr)

	var volumes VolumeSet
	volumes.make(initVolumes, mariadbVolumes, proxyVolumes, wmgrVolumes)

	ss := &appsv1beta2.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1beta2.StatefulSetSpec{
			ServiceName: name,
			Replicas:    &replicas,
			UpdateStrategy: appsv1beta2.StatefulSetUpdateStrategy{
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
						Name: getPersistentVolumeClaimName(constants.QservName),
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

	return ss
}

func GenerateWorkerStatefulSet(cr *qservv1alpha1.Qserv, labels map[string]string) *appsv1beta2.StatefulSet {
	name := cr.Name + "-" + constants.WorkerName
	namespace := cr.Namespace

	const (
		MARIADB = iota
		WMGR
	)

	const INIT = 0

	labels = util.MergeLabels(labels, util.GetLabels(constants.WorkerName, cr.Name))

	var replicas int32 = 2
	storageClass := "standard"
	storageSize := "1G"

	initContainer, initVolumes := getInitContainer(cr, constants.WorkerName)
	mariadbContainer, mariadbVolumes := getMariadbContainer(cr, constants.WorkerName)
	xrootdContainers, xrootdVolumes := getXrootdContainers(cr, constants.WorkerName)
	wmgrContainer, wmgrVolumes := getWmgrContainer(cr)

	// Volumes
	var volumes VolumeSet
	volumes.make(initVolumes, mariadbVolumes, wmgrVolumes, xrootdVolumes)

	ss := &appsv1beta2.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1beta2.StatefulSetSpec{
			PodManagementPolicy: "Parallel",
			ServiceName:         name,
			Replicas:            &replicas,
			UpdateStrategy: appsv1beta2.StatefulSetUpdateStrategy{
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
						Name: getPersistentVolumeClaimName(constants.QservName),
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

	return ss
}

func GenerateXrootdStatefulSet(cr *qservv1alpha1.Qserv, labels map[string]string) *appsv1beta2.StatefulSet {
	namespace := cr.Namespace
	name := util.GetXrootdRedirectorName(cr)

	labels = util.MergeLabels(labels, util.GetLabels(constants.XrootdRedirectorName, cr.Name))

	var replicas int32 = 2

	containers, volumes := getXrootdContainers(cr, constants.XrootdRedirectorName)

	ss := &appsv1beta2.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1beta2.StatefulSetSpec{
			PodManagementPolicy: "Parallel",
			ServiceName:         name,
			Replicas:            &replicas,
			UpdateStrategy: appsv1beta2.StatefulSetUpdateStrategy{
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

	return ss
}
