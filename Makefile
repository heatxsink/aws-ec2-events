GOROOT ?= /usr/local/go
GOBIN ?= go
GOPATH = $(shell pwd)/lib
GO_BIN_STOCK = GOPATH=$(GOPATH) $(GOBIN)
GO_BIN_LINUX = GOARCH=386 GOOS=linux $(GO_BIN_STOCK)
#GO_BIN_TARGET = $(GO_BIN_LINUX)
GO_BIN_TARGET = $(GO_BIN_STOCK)
GODEPS = github.com/heatxsink/goamz/aws \
github.com/heatxsink/goamz/aws \
github.com/heatxsink/go-colour \
github.com/heatxsink/go-yeam \
github.com/heatxsink/go-simpleconfig
EVENTS_SRC = src
EVENTS_BIN = bin
EVENTS_APP = $(EVENTS_SRC)/aws-ec2-events.go
EVENTS_APPS = $(EVENTS_BIN)/aws-ec2-events

all: $(EVENTS_APPS)

$(EVENTS_BIN)/aws-ec2-events: deps
	$(GO_BIN_TARGET) build $(EVENTS_APP)
	mv aws-ec2-events $(EVENTS_BIN)

$(GOPATH)/src/%:
	$(GO_BIN_TARGET) get $*

%_test.go:
	$(GO_BIN_TARGET) test $<

deps: fix-gopath $(patsubst %, $(GOPATH)/src/%, $(GODEPS))

fix-gopath:
	[ -d lib ] || mkdir lib
	[ -d bin ] || mkdir bin
clean:
	rm -rf lib bin
