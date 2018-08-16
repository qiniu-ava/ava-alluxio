SHELL:=/bin/bash
GOROOT:=$(shell echo ${GOROOT})
GOPATH:=$(shell pwd)/tools/golang
DIR:=$(shell pwd)
KAFKA_VERSION:=0.10.2.1

clean:
	@echo TODO && exit 1

update:
	git submodules update

hint:
	@./tools/scripts/linter.sh

build:
	./tools/scripts/build-alluxio.sh

reborn:
	./docker/alluxio.master.sh restart
	@sleep 3
	./docker/alluxio.worker.sh restart
	@sleep 3
	./docker/alluxio.client.sh restart

kafka:
	cd docker/app/kafka && docker build --build-arg KAFKA_VERSION=${KAFKA_VERSION} -t reg-xs.qiniu.io/atlab/avio-kafka .

zookeeper:
	cd docker/app/zookeeper && docker build -t reg-xs.qiniu.io/atlab/avio-zookeeper .

tools:
	@echo TODO && exit 1

tools-clean: tools-golang-clean

tools-golang-clean:
	rm -rf tools/golang/bin/avio tools/golang/bin/linux_amd64

tools-golang: tools-golang-clean tools-golang-macos tools-golang-linux

tools-golang-macos:
	cd tools/golang/src/qiniu.com/app && env GOOS=darwin go install ./...

tools-golang-hack:
	cd tools/golang/src/qiniu.com/app && go install -tags debug ./...

tools-golang-linux:
	cd tools/golang/src/qiniu.com/app && env GOOS=linux go install ./...

tools-golang-deploy: tools-golang-avio-check-version tools-golang
	./tools/scripts/deploy.sh

tools-golang-avio-check-version:
	@./tools/scripts/deploy.sh --check-version
