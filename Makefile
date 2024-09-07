.DEFAULT_GOAL := help
PROJECT_BIN = $(shell pwd)/bin
$(shell [ -f bin ] || mkdir -p $(PROJECT_BIN))
GOBIN = go
PATH := $(PROJECT_BIN):$(PATH)
GOARCH = amd64
LINUX_LDFLAGS = -extldflags '-static' -w -s -buildid=
WINDOWS_LDFLAGS = -extldflags '-static' -w -s -buildid=
GCFLAGS = "all=-trimpath=$(shell pwd) -dwarf=false -l"
ASMFLAGS = "all=-trimpath=$(shell pwd)"
APP = drw6

build: build-bins .crop ## Build all

release: build-bins .crop zip ## Build release

deliv: release
	cp $(PROJECT_BIN)/$(APP).zip ~/dev/windows/ess6/shared

zip:
	rm -rf repository/10-drwbases/common/*
	cp $(PROJECT_BIN)/$(APP).exe .
	zip -rq9 $(PROJECT_BIN)/$(APP).zip $(APP).exe repository
	rm $(APP).exe

docker: ## Build with docker
	docker compose up --build --force-recreate || docker-compose up --build --force-recreate


build-bins: ## Build for windows
	CGO_ENABLED=0 GOOS=windows GOARCH=$(GOARCH) \
	  $(GOBIN) build -ldflags="$(WINDOWS_LDFLAGS)" -trimpath -gcflags=$(GCFLAGS) -asmflags=$(ASMFLAGS) \
	  -o $(PROJECT_BIN)/$(APP).exe cmd/$(APP)/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) \
	  $(GOBIN) build -ldflags="$(LINUX_LDFLAGS)" -trimpath -gcflags=$(GCFLAGS) -asmflags=$(ASMFLAGS) \
	  -o $(PROJECT_BIN)/$(APP) cmd/$(APP)/main.go


	
.crop:
	strip $(PROJECT_BIN)/$(APP)
	objcopy --strip-unneeded $(PROJECT_BIN)/$(APP)
	strip $(PROJECT_BIN)/$(APP).exe
	objcopy --strip-unneeded $(PROJECT_BIN)/$(APP).exe

dev:
	find . -name "*.go" | entr -r make build

help:
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

cert: ## Make cert for sign binaries
	openssl req -newkey rsa:2048 -nodes -keyout server.key -out server.csr -subj "/CN=your.hostname.here"
	echo "subjectAltName = IP:0.0.0.0,IP:127.0.0.1,IP:172.17.0.1,IP:172.252.212.8,IP:45.120.177.178,IP:88.151.117.196" > extfile.cnf
	openssl x509 -req -sha256 -days 365 -in server.csr -signkey server.key -out server.crt -extfile extfile.cnf

	cp server.key cmd/drw6
	cp server.crt cmd/drw6

	rm server.csr \
	   server.key \
	   server.crt \
	   extfile.cnf

