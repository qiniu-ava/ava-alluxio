SHELL:=/bin/bash
GOROOT:=$(shell echo ${GOROOT})
GOPATH:=$(shell pwd)/tools/golang
DIR:=$(shell pwd)

clean:
	@echo TODO && exit 1

update:
	git submodules update

hint:
	@echo TODO && exit 1

build:
	./hack/build-alluxio.sh

reborn:
	./docker/alluxio.master.sh restart
	@sleep 3
	./docker/alluxio.worker.sh restart
	@sleep 3
	./docker/alluxio.client.sh restart

tools:
	@echo TODO && exit 1

tools-clean: tools-golang-clean

tools-golang-clean:
	rm -rf tools/golang/bin/*

tools-golang: tools-golang-clean tools-golang-macos tools-golang-linux

tools-golang-macos:
	cd tools/golang/src/qiniu.com/app && env GOOS=darwin go install ./...

tools-golang-linux:
	cd tools/golang/src/qiniu.com/app && env GOOS=linux go install ./...

tools-golang-deploy: tools-golang-avio-check-version tools-golang
	./tools/scripts/deploy.sh

tools-golang-avio-check-version:
	@./tools/scripts/deploy.sh --check-version
