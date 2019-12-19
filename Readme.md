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