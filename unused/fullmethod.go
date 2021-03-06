package unused

// type ERROR_NOT_FOUND string

import (
	"bytes"
	"log"
	"strconv"

	// "fmt"
	"go/format"
	"html/template"
	"path"
	"unicode"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	googlegen "github.com/golang/protobuf/protoc-gen-go/generator"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pkg/errors"
)

const (
	tmpl = `
// Code generated by protoc-gen-go. DO NOT EDIT.
// source: {{ .FileName }}
package {{ .GoPackage }}
const (
{{- range $service := .Services }}
{{- range $method := $service.Methods }}
	{{ $service.Name }}_{{ $method }} = "/{{ $.Package }}.{{ $service.Name }}/{{ $method }}"
{{- end}}

// ini cobaan maulana
// input
{{- range $input := $service.Input }}
	{{ $service.Name }}_{{ "input" }} = "{{ $input }}"
{{- end}}

// Output
{{- range $output := $service.Output }}
	{{ $service.Name }}_{{"output"}} = "{{ $output }}"
{{- end}}

// Unin
{{- range $unin := $service.Unin }}
	{{ $service.Name }}_{{ $unin }} = "/{{ $.Package }}.{{ $service.Name }}/{{ $unin }}"
{{- end}}

{{- end}}

{{- range $message := .Messages }}
// {{ "test" }} = "{{ $message }}"
{{ $message }}_{{"exampvar"}} = "{{ $message }}"
{{- end}}

//end of message

{{- range $fieldMessage := .FieldMessage }}
{{ $fieldMessage }}_{{"exampvar"}} = "{{ $fieldMessage }}"
{{- end}}


// akhir cobaan
)
var (
	FullMethods = []string{{ "{" }}
{{- range $service := .Services }}
{{- range $method := $service.Methods }}
		{{ $service.Name }}_{{ $method }},
{{- end}}
{{- end}}
	{{ "}" }}
)
`
)

type service struct {
	Name    string
	Methods []string
	Output  []string
	Input   []string
	Unin    []string
}

type data struct {
	FileName     string
	GoPackage    string
	Package      string
	Services     []service
	Messages     []string
	FieldMessage []string
}

type fullMethodsGenerator struct {
	*generator
}

// NewFullMethodsGenerator for new method generator
func NewFullMethodsGenerator() *fullMethodsGenerator {
	return &fullMethodsGenerator{generator: newGenerator()}
}

// Generate for generating file
func (g *fullMethodsGenerator) Generate() error {
	return g.generator.generate(g.generateFile)
}

func (g *fullMethodsGenerator) generateFile(protoFile *descriptor.FileDescriptorProto) (*plugin.CodeGeneratorResponse_File, error) {
	if protoFile.Name == nil {
		return nil, errors.New("missing filename")
	}
	if protoFile.GetOptions().GetGoPackage() == "" {
		return nil, errors.New("missing go_package")
	}

	dat := data{
		FileName:  *protoFile.Name,
		GoPackage: protoFile.GetOptions().GetGoPackage(),
		Package:   protoFile.GetPackage(),
		Services:  make([]service, len(protoFile.Service)),
	}

	// messages := new([]string)
	var messages []string
	var excase []string

	for _, vvv := range protoFile.MessageType {
		messages = append(messages, vvv.GetName())

		// get field of message
		for _, valueOfField := range vvv.Field {
			excase = append(excase, valueOfField.GetName()+"_type_"+valueOfField.GetType().String()+valueOfField.GetLabel().String()+"_"+"_number_"+strconv.Itoa(int(valueOfField.GetNumber())))

			// checkoptions := valueOfField.GetOptions()
			// valueOfField.GetOptions().GetCtype().String()

			// messages = append(messages, v2.GetType().String())
			// messages = append(messages, v2.GetName())
			// messages = append(messages, v2.GetTypeName())
			// +valueOfField.GetType().String()

		}

		// get Options field of messages
		log.Println(vvv.GetOptions().String())

		// get enum field of message
		for _, v3 := range vvv.GetEnumType() {
			// initialize enum type
			excase = append(excase, v3.GetName()+"_enum_type")
			// messages = append(messages, v3.GetName())
			// messages = append(messages, fmt.Sprintf("%s_%s", vvv.GetName(), v3.GetName()))
			// // cara mendapatkan enum fieldnya
			for _, v4 := range v3.GetValue() {
				// messages = append(messages, fmt.Sprintf("%s_%s", vvv.GetName(), v4.GetName()))
				log.Println(v4.GetName(), v4.GetNumber())
			}

			// for _, v4 := range v3.GetReservedName() {
			// 	messages = append(messages, v4)
			// }
			// for _, v5 := range v3.GetOptions().UninterpretedOption {
			// 	messages = append(messages, v5.GetName()
			// }

			// for _, v5 := range v3.GetOptions().ExtensionRangeArray() {
			// 	messages = append(messages, v5)
			// }
			// }

			// for _, _ = range vvv.GetExtension() {

			// 	messages = append(messages, "hellocus")
		}

	}

	dat.Messages = messages
	dat.FieldMessage = excase

	for _, svc := range protoFile.Service {
		methods := make([]string, len(svc.Method))
		input := make([]string, len(svc.Method))
		output := make([]string, len(svc.Method))
		unin := []string{}

		for i, method := range svc.Method {
			methods[i] = ucFirst(*method.Name)
			// research di mulai disini
			input[i] = method.GetInputType()
			output[i] = method.GetOutputType()

			// opt := method.GetOptions()
			// for _, vUnin := range method.GetOptions().GetIdempotencyLevel().String() {

			// }
			// methOpt := method.GetOptions()

			// ref := proto.GetProperties(reflect.TypeOf(methOpt.XXX_Unmarshal))
			// log.Println(methOpt.GetUninterpretedOption(), methOpt.GetIdempotencyLevel(), methOpt.String())
			// log.Println(methOpt.String())
			// proto.Get

			// deprecated
			// opt, errJson := json.Marshal(methOpt.XXX_InternalExtensions)
			// log.Println(errJson)
			// var resultOpt interface{}
			// errs := json.Unmarshal(opt, &resultOpt)
			// log.Println(errs)
			// log.Println(resultOpt, reflect.TypeOf(methOpt.XXX_InternalExtensions))
			// for kOpt, vOpt := range resultOpt {
			// 	log.Println(kOpt, vOpt)
			// }

			// unin = append(unin, method.GetOptions().GetIdempotencyLevel().String())

			// for _, vUnin := range opt.GetUninterpretedOption() {

			// unin = append(unin, vUnin.GetAggregateValue())

			// vUnin2 := vUnin.GetName()
			// for _, vUnin3 := range vUnin2 {
			// 	unin = append(unin, vUnin3.String())
			// }
			// }

		}

		dat.Services = append(dat.Services, service{
			Name:    googlegen.CamelCase(svc.GetName()),
			Methods: methods,
			Input:   input,
			Output:  output,
			Unin:    unin,
		})
	}

	buf := bytes.NewBuffer(nil)

	err := template.Must(template.New("").Parse(tmpl)).Execute(buf, dat)
	if err != nil {
		return nil, err
	}

	g.P(buf.String())

	formatted, err := format.Source(g.Bytes())
	if err != nil {
		return nil, err
	}

	return &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(protoFileBaseName(*protoFile.Name) + ".custom.pb.go"),
		Content: proto.String(string(formatted)),
	}, nil
}

func protoFileBaseName(name string) string {
	if ext := path.Ext(name); ext == ".proto" {
		name = name[:len(name)-len(ext)]
	}
	return name
}

func ucFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
