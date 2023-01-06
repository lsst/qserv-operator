package util

import (
	"text/template"

	"github.com/lsst/qserv-operator/controllers/constants"
)

// TemplateFunctions contain functions used in templates for Qserv configuration files
var TemplateFunctions = template.FuncMap{
	"Iterate":           IterateCount,
	"WorkerDatabaseUrl": WorkerDatabaseURL,
}

// IterateCount return a list of integer in the  for [0, 1, ..., n]
func IterateCount(count uint) []int {
	items := make([]int, count)
	for i := uint(0); i < count; i++ {
		items[i] = int(i)
	}
	return items
}

// GetPtrValue return value if not nil
// else set *value to *defaultValue if defaultValue is not nil
func GetPtrValue(value *string, defaultValue *string) *string {
	if value == nil && defaultValue != nil {
		*value = *defaultValue
	}
	return value
}

// GetValue return value if not empty, else return defaultValue
func GetValue(value string, defaultValue string) string {
	if value == "" {
		value = defaultValue
	}
	return value
}

// HasValue return true if value is in slice, else return false
// FIXME: See assert.Contains for better implementation
func HasValue(value string, slice []constants.ContainerName) bool {
	for _, v := range slice {
		if value == string(v) {
			return true
		}
	}
	return false
}
