SHELL := bash

REPOSITORY ?= localhost
CONTAINER_NAME ?= netrouse
TAG ?= latest

# Build all binaries
build: build-cli build-gui

# Build the binary
build-cli:
	hack/build.sh

# Build the GUI
build-gui: tools
	hack/fyne-metadata.sh
	"$(shell pwd)/bin/fyne" build -o "$(shell pwd)/bin/netrouse-gui" --release ./cmd/gui

# Run the server on port 8080 to quickly test changes
run: build-cli
	bin/netrouse server --log debug

# Build the container image
image:
	podman build -t $(REPOSITORY)/$(CONTAINER_NAME):$(TAG) .

# Build all artifacts used for release, except the container images
release:
	hack/containerized hack/release.sh

# Run unit-tests with race detection and coverage
test:
	go test -v -race -coverprofile=coverprofile.out -coverpkg "./..." ./...

# Update project dependencies
update-deps:
	hack/update-deps.sh

# Generate coverage profile
coverprofile:
	hack/coverprofile.sh

# Run linter
lint:
	golangci-lint run -v --timeout 300s

# Format the code
fmt:
	gofmt -s -w ./cmd ./pkg

# Validate that all generated files are up to date.
validate:
	hack/validate.sh

# Generate all required files
generate: generate-bootstrap generate-swagger

# Generate the bootstrap.css file
generate-bootstrap:
	hack/generate-bootstrap.sh

# Generate Swagger documentation
generate-swagger:
	hack/swagger.sh

# Scan code for vulnerabilities using gosec
gosec:
	gosec ./...

# Clean up generated files
clean:
	hack/clean.sh

# Install the tools required for building the app
tools:
	GOBIN="$(shell pwd)/bin" go install tool

# Show this help message
help:
	@echo "Available targets:"
	@echo ""
	@awk '/^#/{c=substr($$0,3);next}c&&/^[[:alpha:]][[:alnum:]_-]+:/{print substr($$1,1,index($$1,":")),c}1{c=0}' $(MAKEFILE_LIST) | column -s: -t
	@echo ""
	@echo "Run 'make <target>' to execute a specific target."

.PHONY: \
	build \
	build-cli \
	build-gui \
	run \
	image \
	release \
	test \
	update-deps \
	coverprofile \
	lint \
	fmt \
	validate \
	generate \
	generate-bootstrap \
	generate-swagger \
	gosec \
	clean \
	tools \
	help \
	$(NULL)
