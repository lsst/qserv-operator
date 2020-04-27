package util

import "text/template"

var TemplateFunctions = template.FuncMap{
	"Iterate": IterateCount,
}

func IterateCount(count uint) []int {
	items := make([]int, count)
	for i := uint(0); i < count; i++ {
		items[i] = int(i)
	}
	return items
}
