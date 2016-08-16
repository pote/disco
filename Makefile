GPM ?= gpm
DEPENDENCIES = $(firstword $(subst :, ,$(GOPATH)))/up-to-date

all: test

$(DEPENDENCIES): Godeps
	$(GPM) get
	touch $@

test: $(DEPENDENCIES)
	go test

.PHONY: test
