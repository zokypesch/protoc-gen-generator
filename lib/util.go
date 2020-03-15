package lib

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode"
)

func grpcTypeToTs(param string) string {
	switch param {
	case "TYPE_STRING":
		return "string"
	case "TYPE_INT64":
	case "TYPE_INT32":
		return "number"
	case "TYPE_BOOL":
		return "boolean"
	default:
		return "string"
	}

	return "string"
}

func grpcTypeToGo(param string) string {
	switch param {
	case "TYPE_STRING":
		return "string"
	case "TYPE_INT64":
		return "int64"
	case "TYPE_INT32":
		return "int32"
	case "TYPE_BOOL":
		return "bool"
	case "TYPE_DOUBLE":
		return "float64"
	case "TYPE_FLOAT":
		return "float32"
	default:
		return "string"
	}
}

func ucFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func underscoreToGoFormat(s string) string {
	splitV := strings.Split(s, "_")
	result := ""

	for _, v := range splitV {
		result += ucFirst(v)
	}

	return result
}

func ucDown(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

func stringToOpt(param string) []*Option {

	onOpt := strings.Split(param, " ")
	var newOptions []*Option

	// log.Println(param, onOpt)

	for _, vOnOpt := range onOpt {
		splitV := strings.Split(vOnOpt, ":")
		if len(splitV) < 2 {
			// log.Println("skip", vOnOpt)
			continue
		}
		// log.Println(splitV[0], cleanQuote(splitV[1]))
		newOptions = append(newOptions, &Option{
			Code:  splitV[0],
			Name:  splitV[0],
			Value: cleanQuote(splitV[1]),
		})
	}
	return newOptions
}

func getStringFromOptCode(param string) string {
	switch param {
	case "50056":
		return "httpMode"
	case "50057":
		return "isRepository"
	case "72295728":
		return "urlPath"
	case "50058":
		return "ignoreFieldDb"
	case "50059":
		return "isPrimaryKey"
	case "50060":
		return "requiredField"
	case "50061":
		return "requiredType"
	case "50062":
		return "agregator"
	case "50063":
		return "fulltext"
	case "50064":
		return "elastic"
	case "50065":
		return "errorDesc"
	case "50066":
		return "foreignKey"
	case "50067":
		return "whitelist"
	}
	return ""
}

var listStringOfMethod = map[string]string{
	`\x12\x10`: "get",
	`\\x10`:    "post",
	`\x1a\x10`: "put",
	`*\x10`:    "delete",

	`\x12\x0f`: "get",
	`\\x0f`:    "post",
	`\x1a\x0f`: "put",
	`*\x0f`:    "delete",
}

func getHttpModeWithUrl(param string) (string, string) {
	for k, v := range listStringOfMethod {
		if strings.Contains(param, k) {
			r1 := strings.Replace(param, k, "", -1)
			return r1, v
		}
	}
	return "get", ""
}

func getHttpUrl(param string) string {
	i := strings.Index(param, `/`)

	if i == -1 {
		return ""
	}

	chars := param[i:]
	return chars

}

func cleanQuote(param string) string {
	return strings.Replace(param, `"`, "", -1)
}

// urlPath = strings.Replace(vOpt.Value, `\x12\x10`, "", -1)
// urlPath = strings.Replace(urlPath, `\"\x10`, "", -1)
// urlPath = strings.Replace(urlPath, `\x1a\x10`, "", -1)
// urlPath = strings.Replace(urlPath, `*\x10`, "", -1)

func unescape(s string) template.HTML {
	return template.HTML(s)
}

func getFirstService(param []Service) Service {
	if len(param) == 0 {
		return Service{}
	}
	return param[0]
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

var camel = regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)")

func underscore(s string) string {
	var a []string
	for _, sub := range camel.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}
		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}
	return strings.ToLower(strings.Join(a, "_"))
}
