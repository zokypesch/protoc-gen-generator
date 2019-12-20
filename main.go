package main

import (
	"fmt"
	"log"

	lib "github.com/zokypesch/protoc-gen-generator/lib"
	template "github.com/zokypesch/protoc-gen-generator/template"
)

func main() {
	// err := lib.NewFullMethodsGenerator().Generate()
	list := []lib.List{
		template.ListFullMethod,
		template.ListTypeScript,
		template.ListTypeScriptScreen,
		template.ListModelGolang,
		template.ListRepositoryGolang,
		template.ListMasterRepoGolang,
		template.ListHandler,
		template.ListConfig,
		// template.ListYaml,
		// template.ListEnv,
		template.ListService,
		template.ListCoreValidation,
		// template.ListMain,
		// template.ListToml,
		template.ListMainv2,
		// template.ListTomlv2,
	}
	res, err := lib.NewMaster(list).Generate()

	if err != nil {
		log.Println(err)
	}

	for _, v := range res {
		pr := fmt.Sprintf("Execute %s has been successfull", v.Filename)
		log.Println(pr)
	}

	// lib.ExecFile(res)
}
