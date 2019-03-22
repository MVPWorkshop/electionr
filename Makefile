DEP := $(shell command -v dep 2> /dev/null)

install:
	go install ./cmd/legalerd
	go install ./cmd/legalercli

test:
	go test ./...
