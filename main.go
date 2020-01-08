package main

import (
	"fmt"
	"log"

	lib "github.com/zokypesch/protoc-gen-generator/lib"
	template "github.com/zokypesch/protoc-gen-generator/template"
)

func main() {
	// err := lib.NewFullMethodsGenerator().Generate()
	// list := []lib.List{
	// 	template.ListFullMethod,
	// 	template.ListTypeScript,
	// 	template.ListTypeScriptScreen,
	// 	template.ListModelGolang,
	// 	template.ListRepositoryGolang,
	// 	template.ListMasterRepoGolang,
	// 	template.ListHandler,
	// 	template.ListConfig,
	// 	template.ListYaml, //disable for reason
	// 	template.ListEnv,  //disable for reason
	// 	template.ListService,
	// 	// template.ListCoreValidation,
	// 	template.ListMain, //disable for reason
	// 	template.ListToml, //disable for reason
	// 	// template.ListMainv2, //non-disable for reason
	// 	// template.ListTomlv2, //disable for reason
	// }

	list := []lib.List{
		template.ListFullMethod,
		template.ListTypeScript,
		template.ListTypeScriptScreen,
		template.ListModelGolang,
		template.ListRepositoryGolang,
		template.ListMasterRepoGolang,
		template.ListHandler,
		template.ListConfig,
		template.ListService,
		template.ListMainv2,  //non-disable for reason
		template.ListTomlv2,  //disable for reason
		template.ListElastic, //disable for reason
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
