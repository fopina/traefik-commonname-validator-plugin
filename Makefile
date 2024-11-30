.PHONY: lint test vendor clean

export GO111MODULE=on

default: lint test

bootstrap: ## install build deps
	go generate -tags tools tools/tools.go

lint:
	golangci-lint run

test:
	go test -v -cover ./...

yaegi_test:
	yaegi test -v .

vendor:
	go mod vendor

clean:
	rm -rf ./vendor