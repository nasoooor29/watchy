package utils

import (
	"strconv"
	"text/template"
)

func toFloat(v string) float64 {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0
	}

	return f
}

func GetFuncMap() template.FuncMap {
	return template.FuncMap{
		"float": toFloat,
	}
}
