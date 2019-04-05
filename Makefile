DEP := $(shell command -v dep 2> /dev/null)

update_api_docs:
	@statik -src=cmd/electionrcli/swagger-ui -dest=cmd/electionrcli -f

install:
	go install ./cmd/electionrd
	go install ./cmd/electionrcli
	go install ./cmd/gaiadebug

test:
	go test ./...
