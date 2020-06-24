package lib

import (
	// "go/format"
	"fmt"
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

var useTimeStampIO bool

// func (g *Operations) generateFile(protoFile *descriptor.FileDescriptorProto, listTemp List) (*plugin.CodeGeneratorResponse_File, error) {
func (g *Operations) generateFile(protoFile *descriptor.FileDescriptorProto, listTemp List) (*plugin.CodeGeneratorResponse_File, string, Data, error) {
	if protoFile.Name == nil {
		return nil, "", Data{}, errors.New("missing filename")
	}
	if protoFile.GetOptions().GetGoPackage() == "" {
		return nil, "", Data{}, errors.New("missing go_package")
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
	messageAllEs := ""
	messageAllEsSvc := ""
	indexMsgEs := 0
	useElastic := false
	useTimeStamp := false
	useTimeStampIO = false

	// get message in proto
	for iMType, messageType := range protoFile.MessageType {
		messageName := messageType.GetName()
		primaryKeyName := ""
		primaryKeyType := ""

		var newField []*Field

		if iMType == 0 {
			methodsNameAll += ucDown(messageName)
		} else {
			methodsNameAll += "," + ucDown(messageName)
		}

		integMsg := false
		// get field in message
		totalFieldMessage := 1
		for kMessageField, messageField := range messageType.Field {
			typeData := messageField.GetType().String()
			typeDataGo := messageField.GetType().String()
			originalName := messageField.GetName()
			originalType := typeData
			postmanType := "string"
			isFieldMessage := false

			isOptional := false
			if typeData == "TYPE_MESSAGE" || messageField.GetTypeName() != "" {
				onTypeComb := strings.Split(messageField.GetTypeName(), ".")
				typeData = ucDown(onTypeComb[len(onTypeComb)-1:][0])
				typeDataGo = ucFirst(onTypeComb[len(onTypeComb)-1:][0])
				originalType = ucFirst(onTypeComb[len(onTypeComb)-1:][0])
				if typeDataGo == "Timestamp" {
					typeDataGo = "time.Time"
					useTimeStamp = true
				}

				foundIndexChar := strings.Index(typeDataGo, "_")
				if foundIndexChar == -1 {
					isFieldMessage = true
				}

				originalType = typeDataGo
			} else {
				typeData = grpcTypeToTs(typeData)
				typeDataGo = grpcTypeToGo(typeDataGo)
				originalType = typeDataGo
				postmanType = typeDataPostman(typeDataGo)
			}
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
			tagField := ""
			fullText := false
			errorDescription := "this field is required"
			integration := false
			var intConfig *IntegrationConfig
			foreignKeyDb := ""
			associateKeyDb := ""

			// get options in field
			for _, vOptField := range newFieldOptions {
				switch getStringFromOptCode(vOptField.Code) {
				case "ignoreFieldDb":
					if res, err := strconv.Atoi(vOptField.Value); err == nil && res == 1 {
						ignoreGorm = true
						tagField += `gorm:"-"`
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
					tagField += fmt.Sprintf(` validate:"%s"`, requiredType)
				case "fulltext":
					if res, err := strconv.Atoi(vOptField.Value); err == nil && res == 1 {
						fullText = true
					}
				case "errorDesc":
					errorDescription = vOptField.Value
				case "foreignKey":
					foreignKeyDb = vOptField.Value
				case "associateKey":
					associateKeyDb = vOptField.Value
				// for integratiom
				case "integration":
					if res, err := strconv.Atoi(vOptField.Value); err == nil && res == 1 {
						integration = true
						integMsg = true
					}
				}
			}

			if foreignKeyDb != "" || associateKeyDb != "" {
				tagField += fmt.Sprintf(` gorm:"foreignkey:%s;association_foreignkey:%s"`, foreignKeyDb, associateKeyDb)
			}
			// for get integration
			intConfig = integrationCheckOpt(newFieldOptions)
			intConfig.Unique = fmt.Sprintf("%s%s", ucFirst(intConfig.ProtoDomain), ucFirst(intConfig.GrpcMethod))

			newField = append(newField, &Field{
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
				Tag:            tagField,
				FullText:       fullText,
				Index:          kMessageField + 1,
				ErrorDesc:      errorDescription,
				IsFieldMessage: isFieldMessage,
				Integration:    integration,
				IntegrationCfg: intConfig,
				ExtraComma:     totalFieldMessage < len(messageType.Field),
				PostmanType:    postmanType,
			})
			totalFieldMessage++
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
		elastic := false
		for _, vOpt := range newMessageOptions {
			switch getStringFromOptCode(vOpt.Code) {
			case "isRepository":
				if res, err := strconv.Atoi(vOpt.Value); err == nil && res == 1 {
					isRepo = true
				}
			case "elastic":
				if res, err := strconv.Atoi(vOpt.Value); err == nil && res == 1 {
					elastic = true
					useElastic = true
					if indexMsgEs == 0 {
						messageAllEs += "es" + ucFirst(messageName)
						messageAllEsSvc += "es" + ucFirst(messageName) + " core.ESModule"
					} else {
						messageAllEs += "," + "es" + ucFirst(messageName)
						messageAllEsSvc += "," + "es" + ucFirst(messageName) + " core.ESModule"
					}
					indexMsgEs++
				}
			}

		}

		newMessage = append(newMessage, Message{
			Index:          iMType + 1,
			Name:           ucDown(messageName),
			Fields:         newField,
			Options:        newMessageOptions,
			IsRepository:   isRepo,
			PrimaryKeyName: primaryKeyName,
			PrimaryKeyType: primaryKeyType,
			IsElastic:      elastic,
			NumField:       len(newField),
			HasIntegration: integMsg,
		})

	}

	// rebuild message
	for _, vNewMessage := range newMessage {
		for _, vNewFields := range vNewMessage.Fields {
			if vNewFields.IsFieldMessage == true {
				resNewField, foundNewField := findMessage(vNewFields.OriginalType, newMessage)
				if foundNewField {
					vNewFields.MessageTo = resNewField
					vNewFields.MessageToName = ucFirst(resNewField.Name)
				}
			}
		}
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
	datas.NumMessage = len(newMessage)
	datas.MessageAll = messageAllEs
	datas.TimeStamp = useTimeStamp

	var integrationMessage []Message

	useEmptyProto := false
	var whitelist []WhitelistOpt
	// get service in protofile
	for kSvc, svc := range protoFile.Service {
		methods := make([]*Method, len(svc.Method))

		totalMethodSvc := 1
		// get method inside service
		for i, method := range svc.Method {
			// var methodIntegration []Message
			methodsNameAll += ", Output" + ucFirst(*method.Name)

			methOpt := method.GetOptions().String()
			newOptions := stringToOpt(methOpt)
			httpMode := "get"
			urlPath := ""
			var pathPostman []PathPostman
			isAgregator := false
			agregatorMessage := Message{}
			agregatorMessageName := ""
			agregatorGetByPrimary := ""
			agregatorFunction := ""
			isGetAllMessage := false
			isPageParamFound := false
			isLimitParamFound := false

			// get options in service
			for _, vOpt := range newOptions {
				switch getStringFromOptCode(vOpt.Code) {
				case "urlPath":
					// newURL, newMode := getHttpModeWithUrl(vOpt.Value)
					// httpMode = newMode
					urlPath = getHttpUrl(vOpt.Value)
					urlPost := strings.Replace(urlPath, "{id}", "1", -1)
					urlPostSlice := strings.Split(urlPost, "/")

					iSlice := 1
					for _, vSlice := range urlPostSlice {
						extraString := ""
						if iSlice < len(urlPostSlice) {
							extraString = ","
						}
						if len(vSlice) > 0 {
							pathPostman = append(pathPostman, PathPostman{Name: vSlice, Extra: extraString})
						}
						iSlice++

					}
				case "httpMode":
					httpMode = vOpt.Value
				case "agregator":
					isAgregator = true
					msgFun := strings.Split(vOpt.Value, ".")
					agregatorMessage = getMessageByName(msgFun[0], datas.Messages)
					agregatorMessageName = ucFirst(agregatorMessage.Name)
					agregatorGetByPrimary = fmt.Sprintf("GetBy%s", ucFirst(agregatorMessage.PrimaryKeyName))
					if len(msgFun) == 2 {
						agregatorFunction = msgFun[1]
					}
				case "whitelist":
					if res, err := strconv.Atoi(vOpt.Value); err == nil && res == 1 {
						whitelist = append(whitelist, WhitelistOpt{Name: ucFirst(*method.Name), ServiceName: googlegen.CamelCase(svc.GetName())})
					}
				}
			}
			onInputMethod := strings.Split(method.GetInputType(), ".")
			typeInputMethod := ucDown(onInputMethod[len(onInputMethod)-1:][0])
			onOutputMethod := strings.Split(method.GetOutputType(), ".")
			typeOutputMethod := ucDown(onOutputMethod[len(onOutputMethod)-1:][0])

			inputMessage := getMessageByName(typeInputMethod, datas.Messages)
			outputMessage := getMessageByName(typeOutputMethod, datas.Messages)
			inputOutputMessage := getIO(outputMessage, agregatorMessage)
			inputWithAgregator := getIO(inputMessage, agregatorMessage)

			if outputMessage.HasIntegration {
				// check duplication
				foundDupl := false
				for _, vIntMsg := range integrationMessage {
					if vIntMsg.Name == outputMessage.Name {
						foundDupl = true
					}
				}

				if !foundDupl {
					// mapping grpc request to response field
					// mapping grpc response to field response
					for _, dtField := range outputMessage.Fields {
						if dtField.Integration {
							dtField.IntegrationCfg.GrpcRequestMsg = getMessageByName(dtField.IntegrationCfg.GrpcRequestMessage, datas.Messages)
							dtField.IntegrationCfg.GrpcResponseMsg = getMessageByName(dtField.IntegrationCfg.GrpcResponseMessage, datas.Messages)
							dtField.IntegrationCfg.ResultRequest = getIO(outputMessage, dtField.IntegrationCfg.GrpcRequestMsg)
							dtField.IntegrationCfg.ResultResponse = getIO(dtField.MessageTo, dtField.IntegrationCfg.GrpcResponseMsg)
						}
					}
					// methodIntegration = append(methodIntegration, outputMessage)
					integrationMessage = append(integrationMessage, outputMessage)
				}

			}

			if typeInputMethod == "empty" {
				useEmptyProto = true
			}

			for _, vAllMsg := range outputMessage.Fields {
				if strings.ToLower(vAllMsg.Name) == "items" && isAgregator {
					newDataTypeWithAgg := strings.Replace(vAllMsg.TypeDataGo, "[]", "", -1)
					if newDataTypeWithAgg == agregatorMessageName {
						// log.Println(newDataTypeWithAgg, agregatorMessageName)
						isGetAllMessage = true
						break
					}
				}
			}

			for _, vAllMsg := range inputMessage.Fields {
				if strings.ToLower(vAllMsg.Name) == "page" && isAgregator {
					isPageParamFound = true
				} else if strings.ToLower(vAllMsg.Name) == "perpage" && isAgregator {
					isLimitParamFound = true
				}
			}

			methods[i] = &Method{
				Name:                  ucFirst(*method.Name),
				Input:                 typeInputMethod,
				Output:                typeOutputMethod,
				Options:               newOptions,
				HttpMode:              httpMode,
				URLPath:               urlPath,
				InputMessage:          inputMessage,
				OutputMessage:         outputMessage,
				IO:                    inputOutputMessage,
				IsAgregator:           isAgregator,
				AgregatorMessage:      agregatorMessage,
				AgregatorFunction:     agregatorFunction,
				InputWithAgregator:    inputWithAgregator,
				IsGetAllMessage:       isGetAllMessage,
				AgregatorGetByPrimary: agregatorGetByPrimary,
				IsPageLimitFound:      isPageParamFound && isLimitParamFound,
				IORelated:             len(inputOutputMessage.Fields) > 0,
				HasIntegration:        outputMessage.HasIntegration,
				PathPostman:           pathPostman,
				ExtraComma:            totalMethodSvc < len(svc.Method),
				// IntegMessage:          methodIntegration,
			}
			totalMethodSvc++
			isGetAllMessage = false
		}
		// put service in datas
		datas.Services[kSvc] = Service{
			Name:         googlegen.CamelCase(svc.GetName()),
			Methods:      methods,
			MethodsName:  methodsNameAll,
			MessageAllEs: messageAllEsSvc,
			AllMessage:   newMessage,
			Elastic:      useElastic,
		}
	}

	datas.WhiteList = whitelist
	datas.Elastic = useElastic
	datas.UseEmptyProto = useEmptyProto
	datas.IntegrationMessage = integrationMessage
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
			"allowRequest":    allowRequest,
			"currentLoc":      currentLoc,
			"generateGuid":    generateGUID,
			"toupper":         strings.ToUpper,
			"strReplaceParam": strReplaceParam,
		}).
		Parse(curList.Template)).Execute(buf, datas)

	if err != nil {
		log.Println("here error", err)
		return nil, "", datas, err
	}

	// generate protobuffer
	g.Generator.P(buf.String())

	formatted := g.Generator.Bytes()

	// return code generator response
	return &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(ucFirst(protoFileBaseName(*protoFile.Name)) + curList.FileType), // ".custom.pb.go"
		Content: proto.String(string(formatted)),
	}, protoFile.GetOptions().GetGoPackage(), datas, nil
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

func getIO(input Message, output Message) Message {
	var fields []*Field
	index := 1
	for _, v := range input.Fields {
		for _, v2 := range output.Fields {
			if v.Name == v2.Name && v.TypeDataGo == v2.TypeDataGo {
				newField := Field(*v)
				newField.Index = index

				if newField.TypeDataGo == "Timestamp" {
					newField.TypeDataGo = "time.Time"
					useTimeStampIO = true
				}

				fields = append(fields, &newField)
				index++
			} else {
				// DO SOMETHING
			}
		}
	}

	for k, v := range fields {
		if (k + 1) < len(fields) {
			v.ExtraComma = true
		}
	}

	return Message{
		Name:   output.Name,
		Fields: fields,
	}
}

func findMessage(name string, param []Message) (Message, bool) {
	for _, v := range param {
		if strings.ToLower(v.Name) == strings.ToLower(name) {
			newMessage := Message(v)
			return newMessage, true
		}
	}
	return Message{}, false
}

var ignoreField = []string{
	"createdat",
	"updatedat",
	"createdby",
	"updatedby",
	"updateddate",
}

func allowRequest(field string) bool {
	for _, v := range ignoreField {
		if v == strings.ToLower(field) {
			return false
		}
	}
	return true
}

func strReplaceParam(param string) string {
	return strings.Replace(param, "{id}", "1", -1)
}
