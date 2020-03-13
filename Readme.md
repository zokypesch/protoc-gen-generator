#generate


# installnya go install .

# --grpc-gateway_out=logtostderr=true:example
ini aritnya dia generate folder example

protoc -I grpc/example example.proto --go_out=plugins=grpc:example --gofullmethods_out=example --grpc-gateway_out=logtostderr=true:grpc/example -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis 


protoc -I grpc/example example.proto --go_out=plugins=grpc:grpc/example --generator_out=grpc/example --grpc-gateway_out=logtostderr=true:grpc/example -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis 

# templatenya
ada convention name jadi ...
nama foldernya harus `protoc-gen-` {nama-folder-kt}
supaya kebaca sama protoc-gen library ].
nah syarat lainnya harus di install

# klao mau uninstall
go clean -i github.com/zokypesch/protoc-gen-generator

# liat servicenya
ls $GOPATH/bin/protoc-gen-generator

protoc -I grpc/example example.proto \
--generator_out=grpc/example \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
--go_out=plugins=grpc:grpc/example \
--grpc-gateway_out=logtostderr=true:grpc/example

- type data bool (done)
- tambahin field typedata go: disesuaikan (done)
- bikin message option tablename: isinya nama table
- bkin model
- bikin repository
- Bikin handler to isinya function sm return default
- bikin main.go

- master repository (done)
- handler + validator (done)
- grpc validator (done)
- aggregator + udh nembak ke repo
- generate config (done)
- generate env (done)
- generate yaml (done)
- generate main

{{- if $method.IsAgregator}}

{{- if eq $method.AgregatorFunction "Create"}}
	res, err := svc.repo.{{ ucfirst AgregatorMessage.Name }}.Create()
{{- end}}

{{- range $agr := $method.AgregatorMessage.Fields }}
{{- end}}

{{- end}}

env GOOS=linux GOARCH=arm go build -v ./
env GOOS=windows GOARCH=arm go build -v ./
sangkuriangV2 grpc/proto/simple simple grpc/pb/simple

copy file to go/bin

# create in your bash profile
sangkuriang() {
  protoc -I $1 $2.proto --generator_out=$3 -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:$3 --grpc-gateway_out=logtostderr=true:$3
  dep ensure -v
  echo "files has been generate"
}

# how to generate
- add folder proto in folder grpc/proto/simple
- files name simple.proto
```
syntax = "proto3";

option go_package = "simple";
package simple;

import "google/protobuf/empty.proto";
import "google/api/annotations.proto";
import "google/protobuf/descriptor.proto";

extend google.protobuf.MethodOptions {
    string httpMode = 50056;
    string agregator = 50062;
    string testCAse = 50063;
}

extend google.protobuf.MessageOptions {
    bool isRepo = 50057;
}

extend google.protobuf.FieldOptions {
    bool ignoreFieldDb = 50058;
    bool isPrimaryKey = 50059;
    bool required = 50060;
    string required_type = 50061;
}

service simple {
    rpc HelloWorld(google.protobuf.Empty) returns(HelloWorldMessage) {
        option (google.api.http) = {
            get: "/v1/example/echo"
        };
        option(httpMode) = "get";
        option(testCAse) = "test-api-payment";
    };

    rpc HelloWorldGetWithParam(Test) returns(HelloWorldMessage) {
        option (google.api.http) = {
            get: "/v1/example/echo"
        };
    };

    rpc HelloTest(Test) returns(HelloWorldMessage) {
        option (google.api.http) = {
            post: "/v1/example/echo",
            body: '*'
        };
        option(httpMode) = "post";
    };
    rpc HelloBring(HelloWorldMessage) returns(HelloWorldMessage) {
        option (google.api.http) = {
            delete: "/v1/example/echo"
        };
        option(httpMode) = "delete";
        option(agregator) = "HelloWorldMessage.Create";
    };
}

message HelloWorldMessage {
    option (isRepo) = true;
    string name = 1 [(isPrimaryKey) = true, (required) = true, (required_type) ="min_max*5.10"];
    string message = 2 [(required) = true, (required_type) ="not_empty_string"];
    repeated Test test = 3 [(ignoreFieldDb) = true];
    Single single = 4;
    string email = 5 [(required) = true, (required_type) ="email"];
}

message Test {
    string data_test = 1 [(isPrimaryKey) = true];
    int64 numbers_data = 2;
    bool ex_bool = 3;
    option (isRepo) = true;
}

enum TType {
    CREATE = 0;
    UPDATE = 1;
    DELETE = 2;
}

message Single {
    string who = 1;
    TType ttypes = 2; 
}
```

- run this command
sangkuriang grpc/proto/simple simple grpc/pb/simple
- boooommmmmm


```
{{- range $field := $method.InputMessage.Fields }}
{{- if $field.RequiredOption}}
	if err := core.Validate("{{ $field.RequiredType }}", in.{{ ucfirst $field.Name }}, "{{ ucfirst $field.Name }}"); err != nil {
		return &pb.{{ ucfirst $method.Output }}{}, err
	}
{{- end}}
{{- end}}

repo backup
func (repo *{{ ucfirst $msg.Name }}Repository) Update({{ ucdown $msg.PrimaryKeyName }} {{ $msg.PrimaryKeyType }}, payload *model.{{ ucfirst $msg.Name }}) (*model.{{ ucfirst $msg.Name }}, error) {
	err := repo.db.Model(&model.{{ ucfirst $msg.Name }}{{unescape "{"}}{{ $msg.PrimaryKeyName }}: {{ ucdown $msg.PrimaryKeyName }}{{unescape "}"}}).Update(payload).Error

	return payload, err
}

func (repo *{{ ucfirst $msg.Name }}Repository) Delete({{ ucdown $msg.PrimaryKeyName }} {{ $msg.PrimaryKeyType }}) error {
	err := repo.db.Delete(model.{{ ucfirst $msg.Name }}{{unescape "{"}}{{unescape "}"}}, "{{ underscore $msg.PrimaryKeyName }} = ?", {{ ucdown $msg.PrimaryKeyName }}).Error

	return err
}

Delete({{ $msg.PrimaryKeyName }} {{ $msg.PrimaryKeyType }}) error

GetBy{{ $msg.PrimaryKeyName }}({{ $msg.PrimaryKeyName }} {{ $msg.PrimaryKeyType }}) *model.{{ ucfirst $msg.Name }}
```

# generate option
protoc --go_out=. proto/options.proto

# notes for using decorator
// Query count
queryCount := r.db.Select("count(id) count").Model(Price{})
decorCount := decorator.NewServiceDecorator(queryCount, q)
queryCount, _ = decorCount.AppendWhere()
if err = queryCount.Count(&pagingResp.Total).Error; err != nil {
    return nil, nil, err
}