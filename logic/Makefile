include platform.mk

NAME=logic
IMAGE_NAME=zim/$(NAME)
GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_TAG=$(shell git describe --abbrev=0 --tags --always --match "v*")
GIT_IMPORT=github.com/fztcjjl/zim/logic/cmd
BUILD_DATE=$(shell date +%s)
LDFLAGS=-X $(GIT_IMPORT).BuildDate=$(BUILD_DATE) -X $(GIT_IMPORT).GitCommit=$(GIT_COMMIT) -X $(GIT_IMPORT).GitTag=$(GIT_TAG)
IMAGE_TAG=$(GIT_TAG)-$(GIT_COMMIT)

.PHONY: build clean vet test docker
build:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) go build -a -installsuffix cgo -ldflags "-s -w ${LDFLAGS}"

docker:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(IMAGE_NAME):latest
#	docker push $(IMAGE_NAME):$(IMAGE_TAG)
#	docker push $(IMAGE_NAME):latest

vet:
	go vet ./...

test: vet
	go test -v ./...

clean:
	rm -rf ./logic

