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

KIND_CLUSTER := leo-cluster 

kind-up:
	kind create cluster \
		--name ${KIND_CLUSTER} \
		--config leo/kind/kind-config.yaml
		
kind-down:
		kind delete cluster --name $(KIND_CLUSTER)

kind-status:
		kubectl get nodes -o wide
		kubectl get svc -o wide
		kubectl get pods -o wide --watch --all-namespaces