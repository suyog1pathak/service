# Image URL to use all building/pushing image targets
artifactURL ?= services
TAG ?= latest

build:
	go build -o bin/manager cmd/run-services.go

docker-build:
	docker build -t ${artifactURL}:${TAG} .

docker-push:
	docker push ${artifactURL}:${TAG}