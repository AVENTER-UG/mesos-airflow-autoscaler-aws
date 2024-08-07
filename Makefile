#Dockerfile vars

#vars
IMAGENAME=mesos-airflow-autoscaler-aws
REPO=avhost
BRANCH=$(shell git symbolic-ref --short HEAD | xargs basename)
TAG=$(shell git describe)
BUILDDATE=$(shell date -u +%Y%m%d)
IMAGEFULLNAME=${REPO}/${IMAGENAME}
LASTCOMMIT=$(shell git log -1 --pretty=short | tail -n 1 | tr -d " " | tr -d "UPDATE:")

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

ifeq (${BRANCH}, master) 
        BRANCH=latest
endif

ifneq ($(shell echo $(LASTCOMMIT) | grep -E '^v([0-9]+\.){0,2}(\*|[0-9]+)'),)
	BRANCH=${LASTCOMMIT}
endif

build:
	@echo ">>>> Build Docker branch: " ${BRANCH}_${BUILDDATE}
	@docker build --build-arg TAG=${TAG} --build-arg BUILDDATE=${BUILDDATE} -t ${IMAGEFULLNAME}:${BRANCH}_${BUILDDATE} .

build-bin:
	@echo ">>>> Build binary"
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X main.BuildVersion=${BUILDDATE} -X main.GitVersion=${TAG} -extldflags \"-static\"" .

push:
	@echo ">>>> Publish docker image: " ${BRANCH}_${BUILDDATE}
	@docker buildx create --use --name buildkit	
	@docker buildx build --sbom=true --provenance=true --push --platform linux/amd64 --build-arg TAG=${TAG} --build-arg BUILDDATE=${BUILDDATE} -t ${IMAGEFULLNAME}:${BRANCH}_${BUILDDATE} .
	@docker buildx build --sbom=true --provenance=true --push --platform linux/amd64 --build-arg TAG=${TAG} --build-arg BUILDDATE=${BUILDDATE} -t ${IMAGEFULLNAME}:${BRANCH} .
	@docker buildx build --sbom=true --provenance=true --push --platform linux/amd64 --build-arg TAG=${TAG} --build-arg BUILDDATE=${BUILDDATE} -t ${IMAGEFULLNAME}:latest .
	@docker buildx rm buildkit	

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
	grype --add-cpes-if-none .

go-fmt:
	@gofmt -w .

version:
	@echo ">>>> Generate version file"
	@echo "[{ \"version\":\"${TAG}\", \"builddate\":\"${BUILDDATE}\" }]" > .version.json
	@cat .version.json
	@echo "Saved under .version.json"

imagecheck: build
	trivy image ${IMAGEFULLNAME}:${BRANCH}_${BUILDDATE}		

check: go-fmt sboom seccheck imagecheck
all: check build version sboom 
