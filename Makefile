.PHONY: all

all:
	GO111MODULE=on go build -ldflags "-X main.appVersion=$(shell gov)"
