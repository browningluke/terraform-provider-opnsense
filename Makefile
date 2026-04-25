.PHONY: mkdocs build-local build fmt test testacc

# Variables for testacc
PKG ?=
TEST ?=

# Generate documentation from code
docs:
	go generate ./...

# Format Go code
fmt:
	gofmt -w .

# Run unit tests
test:
	go test ./...

# Run acceptance tests
# Usage:
#   make testacc                                    - Run all tests
#   make testacc PKG=firewall                       - Run tests for firewall package
#   make testacc PKG=firewall TEST=TestAccFirewall  - Run specific test in firewall package
testacc:
ifdef PKG
	TF_ACC=1 go test -v -p 1 $(if $(TEST),-run $(TEST)) ./internal/service/$(PKG)/...
else
	TF_ACC=1 go test -v -p 1 $(if $(TEST),-run $(TEST)) ./...
endif

# Build provider binary to local Terraform plugins directory
# This installs the provider at: dev.io/browningluke/opnsense v1.0.0
# Detects OS and architecture automatically
build-local:
	go build -o ~/.terraform.d/plugins/dev.io/browningluke/opnsense/1.0.0/$$(go env GOOS)_$$(go env GOARCH)/terraform-provider-opnsense .

# Build provider binary to current directory
build:
	go build -o terraform-provider-opnsense .
