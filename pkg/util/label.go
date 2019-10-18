package util

import (
	"github.com/lsst/qserv-operator/pkg/constants"
)

// MergeLabels merges all the label maps received as argument into a single new label map.
func MergeLabels(allLabels ...map[string]string) map[string]string {
	res := map[string]string{}

	for _, labels := range allLabels {
		if labels != nil {
			for k, v := range labels {
				res[k] = v
			}
		}
	}
	return res
}

// GetLabels returns the labels for the component with specific role
func GetLabels(component constants.ComponentName, role string) map[string]string {
	return generateComponentLabels(component, role)
}

func generateComponentLabels(component constants.ComponentName, role string) map[string]string {
	componentStr := string(component)
	return map[string]string{
		"app":       constants.AppLabel,
		"component": componentStr,
		"instance":  role,
	}
}

func GetContainerLabels(container constants.ContainerName, role string) map[string]string {
	return generateContainerLabels(container, role)
}

func generateContainerLabels(container constants.ContainerName, role string) map[string]string {
	return map[string]string{
		"app":       constants.AppLabel,
		"container": string(container),
		"instance":  role,
	}
}
