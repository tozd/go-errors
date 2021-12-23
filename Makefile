.PHONY: lint lint-ci fmt fmt-ci test test-ci clean release lint-docs audit encrypt decrypt sops

lint:
	golangci-lint run --timeout 4m --color always

# TODO: Output both formats at the same time, once it is supported.
# See: https://github.com/golangci/golangci-lint/issues/481
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

release:
	npx --yes --package 'release-it@14.11.6' --package 'git+https://github.com/mitar/keep-a-changelog.git#better-gitlab' -- release-it

lint-docs:
	npx --yes --package 'markdownlint-cli@~0.30.0' -- markdownlint --ignore-path .gitignore --ignore testdata/ '**/*.md'

audit:
	go list -json -deps | nancy sleuth --skip-update-check

encrypt:
	gitlab-config sops -- --encrypt --mac-only-encrypted --in-place --encrypted-comment-regex sops:enc .gitlab-conf.yml

decrypt:
	SOPS_AGE_KEY_FILE=keys.txt gitlab-config sops -- --decrypt --in-place .gitlab-conf.yml

sops:
	SOPS_AGE_KEY_FILE=keys.txt gitlab-config sops -- .gitlab-conf.yml
