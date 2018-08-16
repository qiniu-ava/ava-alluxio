// +build !debug

package main

import (
	"qiniu.com/app/avio-apiserver/compile"
)

func main() {
	bootCompile()
	bootMain()
}

func bootCompile() {
	compile.Mode = "release"
	compile.MongoAuthConfigPath = "/root/config/mongoauth.json"
}

func bootMain() {
	Boot()
}
