package util

import "text/template"

var TemplateFunctions = template.FuncMap{
	"Iterate": IterateCount,
}

func IterateCount(count int) []int {
	items := make([]int, count)
	for i := 0; i < count; i++ {
		items[i] = i
	}
	return items
}
