VERSION := $(shell cat VERSION)

PKG:= kcl-lang.io/cli
LDFLAGS := -X $(PKG)/pkg/version.version=$(VERSION)
TAGS := rpc
COVER_FILE			?= coverage.out
SOURCE_PATHS		?= ./pkg/...
MAIN_FILE := ./cmd/kcl/main.go

GO ?= go

.PHONY: run
run:
	go run $(MAIN_FILE) run ./examples/kubernetes.k

.PHONY: format
format:
	test -z "$$(find . -type f -o -name '*.go' -exec gofmt -d {} + | tee /dev/stderr)" || \
	test -z "$$(find . -type f -o -name '*.go' -exec gofmt -w {} + | tee /dev/stderr)"

.PHONY: lint
lint:
	scripts/update-gofmt.sh
	scripts/verify-gofmt.sh
	scripts/verify-govet.sh

.PHONY: build
build: lint
	mkdir -p bin/
	$(GO_BUILD_ENV) go build -o bin/kcl -ldflags="$(LDFLAGS)" $(MAIN_FILE)

.PHONY: test
test:
	go test -v ./...

.PHONY: e2e-test
e2e-test:
	./examples/test.sh
	go test -gcflags=all=-l -timeout=20m `go list ./cmd/...` -coverprofile $(COVER_FILE) ${TEST_FLAGS} -v

.PHONY: cover
cover: ## Generates coverage report
	go test -gcflags=all=-l -timeout=20m `go list $(SOURCE_PATHS)` -coverprofile $(COVER_FILE) ${TEST_FLAGS} -v

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: bootstrap
bootstrap:
	go mod download
	command -v golint || GO111MODULE=off go get -u golang.org/x/lint/golint

.PHONY: docker-run-release
docker-run-release: export pkg=/go/src/github.com/kcl-lang/cli
docker-run-release:
	git checkout main
	git push
	docker run -it --rm -e GITHUB_TOKEN -v $(shell pwd):$(pkg) -w $(pkg) golang:1.24 make bootstrap release

.PHONY: dist
dist: export COPYFILE_DISABLE=1 #teach OSX tar to not put ._* files in tar archive
dist: export CGO_ENABLED=0
dist:
	rm -rf build/kcl/* release/*
	mkdir -p build/kcl/bin release/
	cp -f README.md LICENSE build/kcl
	GOOS=linux GOARCH=amd64 $(GO) build -o build/kcl/bin/kcl -trimpath -ldflags="$(LDFLAGS)" -tags="$(TAGS)" $(MAIN_FILE)
	tar -C build/ -zcvf $(CURDIR)/release/kcl-linux-amd64.tgz kcl/
	GOOS=linux GOARCH=arm64 $(GO) build -o build/kcl/bin/kcl -trimpath -ldflags="$(LDFLAGS)" -tags="$(TAGS)" $(MAIN_FILE)
	tar -C build/ -zcvf $(CURDIR)/release/kcl-linux-arm64.tgz kcl/
	GOOS=darwin GOARCH=amd64 $(GO) build -o build/kcl/bin/kcl -trimpath -ldflags="$(LDFLAGS)" -tags="$(TAGS)" $(MAIN_FILE)
	tar -C build/ -zcvf $(CURDIR)/release/kcl-macos-amd64.tgz kcl/
	GOOS=darwin GOARCH=arm64 $(GO) build -o build/kcl/bin/kcl -trimpath -ldflags="$(LDFLAGS)" -tags="$(TAGS)" $(MAIN_FILE)
	tar -C build/ -zcvf $(CURDIR)/release/kcl-macos-arm64.tgz kcl/
	rm build/kcl/bin/kcl
	GOOS=windows GOARCH=amd64 $(GO) build -o build/kcl/bin/kcl.exe -trimpath -ldflags="$(LDFLAGS)" -tags="$(TAGS)" $(MAIN_FILE)
	tar -C build/ -zcvf $(CURDIR)/release/kcl-windows-amd64.tgz kcl/

.PHONY: release
release: lint dist
	scripts/release.sh v$(VERSION)

.PHONY: tag
tag:
	scripts/tag.sh v$(VERSION)

.PHONY: e2e
e2e: ## Run e2e test
	scripts/e2e/e2e.sh

.PHONY: e2e-init
e2e-init:
	scripts/e2e/e2e-init.sh $(TS)

GO_VERSION := $(shell awk '/^go /{print $$2}' go.mod|head -n1)

GIT_TAG ?= $(shell git describe --tags --dirty --always)
# Image URL to use all building/pushing image targets
PLATFORMS ?= linux/amd64,linux/arm64
DOCKER_BUILDX_CMD ?= docker buildx
IMAGE_BUILD_CMD ?= $(DOCKER_BUILDX_CMD) build
BASE_IMAGE ?= debian:12-slim
BUILDER_IMAGE ?= golang:$(GO_VERSION)
CGO_ENABLED ?= 0

IMAGE_BUILD_EXTRA_OPTS ?=

IMAGE_REGISTRY ?= ghcr.io/kcl-lang
IMAGE_NAME := kcl
IMAGE_REPO ?= $(IMAGE_REGISTRY)/$(IMAGE_NAME)
IMAGE_TAG ?= $(IMAGE_REPO):$(GIT_TAG)

ifdef EXTRA_TAG
IMAGE_EXTRA_TAG ?= $(IMAGE_REPO):$(EXTRA_TAG)
endif
ifdef IMAGE_EXTRA_TAG
IMAGE_BUILD_EXTRA_OPTS += -t $(IMAGE_EXTRA_TAG)
endif

# Build the multiplatform container image locally and push to repo.
.PHONY: image-local-push
image-local-push: PUSH=--push
image-local-push: image-local-build

# Build the multiplatform container image locally.
.PHONY: image-local-build
image-local-build:
	BUILDER=$(shell $(DOCKER_BUILDX_CMD) create --use)
	$(MAKE) image-build PUSH=$(PUSH)
	$(DOCKER_BUILDX_CMD) rm $$BUILDER

.PHONY: image-push
image-push: PUSH=--push
image-push: image-build

image-build:
	$(IMAGE_BUILD_CMD) -t $(IMAGE_TAG) \
		--platform=$(PLATFORMS) \
		--build-arg BASE_IMAGE=$(BASE_IMAGE) \
		--build-arg BUILDER_IMAGE=$(BUILDER_IMAGE) \
		--build-arg CGO_ENABLED=$(CGO_ENABLED) \
		$(PUSH) \
		$(IMAGE_BUILD_EXTRA_OPTS) ./
