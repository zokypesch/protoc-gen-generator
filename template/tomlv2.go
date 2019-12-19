package template

import lib "github.com/zokypesch/protoc-gen-generator/lib"

var tmplTomlv2 = `

[[constraint]]
  name = "github.com/kelseyhightower/envconfig"
  version = "1.3.0"

[[constraint]]
  name = "github.com/jinzhu/gorm"
  version = "1.9.1"

[prune]
  go-tests = true
  unused-packages = true
`

var ListTomlv2 = lib.List{
	FileType: ".toml",
	Template: tmplTomlv2,
	Location: "./",
	Lang:     "toml",
}
