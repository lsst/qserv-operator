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
func GetLabels(component constants.ComponentName, cr_name string) map[string]string {
	return generateComponentLabels(component, cr_name)
}

func generateComponentLabels(component constants.ComponentName, cr_name string) map[string]string {
	componentStr := string(component)
	return map[string]string{
		"app":       constants.AppLabel,
		"component": componentStr,
		"instance":  cr_name,
	}
}

func GetContainerLabels(container constants.ContainerName, cr_name string) map[string]string {
	return generateContainerLabels(container, cr_name)
}

func generateContainerLabels(container constants.ContainerName, cr_name string) map[string]string {
	return map[string]string{
		"app":       constants.AppLabel,
		"container": string(container),
		"instance":  cr_name,
	}
}
