.PHONY: lint lint-ci fmt fmt-ci test test-ci clean

lint:
	golangci-lint run --timeout 4m --color always

lint-ci:
	-golangci-lint run --timeout 4m --color always
	golangci-lint run --timeout 4m --out-format code-climate > codeclimate.json

fmt:
	go mod tidy
	gofumpt -w *.go
	goimports -w -local gitlab.com/tozd/go/errors *.go

fmt-ci: fmt
	git diff --exit-code --color=always

test:
	gotestsum --format pkgname --packages ./... -- -race -timeout 10m -cover -covermode atomic

test-ci:
	gotestsum --format pkgname --packages ./... --junitfile tests.xml -- -race -timeout 10m -coverprofile=coverage.txt -covermode atomic
	gocover-cobertura < coverage.txt > coverage.xml
	go tool cover -html=coverage.txt -o coverage.html

clean:
	rm -f coverage.* codeclimate.json tests.xml
