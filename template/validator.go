package template

import lib "github.com/zokypesch/protoc-gen-generator/lib"

var tmplCoreValidation = `package core

import (
	"fmt"
	"strings"
	"strconv"
)

func Validate(types string, value interface{}, fieldName string) error {
	s := strings.Split(types, "*")
	validation := s[0]
	paramValidation := ""
	if len(s) == 2 {
		paramValidation = s[1]
	}

	switch(validation){
	case "not_empty_string":
		if value.(string) == "" {
			return fmt.Errorf("cannot empty string")
		}
	case "min_max":
		paramsMinMax := strings.Split(paramValidation, ".")
		if len(paramsMinMax) {{unescape "<"}} 2 {
			return fmt.Errorf("wrong number param")
		}
		min, _ := strconv.Atoi(paramsMinMax[0])
		max, _ := strconv.Atoi(paramsMinMax[1])
		if len(value.(string)) {{unescape "<"}} min {
			return fmt.Errorf("length " + fieldName + " less than " + strconv.Itoa(min))
		}

		if len(value.(string)) {{unescape ">"}} max {
			return fmt.Errorf("length " + fieldName + " more than " + strconv.Itoa(max))
		}
	}
	return nil
}
`

var ListCoreValidation = lib.List{
	FileType: ".validation.go",
	Template: tmplCoreValidation,
	Location: "./core/",
	Lang:     "go",
}
