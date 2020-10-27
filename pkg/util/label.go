package util

import (
	"github.com/lsst/qserv-operator/pkg/constants"
)

// MergeLabels merges all the label maps received as argument into a single new label map.
func MergeLabels(allLabels ...map[string]string) map[string]string {
	res := map[string]string{}

	for _, labels := range allLabels {
		for k, v := range labels {
			res[k] = v
		}
	}
	return res
}

// GetLabels returns the labels for the component with specific role
func GetLabels(component constants.PodClass, crName string) map[string]string {
	return generatePodLabels(component, crName)
}

func generatePodLabels(component constants.PodClass, crName string) map[string]string {
	componentStr := string(component)
	return map[string]string{
		"app":       constants.AppLabel,
		"component": componentStr,
		"instance":  crName,
	}
}

// GetContainerLabels returns the labels for containers
func GetContainerLabels(container constants.ContainerName, crName string) map[string]string {
	return generateContainerLabels(container, crName)
}

func generateContainerLabels(container constants.ContainerName, crName string) map[string]string {
	return map[string]string{
		"app":       constants.AppLabel,
		"container": string(container),
		"instance":  crName,
	}
}
