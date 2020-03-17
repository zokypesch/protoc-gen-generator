package template

import lib "github.com/zokypesch/protoc-gen-generator/lib"

var tmplGitIgnore = `vendor/

`

var ListGitIgone = lib.List{
	FileType: ".gitignore",
	Template: tmplGitIgnore,
	Location: "./",
	Lang:     "git",
}
