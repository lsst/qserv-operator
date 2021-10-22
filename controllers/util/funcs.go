package util

import "text/template"

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

// NoFunc does nothing and returns nil
var NoFunc = func() error {
	return nil
}
