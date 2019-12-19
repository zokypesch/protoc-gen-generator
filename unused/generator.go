package unused

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	googlegen "github.com/golang/protobuf/protoc-gen-go/generator"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pkg/errors"
)

type fileGenerator func(protoFile *descriptor.FileDescriptorProto) (*plugin.CodeGeneratorResponse_File, error)

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
	}
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
