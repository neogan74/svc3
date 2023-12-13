SHELL := /bin/bash

run:
	go run ./app/services/sales-api/main.go | go run ./app/services/tooling/logfmt/main.go  

admin-genkey:
	go run app/tooling/admin/main.go

build:
	# go build -ldflags "-X main.build=${VCS_REF}"
	go build -ldflags "-X main.build=${BUILD_REF}" ./app/services/sales-api/main.go

# Building containers

VERSION := 1.2

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
	cd leo/k8s/kind/sales-pod; kustomize edit set image sales-api-image=sales-api3:${VERSION}
	kind load docker-image sales-api3:${VERSION} --name ${KIND_CLUSTER}

k8s-apply:
	cat leo/k8s/base/sales-pod.yaml | kubectl apply -f -

kustomize-apply:
	kustomize build leo/k8s/kind/sales-pod | kubectl apply -f -


k8s-logs:
	kubectl logs -n leo-sales -l app=leo-sales --all-containers=true -f --tail=100

k8s-logs-pretty:
	kubectl logs -n leo-sales -l app=leo-sales --all-containers=true -f --tail=100 | go run ./app/services/tooling/logfmt/main.go 

k8s-restart-leo-sales:
	kubectl rollout restart deployment leo-sales -n leo-sales 

kind-update: all kind-load-image kustomize-apply k8s-restart-leo-sales

kind-update-apply: all kind-update kustomize-apply

kind-describe:
	kubectl describe nodes 
	kubectl describe svc 
	kubectl describe pod -l app=leo-sales -n leo-service

tidy:
	go mod tidy 
	go mod vendor

expvarsmon:
	~/go/bin/expvarmon -ports=":4000"
metrics:
	~/go/bin/expvarmon -ports="localhost:4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"


install-expvarmon:
	go install github.com/divan/expvarmon@latest

stress:
	hey -m GET -c 100 -n 10000 http://localhost:3000/v1/test