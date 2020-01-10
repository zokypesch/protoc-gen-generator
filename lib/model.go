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
	Name               string
	Input              string
	Output             string
	Options            []*Option
	HttpMode           string
	URLPath            string
	InputMessage       Message
	OutputMessage      Message
	IO                 Message
	IsAgregator        bool
	AgregatorMessage   Message
	AgregatorFunction  string
	InputWithAgregator Message
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
	Fields         []Field
	Options        []*Option
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
}

// Data for struct list of data
type Data struct {
	FileName   string
	Src        string
	GoPackage  string
	Package    string
	Services   []Service
	Messages   []Message
	NumMessage int
	MessageAll string
	Enums      []*Enum
	Elastic    bool
}
