package template

import lib "github.com/zokypesch/protoc-gen-generator/lib"

var tmplToml = `ignored = ["gitlab.com/ruangguru/source/shared-lib*"]

required = [
  "github.com/sirupsen/logrus",
  "github.com/satori/go.uuid",
  "github.com/gin-gonic/gin",
  "github.com/grpc/grpc-go/status",
  "google.golang.org/grpc/reflection",
  "github.com/prometheus/client_golang/prometheus/promhttp"
]

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

var ListToml = lib.List{
	FileType: ".toml",
	Template: tmplToml,
	Location: "./",
	Lang:     "toml",
}
