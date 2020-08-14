package template

const Makefile = `appname = {{.Project}}
img =  {{.DockerImg}}:$(version)
imgdev = {{.DockerImg}}dev:$(version)

export docker_image ?= $(imgdev)

wd = $(shell pwd)

composerundev = docker-compose run --rm {{.Project}}

dockerrunbase = docker run --rm $(vols)
rundev = $(dockerrunbase) $(imgdev)

cov = coverage.out
covhtml = coverage.html

pkg ?= ./...
testtimeout ?= 30s
testflag ?= -race -timeout $(timeout) -coverprofile=$(cov) $(flag)
gotest = go test -failfast $(pkg) $(testflag) $(if $(testcase),-run "$(testcase)")

all: static-analysis test test-integration dev-build build

.PHONY: help
help: ## display this help
	@ echo "Please use 'make <target>' where <target> is one of:"
	@ echo
	@ grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-16s\033[0m - %s\n", $$1, $$2}'
	@ echo

.PHONY: build
build: ## build final image
	docker build --build-arg VERSION=$(version) . -t $(img)

.PHONY: dev-build
dev-build: ## build dev image
	docker build . --target base -t $(imgdev)

.PHONY: test
test: dev-build ## Run unit tests, set testcase=<testcase> or flag=-v if you need them
	$(composerundev) $(gotest)

.PHONY: test-integration
test-integration: override testflag+=-tags=integration ## Run integration tests, set pkg=<pkg to test> testcase=<testcase> or flag=-v if you need them
test-integration: dev-build
	$(composerundev) $(gotest)

.PHONY: coverage
coverage: override vols+=-v $(wd):/app ## show test coverage
coverage: test
	$(rundev) go tool cover -html=$(cov) -o=$(covhtml)
	xdg-open coverage.html

.PHONY: lint
lint: override vols+=-v $(wd):/app ## run golang-ci lint
lint:
	$(dockerrunbase) -w /app golangci/golangci-lint:{{.CILintVersion}} \
		golangci-lint run --color always --enable-all \
		./...

.PHONY: fmt
fmt: dev-build ## run gofmt
	$(rundev) gofmt -w -s -l .

.PHONY: static-analysis
static-analysis: fmt lint ## run gofmt and golangci-lint

.PHONY: run
run: override docker_image=$(img)
run: override docker_params+=--service-ports
run: build  ## run the code with given params
	$(composerundev)
`
