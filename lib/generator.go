package lib

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	googlegen "github.com/golang/protobuf/protoc-gen-go/generator"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pkg/errors"
)

type fileGenerator func(protoFile *descriptor.FileDescriptorProto) (*plugin.CodeGeneratorResponse_File, error)

type fileGeneratorMulti func(protoFile *descriptor.FileDescriptorProto, listTemp List) (*plugin.CodeGeneratorResponse_File, string, Data, error)

type generator struct {
	*googlegen.Generator
	reader io.Reader
	writer io.Writer
}

func newGenerator() *generator {
	return &generator{
		Generator: googlegen.New(),
		reader:    os.Stdin,
		writer:    os.Stdout,
		// writer: ,
	}
}

func (g *generator) generateMulti(generateFile fileGeneratorMulti, listParam []List) ([]fileAfterExecute, error) {
	var after []fileAfterExecute
	err := readRequest(g.reader, g.Request)
	if err != nil {
		return after, err
	}

	g.CommandLineParameters(g.Request.GetParameter())
	g.WrapTypes()
	g.SetPackageNames()
	g.BuildTypeNameMap()
	g.GenerateAllFiles()

	for _, protoFile := range g.Request.ProtoFile {
		if len(protoFile.GetService()) < 1 {
			continue
		}
		for _, genFile := range listParam {
			g.Reset()
			response := &plugin.CodeGeneratorResponse{}
			file, pkgName, datas, err := generateFile(protoFile, genFile)
			if err != nil {
				return after, err
			}
			response.File = append(response.File, file)
			// errWrite := writeResponse(g.writer, response)
			res, errWrite := writeResponseWithList(g.writer, response, genFile, pkgName, datas)
			if errWrite != nil {
				return after, errWrite
			}
			after = append(after, fileAfterExecute{
				Filename: res,
				PkgName:  pkgName,
				Location: genFile.Location,
			})
		}
	}
	return after, nil

}

func (g *generator) generate(generateFile fileGenerator) error {
	err := readRequest(g.reader, g.Request)
	if err != nil {
		return err
	}

	g.CommandLineParameters(g.Request.GetParameter())
	g.WrapTypes()
	g.SetPackageNames()
	g.BuildTypeNameMap()
	g.GenerateAllFiles()
	g.Reset()

	response := &plugin.CodeGeneratorResponse{}
	for _, protoFile := range g.Request.ProtoFile {
		if len(protoFile.GetService()) < 1 {
			continue
		}
		file, err := generateFile(protoFile)
		if err != nil {
			return err
		}
		response.File = append(response.File, file)
	}

	return writeResponse(g.writer, response)
}

func readRequest(r io.Reader, request *plugin.CodeGeneratorRequest) error {
	input, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Wrap(err, "error while reading input")
	}

	if err = proto.Unmarshal(input, request); err != nil {
		return errors.Wrap(err, "error while parsing input proto")
	}

	if len(request.FileToGenerate) == 0 {
		return errors.New("no files to generate")
	}

	return nil
}

func writeResponse(w io.Writer, response *plugin.CodeGeneratorResponse) error {
	output, err := proto.Marshal(response)
	if err != nil {
		return errors.Wrap(err, "failed to marshal output proto")
	}

	_, err = w.Write(output)
	if err != nil {
		return errors.Wrap(err, "failed to write output proto")
	}

	return nil
}

func writeResponseWithList(w io.Writer, response *plugin.CodeGeneratorResponse, list List, pkgName string, datas Data) (string, error) {
	_, err := proto.Marshal(response)
	// _, err := proto.Marshal(response)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal output proto")
	}

	fileName := response.GetFile()[0].GetName()
	location := list.Location

	if list.Location == "" {
		location = "./grpc/pb/" + pkgName + "/"
		if _, err := os.Stat("./grpc/pb"); os.IsNotExist(err) {
			os.Mkdir("./grpc/pb", 0755)
		}

		if _, err := os.Stat("./grpc/pb/" + pkgName); os.IsNotExist(err) {
			os.Mkdir("./grpc/pb/"+pkgName, 0755)
		}

	}

	if list.Lang == "go" {
		fileName = strings.ToLower(fileName)
	} else if list.Lang == "yaml" {
		fileName = "service.yaml"
	} else if list.Lang == "toml" {
		fileName = "Gopkg.toml"

	} else if list.Lang == "docker" {
		fileName = "Dockerfile"

	} else if list.Lang == "readme" {
		fileName = "Readme.md"

	} else if list.Lang == "elastic" {
		if !datas.Elastic {
			return "", nil
		}
		fileName = "elastic.go"
	}

	content := response.GetFile()[0].GetContent()

	if list.ReplaceQuote {
		content = strings.Replace(content, "'", "`", -1)
	}

	foundIndex := strings.Index(location, "%")
	if foundIndex > -1 {
		location = fmt.Sprintf(location, strings.ToLower(datas.Services[0].Name))
	}

	if _, err := os.Stat(location); os.IsNotExist(err) {
		os.Mkdir(location, 0755)
	}

	fileDestPrefix, _ := filepath.Abs(location + fileName)

	d1 := []byte(content)
	err = ioutil.WriteFile(fileDestPrefix, d1, 0644)

	if err != nil {
		return "", errors.Wrap(err, "failed to write output proto")
	}

	return fileName, nil
}

type fileAfterExecute struct {
	Filename string
	PkgName  string
	Location string
}

// ExecFile for execute file
func ExecFile(files []fileAfterExecute) {
	time.Sleep(1 * time.Second)
	for _, file := range files {
		if file.Location == "" {
			continue
		}
		fileCurrentPrefix, _ := filepath.Abs("./grpc/" + file.PkgName + "/" + file.Filename)
		fileDestPrefix, _ := filepath.Abs(file.Location + file.Filename)

		_, errCp := copy(fileCurrentPrefix, fileDestPrefix)
		log.Println("call here", errCp)
		os.Remove(fileCurrentPrefix)
	}

}
