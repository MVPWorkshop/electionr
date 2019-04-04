DEP := $(shell command -v dep 2> /dev/null)

install:
	go install ./cmd/electionrd
	go install ./cmd/electionrcli
	go install ./cmd/gaiadebug

test:
	go test ./...
