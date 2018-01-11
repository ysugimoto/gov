.PHONY: all

all:
	go build -ldflags "-X main.appVersion=$(shell gov)"
