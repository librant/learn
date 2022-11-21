TARGET := cluster-console

GOPATH ?= $(shell go env GOPATH)
GOROOT ?= $(shell go env GOROOT)
GO      = ${GOROOT}/bin/go

ALL: tool BUILD

LINUX:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${TARGET}
BUILD:
	@echo "begin build" && export GO111MODULE=on
	$(GO) build $(GO_BUILD_FLAG) -o ${TARGET}

tool:
	$(GO) fmt ./
	$(GO) vet ./

clean:
	@if [ -f ${TARGET} ]; then rm ${TARGET}; fi