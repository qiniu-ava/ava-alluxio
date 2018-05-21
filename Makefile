SHELL=/bin/bash

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
