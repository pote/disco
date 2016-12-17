GPM ?= gpm
DEPENDENCIES = $(firstword $(subst :, ,$(GOPATH)))/up-to-date

all: test

$(DEPENDENCIES): Godeps
	$(GPM) get
	touch $@

build: $(DEPENDENCIES)
	go build .
	go vet .

test: build
	go test

.PHONY: test
