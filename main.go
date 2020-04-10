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
	// 	template.ListYaml,
	// 	template.ListEnv,
	// 	template.ListService,
	// 	template.ListMain,
	// 	template.ListToml,
	// 	template.ListElastic,
	// }

	log.Println("starting generate")
	// list := []lib.List{
	// 	template.ListFullMethod,
	// 	template.ListTypeScript,
	// 	template.ListTypeScriptScreen,
	// 	template.ListModelGolang,
	// 	template.ListRepositoryGolang,
	// 	template.ListMasterRepoGolang,
	// 	template.ListHandler,
	// 	template.ListConfig,
	// 	template.ListService,
	// 	template.ListMainv2,
	// 	template.ListTomlv2,
	// 	template.ListElastic,
	// 	template.ListTypeScriptValidation,
	// }

	// v3 for prakerja
	list := []lib.List{
		template.ListFullMethod,
		template.ListModelGolang,
		template.ListRepositoryGolang,
		template.ListMasterRepoGolang,
		template.ListHandler,
		template.ListConfig,
		template.ListService,
		template.ListMainv3,
		// template.ListTomlv2,
		template.ListDocker,
		template.ListReadme,
		template.ListIntegration,
		template.ListGitIgone,
		template.ListPostman,
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
