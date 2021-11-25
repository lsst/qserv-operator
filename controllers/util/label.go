package util

import (
	"github.com/lsst/qserv-operator/controllers/constants"
)

// GetInstanceLabels returns the labels for a Qserv instance
func GetInstanceLabels(crName string) map[string]string {
	return map[string]string{
		"app":                          constants.QservName,
		"app.kubernetes.io/managed-by": "qserv-operator",
		"instance":                     crName,
	}
}

// GetComponentLabels returns the labels for the component with specific role
func GetComponentLabels(component constants.PodClass, crName string) map[string]string {
	labels := GetInstanceLabels(crName)
	labels["component"] = string(component)
	return labels
}

// GetContainerLabels returns the labels for containers
func GetContainerLabels(container constants.ContainerName, crName string) map[string]string {
	labels := GetInstanceLabels(crName)
	labels["component"] = string(container)
	return labels
}
