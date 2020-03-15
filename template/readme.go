package template

import lib "github.com/zokypesch/protoc-gen-generator/lib"

var tmplReadme = `# {{ ucfirst (getFirstService .Services).Name }}

# How to build to docker ?
docker build -t {your_username}/{repo_name} .

# login to docker ?
username: {your_username}
password: {your_password}

# push to cloud ?
docker push {your_username}/{repo_name}

# how to auto generate proto using db ??
proto-tools -cmd=gen-proto-db host=localhost name={{ ucdown (getFirstService .Services).Name }} user=root password=

# How to regenerate ??
please split your own code within auto generate code, because auto generated code will be replace your own code or you can mark it by using type in code // DO_NOT_REPLACE
sangkuriang grpc/proto/{{ ucdown (getFirstService .Services).Name }} {{ ucdown (getFirstService .Services).Name }} grpc/pb/{{ ucdown (getFirstService .Services).Name }} 
`

var ListReadme = lib.List{
	FileType:     "Readme",
	Template:     tmplReadme,
	Location:     "./",
	Lang:         "readme",
	ReplaceQuote: false,
}
