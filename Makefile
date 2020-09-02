APP:=mailsender
SHORT:=mailsender
APP_ENTRY_POINT:=cmd/mailsender.go
BUILD_OUT_DIR:=./
COMMON_PATH	?= $(shell pwd)

DOCKER_REGISTRY=andskur

GOOS	:=
GOARCH	:=
ifeq ($(OS),Windows_NT)
	GOOS =windows
	ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
		OSFLAG =amd64
	endif
	ifeq ($(PROCESSOR_ARCHITECTURE),x86)
		OSFLAG =ia32
	endif
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		GOOS =linux
	endif
	ifeq ($(UNAME_S),Darwin)
		GOOS =darwin
	endif
		UNAME_P := $(shell uname -m)
	ifeq ($(UNAME_P),x86_64)
		GOARCH =amd64
	endif
	ifneq ($(filter %86,$(UNAME_P)),)
		GOARCH =386
	endif
	ifneq ($(filter arm%,$(UNAME_P)),)
		GOARCH =arm
	endif
endif

TAG 		:= $(shell git describe --abbrev=0 --tags)
COMMIT		:= $(shell git rev-parse HEAD)
BRANCH		?= $(shell git rev-parse --abbrev-ref HEAD)
REMOTE		:= $(shell git config --get remote.origin.url)
BUILD_DATE	:= $(shell date +'%Y-%m-%dT%H:%M:%SZ%Z')

RELEASE :=
ifeq ($(TAG),)
	RELEASE := $(COMMIT)
else
	RELEASE := $(TAG)
endif

CONTAINER_IMAGE := $(DOCKER_REGISTRY)/$(APP):$(RELEASE)

LDFLAGS += -X $(GITVER_PKG).ServiceName=$(SHORT)
LDFLAGS += -X $(GITVER_PKG).CommitTag=$(TAG)
LDFLAGS += -X $(GITVER_PKG).CommitSHA=$(COMMIT)
LDFLAGS += -X $(GITVER_PKG).CommitBranch=$(BRANCH)
LDFLAGS += -X $(GITVER_PKG).OriginUrl=$(REMOTE)
LDFLAGS += -X $(GITVER_PKG).BuildDate=$(BUILD_DATE)

all: clean build

run: tidy
	go run -race $(APP_ENTRY_POINT) serve

build:
	env CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_OUT_DIR)/$(APP) $(APP_ENTRY_POINT)

clean:
	rm -f $(APP)

tidy:
	GOPRIVATE=$(GOPRIVATE) go mod tidy

update:
	GOPRIVATE=$(GOPRIVATE) go get -u ./...

image: GOOS =linux
image: GOARCH =amd64
image: build
	docker build -t $(CONTAINER_IMAGE) .

tag:
	docker tag $(CONTAINER_IMAGE) $(DOCKER_REGISTRY)/$(APP):latest

image_latest: image tag

push: image tag
	docker push $(DOCKER_REGISTRY)/$(APP):latest

container: GOOS =linux
container: GOARCH =amd64
container: image
	docker stop $(CONTAINER_IMAGE) || true && docker rm $(CONTAINER_IMAGE) || true
	docker run --name --rm $(CONTAINER_IMAGE)