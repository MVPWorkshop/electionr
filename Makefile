DEP := $(shell command -v dep 2> /dev/null)

install:
	go install ./cmd/legalerd
	go install ./cmd/legalercli
	go install ./cmd/gaiadebug

test:
	go test ./...
