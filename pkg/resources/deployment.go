package qserv

import (
	qservv1alpha1 "github.com/lsst/qserv-operator/api/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
	"github.com/lsst/qserv-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateDashboardDeployment generate deployment specification for Qserv Web
func GenerateDashboardDeployment(cr *qservv1alpha1.Qserv) *appsv1.Deployment {
	name := cr.Name + "-" + string(constants.Dashboard)
	namespace := cr.Namespace
	labels := util.GetComponentLabels(constants.Dashboard, cr.Name)

	var replicas int32 = 1

	dashboardContainer, dashboardVolumes := getDashboardContainer(cr)

	var volumes VolumeSet
	volumes.make(dashboardVolumes)

	ss := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					Affinity: &cr.Spec.Czar.Affinity,
					Containers: []v1.Container{
						dashboardContainer,
					},
					Volumes: volumes.toSlice(),
				},
			},
		},
	}

	ss.Spec.Template.Spec.Tolerations = cr.Spec.Tolerations

	return ss
}
