.PHONY: generate lint breaking test check

generate:
	buf generate

lint:
	buf lint

breaking:
	buf breaking --against '.git#branch=master'

test:
	go test ./...

check: lint breaking test
