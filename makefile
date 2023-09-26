SHELL := /bin/bash

run:
	go run main.go

build:
	# go build -ldflags "-X main.build=${VCS_REF}"
	go build -ldflags "-X main.build=${BUILD_REF}"

# Building containers

VERSION := 0.3

all: docker-sales

docker-sales:
	docker  build \
		-f leo/docker/Dockerfile.sales-api \
		-t sales-api3:${VERSION} \
		--build-arg VCS_REF=${VERSION} \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%d%H:%M:%SZ"` \
		.

KIND_CLUSTER := leo-cluster 

kind-up:
	kind create cluster \
		--name ${KIND_CLUSTER} \
		--config leo/k8s/kind/kind-config.yaml

kind-down:
		kind delete cluster --name $(KIND_CLUSTER)

kind-status:
		kubectl get nodes -o wide
		kubectl get svc -o wide
		kubectl get pods -o wide --watch --all-namespaces

kind-load-image:
	kind load docker-image sales-api3:${VERSION} --name ${KIND_CLUSTER}

k8s-apply:
	cat leo/k8s/base/sales-pod.yaml | kubectl apply -f -

kustomize-apply:
	kustomize build leo/k8s/kind/sales-pod | kubectl apply -f -


k8s-logs:
	kubectl logs -n leo-service -l app=leo-service --all-containers=true -f --tail=100

k8s-restart-leo-sales:
	kubectl rollout restart deployment leo-sales -n leo-sales 

kind-update: all kind-load-image k8s-restart-leo-sales

kind-describe:
	kubectl describe nodes 
	kubectl describe svc 
	kubectl describe pod -l app=leo-service -n leo-service

tidy:
	go mod tidy 
	go mod vendor