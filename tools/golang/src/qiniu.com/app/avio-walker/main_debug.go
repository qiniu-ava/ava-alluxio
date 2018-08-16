// +build debug

package main

import (
	"qiniu.com/app/avio-walker/compile"
)

func main() {
	bootCompile()
	bootMain()
}

func bootCompile() {
	compile.Mode = "debug"
}

func bootMain() {
	Boot()
}
