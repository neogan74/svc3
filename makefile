SHELL := /bin/bash

run:
	go run main.go

build:
	# go build -ldflags "-X main.build=${VCS_REF}"
	go build -ldflags "-X main.build=${BUILD_REF}"

# Building containers

VERSION := 0.1

all: svc

svc:
	docker  build \
		-f leo/docker/Dockerfile \
		-t svc-arm64:${VERSION} \
		--build-arg VCS_REF=${VERSION} \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%d%H:%M:%SZ"` \
		.