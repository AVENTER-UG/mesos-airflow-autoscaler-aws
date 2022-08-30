#Dockerfile vars

#vars
IMAGENAME=mesos-airflow-autoscaler-aws
REPO=avhost
BRANCH=${shell git rev-parse --abbrev-ref HEAD}
TAG=latest
BUILDDATE=${shell date -u +%Y-%m-%dT%H:%M:%SZ}
IMAGEFULLNAME=${REPO}/${IMAGENAME}

.PHONY: help build all docs

help:
	    @echo "Makefile arguments:"
	    @echo ""
	    @echo "Makefile commands:"
	    @echo "build"
	    @echo "all"
			@echo "docs"
			@echo "publish"
			@echo ${TAG}

.DEFAULT_GOAL := all

build:
	@echo ">>>> Build Docker"
	@docker build --build-arg TAG=${TAG} --build-arg BUILDDATE=${BUILDDATE} -t ${IMAGEFULLNAME}:${TAG} .

build-bin:
	@echo ">>>> Build binary"
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X main.BuildVersion=${BUILDDATE} -X main.GitVersion=${TAG} -extldflags \"-static\"" .

publish:
	@echo ">>>> Publish docker image"
	@docker buildx build --push --platform linux/arm64,linux/amd64 --build-arg TAG=${TAG} --build-arg BUILDDATE=${BUILDDATE} -t ${IMAGEFULLNAME}:${TAG} .
	@docker buildx build --push --platform linux/arm64,linux/amd64 --build-arg TAG=${TAG} --build-arg BUILDDATE=${BUILDDATE} -t ${IMAGEFULLNAME}:latest .

update-gomod:
	go get -u
	go mod tidy

docs:
	@echo ">>>> Build docs"
	$(MAKE) -C $@

sboom:
	syft dir:. > sbom.txt
	syft dir:. -o json > sbom.json

seccheck:
	gosec --exclude G104 --exclude-dir ./vendor ./... 

version:
	@echo ">>>> Generate version file"
	@echo "[{ \"version\":\"${TAG}\", \"builddate\":\"${BUILDDATE}\" }]" > .version.json
	@cat .version.json
	@echo "Saved under .version.json"


all: build version
