// +build !debug

package main

import (
	"qiniu.com/app/avio-walker/compile"
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
