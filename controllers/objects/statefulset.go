package objects

import (
	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/controllers/constants"
	"github.com/lsst/qserv-operator/controllers/util"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
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

// GenerateXrootdStatefulSet generate statefulset specification for xrootd redirectors
func GenerateXrootdStatefulSet(cr *qservv1beta1.Qserv) *appsv1.StatefulSet {
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
