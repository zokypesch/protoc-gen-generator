package lib

import (
	// "go/format"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"bytes"
	"html/template"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	googlegen "github.com/golang/protobuf/protoc-gen-go/generator"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pkg/errors"
)

// Operations struct of list data
type Operations struct {
	Datas       Data
	Generator   *generator
	List        []List
	CurrentList List
}

// List of list and generate name
type List struct {
	FileType     string
	Template     string
	Location     string
	Lang         string
	ReplaceQuote bool
}

// NewMaster for new master
func NewMaster(list []List) *Operations {
	return &Operations{Generator: newGenerator(), List: list}
}

// Generate for generating file
func (g *Operations) Generate() ([]fileAfterExecute, error) {
	//multi
	res, err := g.Generator.generateMulti(g.generateFile, g.List)
	if err != nil {
		// log.Println("error here", err)
		return nil, err
	}

	// single
	// err := g.Generator.generate(g.generateFile)
	// if err != nil {
	// 	log.Println("error here", err)
	// 	return err
	// }

	log.Println("generate success")
	return res, nil
}

func (g *Operations) setCurrentList(index int) {
	// log.Println(g.List[index])
	g.CurrentList = g.List[index]
}

// func (g *Operations) generateFile(protoFile *descriptor.FileDescriptorProto, listTemp List) (*plugin.CodeGeneratorResponse_File, error) {
func (g *Operations) generateFile(protoFile *descriptor.FileDescriptorProto, listTemp List) (*plugin.CodeGeneratorResponse_File, string, error) {
	if protoFile.Name == nil {
		return nil, "", errors.New("missing filename")
	}
	if protoFile.GetOptions().GetGoPackage() == "" {
		return nil, "", errors.New("missing go_package")
	}

	// initial message
	datas := Data{
		FileName:  *protoFile.Name,
		GoPackage: protoFile.GetOptions().GetGoPackage(),
		Package:   protoFile.GetPackage(),
		Services:  make([]Service, len(protoFile.Service)),
		Messages:  make([]Message, len(protoFile.MessageType)),
	}

	var newMessage []Message
	var enums []*Enum
	methodsNameAll := ""

	// get message in proto
	for iMType, messageType := range protoFile.MessageType {
		messageName := messageType.GetName()
		primaryKeyName := ""
		primaryKeyType := ""

		var newField []Field

		if iMType == 0 {
			methodsNameAll += ucDown(messageName)
		} else {
			methodsNameAll += "," + ucDown(messageName)
		}

		// get field in message
		for _, messageField := range messageType.Field {
			typeData := messageField.GetType().String()
			typeDataGo := messageField.GetType().String()
			originalName := messageField.GetName()
			originalType := typeData

			isOptional := false

			if typeData == "TYPE_MESSAGE" || messageField.GetTypeName() != "" {
				onTypeComb := strings.Split(messageField.GetTypeName(), ".")
				typeData = ucDown(onTypeComb[len(onTypeComb)-1:][0])
				typeDataGo = ucFirst(onTypeComb[len(onTypeComb)-1:][0])
				originalType = typeDataGo
			} else {
				typeData = grpcTypeToTs(typeData)
				typeDataGo = grpcTypeToGo(typeDataGo)
				originalType = typeDataGo
			}
			// log.Println(typeData, messageField.GetTypeName(), messageField.GetOptions())
			isRepeated := false

			if "LABEL_REPEATED" == messageField.GetLabel().String() {
				isRepeated = true
				typeData += "[]"
				typeDataGo = "[]" + typeDataGo
			} else if "LABEL_OPTIONAL" == messageField.GetLabel().String() {
				isOptional = true
			}
			newFieldOptions := stringToOpt(messageField.GetOptions().String())
			ignoreGorm := false
			isPrimaryKey := false
			isRequiredField := false
			requiredType := ""

			// get options in field
			for _, vOptField := range newFieldOptions {
				switch getStringFromOptCode(vOptField.Code) {
				case "ignoreFieldDb":
					if res, err := strconv.Atoi(vOptField.Value); err == nil && res == 1 {
						ignoreGorm = true
					}
				case "isPrimaryKey":
					isPrimaryKey = true
					primaryKeyName = underscoreToGoFormat(messageField.GetName())
					primaryKeyType = typeDataGo
				case "requiredField":
					if res, err := strconv.Atoi(vOptField.Value); err == nil && res == 1 {
						isRequiredField = true
					}
				case "requiredType":
					requiredType = vOptField.Value
				}
			}
			// log.Println(isRequiredField, " jalanin")
			newField = append(newField, Field{
				Name:           messageField.GetJsonName(),
				NameGo:         underscoreToGoFormat(messageField.GetName()),
				TypeData:       typeData,
				TypeDataGo:     typeDataGo,
				IsRepeated:     isRepeated,
				IsOptional:     isOptional,
				IgnoreGorm:     ignoreGorm,
				IsPrimaryKey:   isPrimaryKey,
				OriginalName:   originalName,
				OriginalType:   originalType,
				RequiredOption: isRequiredField,
				RequiredType:   requiredType,
			})
			// log.Println(typeData)
		}

		// get enum declare in message protofile
		for _, typEnum := range messageType.GetEnumType() {
			listOptEnum := make([]*Option, len(typEnum.GetValue()))

			for kValEnum, valEnum := range typEnum.GetValue() {
				listOptEnum[kValEnum] = &Option{
					Name:  valEnum.GetName(),
					Code:  valEnum.GetName(),
					Value: strconv.Itoa(int(valEnum.GetNumber())),
				}
			}
			enums = append(enums, &Enum{
				Name:    ucDown(typEnum.GetName()),
				Options: listOptEnum,
			})
		}

		// get message options
		newMessageOptions := stringToOpt(messageType.GetOptions().String())
		isRepo := false
		for _, vOpt := range newMessageOptions {
			switch getStringFromOptCode(vOpt.Code) {
			case "isRepository":
				if res, err := strconv.Atoi(vOpt.Value); err == nil && res == 1 {
					isRepo = true
				}
			}
		}

		newMessage = append(newMessage, Message{
			Name:           ucDown(messageName),
			Fields:         newField,
			Options:        newMessageOptions,
			IsRepository:   isRepo,
			PrimaryKeyName: primaryKeyName,
			PrimaryKeyType: primaryKeyType,
		})

		// log.Println(messageType.GetName())
	}

	// get enum ini protofile
	for _, typEnum := range protoFile.GetEnumType() {
		listOptEnum := make([]*Option, len(typEnum.GetValue()))

		for kValEnum, valEnum := range typEnum.GetValue() {
			listOptEnum[kValEnum] = &Option{
				Name:  valEnum.GetName(),
				Code:  valEnum.GetName(),
				Value: strconv.Itoa(int(valEnum.GetNumber())),
			}
		}
		enums = append(enums, &Enum{
			Name:    ucDown(typEnum.GetName()),
			Options: listOptEnum,
		})
	}

	datas.Messages = newMessage
	datas.Enums = enums

	// check enums
	// for _, vCheckEnum := range datas.Enums {
	// 	log.Println(vCheckEnum.Name, vCheckEnum.Option, vCheckEnum.Value)
	// }

	// get service in protofile
	for kSvc, svc := range protoFile.Service {
		methods := make([]*Method, len(svc.Method))

		// get method inside service
		for i, method := range svc.Method {
			// log.Println(method, i)
			// if i == 0 {
			// 	methodsNameAll += "Output" + ucFirst(*method.Name)
			// } else {
			// 	methodsNameAll += ", Output" + ucFirst(*method.Name)
			// }
			methodsNameAll += ", Output" + ucFirst(*method.Name)

			methOpt := method.GetOptions().String()
			newOptions := stringToOpt(methOpt)
			httpMode := "get"
			urlPath := ""
			isAgregator := false
			agregatorMessage := Message{}
			agregatorFunction := ""

			// get options in service
			for _, vOpt := range newOptions {
				switch getStringFromOptCode(vOpt.Code) {
				case "urlPath":
					newURL, newMode := getHttpModeWithUrl(vOpt.Value)
					httpMode = newMode
					urlPath = newURL
				case "agregator":
					isAgregator = true
					msgFun := strings.Split(vOpt.Value, ".")
					agregatorMessage = getMessageByName(msgFun[0], datas.Messages)
					if len(msgFun) == 2 {
						agregatorFunction = msgFun[1]
					}
				}
			}
			// log.Println(httpMode, urlPath)
			onInputMethod := strings.Split(method.GetInputType(), ".")
			typeInputMethod := ucDown(onInputMethod[len(onInputMethod)-1:][0])
			onOutputMethod := strings.Split(method.GetOutputType(), ".")
			typeOutputMethod := ucDown(onOutputMethod[len(onOutputMethod)-1:][0])

			inputMessage := getMessageByName(typeInputMethod, datas.Messages)
			outputMessage := getMessageByName(typeOutputMethod, datas.Messages)

			methods[i] = &Method{
				Name:              ucFirst(*method.Name),
				Input:             typeInputMethod,
				Output:            typeOutputMethod,
				Options:           newOptions,
				HttpMode:          httpMode,
				URLPath:           urlPath,
				InputMessage:      inputMessage,
				OutputMessage:     outputMessage,
				IsAgregator:       isAgregator,
				AgregatorMessage:  agregatorMessage,
				AgregatorFunction: agregatorFunction,
			}
		}

		// for _, vM := range methods {
		// 	for _, vOp := range vM.Options {
		// 		log.Println(vOp.Code, vOp.Value)
		// 	}
		// }

		// put service in datas
		datas.Services[kSvc] = Service{
			Name:        googlegen.CamelCase(svc.GetName()),
			Methods:     methods,
			MethodsName: methodsNameAll,
		}
		// datas.Services = append(datas.Services, Service{
		// 	Name:    googlegen.CamelCase(svc.GetName()),
		// 	Methods: methods,
		// })
	}

	// for _, svcTest := range datas.Services {
	// 	log.Println(svcTest.Name, svcTest.Methods)
	// }

	// get current folder path for assign src
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	gopath += "/src/"
	filePath, _ := filepath.Abs("./")
	filePath = strings.Replace(filePath, gopath, "", -1)

	datas.Src = filePath
	// generate process
	curList := listTemp
	// log.Println(curList.FileType, "has been executed")
	buf := bytes.NewBuffer(nil)

	// mapping datas to template
	err := template.Must(template.New("").Funcs(
		template.FuncMap{
			"unescape":        unescape,
			"ucfirst":         ucFirst,
			"ucdown":          ucDown,
			"protoRemove":     protoFileBaseName,
			"getFirstService": getFirstService,
			"upper":           strings.ToUpper,
			"underscore":      underscore,
		}).
		Parse(curList.Template)).Execute(buf, datas)

	if err != nil {
		log.Println("here error", err)
		return nil, "", err
	}

	// generate protobuffer
	g.Generator.P(buf.String())

	// check format golang
	// formatted, err := format.Source(g.Generator.Bytes())
	// if err != nil {
	// 	return nil, err
	// }
	formatted := g.Generator.Bytes()

	// return code generator response
	return &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(ucFirst(protoFileBaseName(*protoFile.Name)) + curList.FileType), // ".custom.pb.go"
		Content: proto.String(string(formatted)),
	}, protoFile.GetOptions().GetGoPackage(), nil
}

func protoFileBaseName(name string) string {
	if ext := path.Ext(name); ext == ".proto" {
		name = name[:len(name)-len(ext)]
	}
	return name
}

func getMessageByName(name string, messages []Message) Message {
	for k, v := range messages {
		if ucFirst(name) == ucFirst(v.Name) {
			return messages[k]
		}
	}
	return Message{}
}
