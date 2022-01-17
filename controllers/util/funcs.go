package util

import (
	"text/template"

	"github.com/lsst/qserv-operator/controllers/constants"
)

// TemplateFunctions contain functions used in templates for Qserv configuration files
var TemplateFunctions = template.FuncMap{
	"Iterate": IterateCount,
}

// IterateCount return a list of integer in the  for [0, 1, ..., n]
func IterateCount(count uint) []int {
	items := make([]int, count)
	for i := uint(0); i < count; i++ {
		items[i] = int(i)
	}
	return items
}

func GetValue(value string, defaultValue string) string {
	if value == "" {
		value = defaultValue
	}
	return value
}

func HasValue(value string, slice []constants.ContainerName) bool {
	for _, v := range slice {
		if value == string(v) {
			return true
		}
	}
	return false
}
