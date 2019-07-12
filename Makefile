#!/usr/bin/make
# GO111MODULE=on

.PHONY: image test

binary: *.go go.* internal
	export GO111MODULE=on
	go build -o ./_output/aether

image:
	docker build -t compute-aether .

test:
	docker build -t aether-tests -f ./test/container/Dockerfile ./test/container
