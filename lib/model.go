package lib

// Service list of services
type Service struct {
	Name         string
	Methods      []*Method
	MethodsName  string
	MessageAllEs string
	AllMessage   []Message
	Elastic      bool
}

// Method list of method inside service
type Method struct {
	Name                  string
	Input                 string
	Output                string
	Options               []*Option
	HttpMode              string
	URLPath               string
	InputMessage          Message
	OutputMessage         Message
	IO                    Message
	IsAgregator           bool
	AgregatorMessage      Message
	AgregatorFunction     string
	AgregatorGetByPrimary string
	InputWithAgregator    Message
	IsGetAllMessage       bool
	IsPageLimitFound      bool
	IORelated             bool
	MessageRequest        Message
	MessageResponse       Message
	HasIntegration        bool
	PathPostman           []PathPostman
	ExtraComma            bool
	// IntegMessage          []Message
}

// PathPostman for path postman
type PathPostman struct {
	Name  string
	Extra string
}

// Option for optional
type Option struct {
	Code  string
	Name  string
	Value string
}

// Message for messages
type Message struct {
	Index          int
	Name           string
	IsRepository   bool
	IsElastic      bool
	PrimaryKeyName string
	PrimaryKeyType string
	NumField       int
	Fields         []*Field
	Options        []*Option
	Domain         string
	HasIntegration bool
}

// Enum for messaging enum
type Enum struct {
	Name    string
	Options []*Option
}

// Field for field in messages
type Field struct {
	Index          int
	OriginalName   string
	OriginalType   string
	Name           string
	NameGo         string
	TypeData       string
	TypeDataGo     string
	IsRepeated     bool
	IsOptional     bool
	DefaultValue   string
	IgnoreGorm     bool
	IsPrimaryKey   bool
	RequiredOption bool
	RequiredType   string
	Tag            string
	FullText       bool
	ErrorDesc      string
	IsFieldMessage bool
	MessageTo      Message
	MessageToName  string
	Integration    bool
	IntegrationCfg *IntegrationConfig
	ExtraComma     bool
	PostmanType    string
}

// Data for struct list of data
type Data struct {
	FileName           string
	Src                string
	GoPackage          string
	Package            string
	Services           []Service
	Messages           []Message
	NumMessage         int
	MessageAll         string
	Enums              []*Enum
	Elastic            bool
	TimeStamp          bool
	UseEmptyProto      bool
	WhiteList          []WhitelistOpt
	IntegrationMessage []Message
}

// WhitelistOpt for the struct whitelist
type WhitelistOpt struct {
	Name        string
	ServiceName string
}

// IntegrationConfig for integration config
type IntegrationConfig struct {
	ProtoFileLoc        string
	GrpcMethod          string
	GrpcAddress         string
	ProtoDomain         string
	GrpcPort            string
	GrpcRequestName     string
	GrpcResponseName    string
	GrpcRequestMessage  string
	GrpcResponseMessage string
	GrpcRequestMsg      Message
	GrpcResponseMsg     Message
	ResultRequest       Message
	ResultResponse      Message
	Unique              string
}
