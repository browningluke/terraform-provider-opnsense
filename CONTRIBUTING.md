# Contributor Guides

Thank you for your considering contributing to this Terraform provider! This document provides helpful guidelines and instructions for contributing to this project.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Development Setup](#development-setup)
  - [1. Fork and Clone](#1-fork-and-clone)
  - [2. Local Dependencies](#2-local-dependencies)
  - [3. Install Dependencies](#3-install-dependencies)
  - [4. Build the Provider](#4-build-the-provider)
- [Making Changes](#making-changes)
  - [Branch Naming](#branch-naming)
  - [Code Organization](#code-organization)
  - [Generating Documentation](#generating-documentation)
  - [Code Formatting](#code-formatting)
- [Testing](#testing)
  - [Unit Tests](#unit-tests)
  - [Acceptance Tests](#acceptance-tests)
  - [Writing Tests](#writing-tests)
- [Code Standards](#code-standards)
  - [Go Best Practices](#go-best-practices)
  - [Provider-Specific Standards](#provider-specific-standards)
  - [Documentation Requirements](#documentation-requirements)
  - [Error Messages](#error-messages)
- [Pull Request Process](#pull-request-process)
  - [Before Submitting](#before-submitting)
  - [PR Guidelines](#pr-guidelines)
  - [Review Process](#review-process)
- [Project Structure](#project-structure)
  - [Key Concepts](#key-concepts)
- [Additional Resources](#additional-resources)
- [Getting Help](#getting-help)

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.24 or later** - [Download](https://golang.org/dl/)
- **Terraform 1.10+** - [Download](https://developer.hashicorp.com/terraform/install)
- **OPNsense instance** - Required for running acceptance tests (can be a VM, physical device, or test environment)
- **Make** - For using the provided Makefile commands (optional but recommended)

You should also be familiar with:
- Go programming fundamentals
- Terraform provider development basics
- The Terraform Plugin Framework (v6)
- Basic OPNsense configuration and API usage

## Development Setup

### 1. Fork and Clone

If you are adding a new feature, or modifying an existing resource, fork the `browningluke/opnsense-go` repo on Github, then clone that fork:
```bash
$ git clone https://github.com/YOUR-USERNAME/opnsense-go.git
```

Fork this repository on GitHub, then clone your fork:

```bash
$ git clone https://github.com/YOUR-USERNAME/terraform-provider-opnsense.git
$ cd terraform-provider-opnsense
```

### 2. Local Dependencies

This provider depends *heavily* on [`opnsense-go`](https://github.com/browningluke/opnsense-go) as its API client. Development will be made much easier if you add the following replace directive to your `go.mod` file:
```go
replace github.com/browningluke/opnsense-go => ../opnsense-go
```

The provider will automatically use your local version. If you're only working on the provider, the published version will be used via Go modules.

### 3. Install Dependencies

```bash
$ go mod download
```

### 4. Build the Provider

Build the provider binary:

```bash
# Build to current directory
$ make build

# Or build and install to your local Terraform plugins directory
# This makes the provider available at: dev.io/browningluke/opnsense v1.0.0
$ make build-local
```

After running `make build-local`, you can use the provider in your Terraform configurations with:

```hcl
terraform {
  required_providers {
    opnsense = {
      source  = "dev.io/browningluke/opnsense"
      version = "1.0.0"
    }
  }
}
```

## Making Changes

### Branch Naming

Create a feature branch from `main` for your changes:

```bash
$ git checkout -b feature/your-feature-name
```

Use descriptive branch names that indicate the type of change:
- `feature/` - New features or resources
- `fix/` - Bug fixes
- `docs/` - Documentation changes
- `refactor/` - Code refactoring

### Code Organization

This project follows a **service module pattern**. Each OPNsense "Module" maps to a service, and has its own directory in `internal/service/`:

```
internal/service/[service]/
├── exports.go              # Factory functions
├── [resource]_resource.go  # Resource CRUD implementation
├── [resource]_data_source.go
├── [resource]_schema.go    # Schema + conversion functions
└── [resource]_test.go      # Acceptance tests
```

There are two scenarios for adding new functionality:

#### Adding a Resource to an Existing Service

If the OPNsense API endpoint you're implementing belongs to an existing service module (e.g., adding a new firewall resource to the `firewall` service):

1. **Create the resource files** in the appropriate service directory:
   - `[resource]_resource.go` - Resource CRUD implementation
   - `[resource]_data_source.go` - Data source (if applicable)
   - `[resource]_schema.go` - Schema and conversion functions
   - `[resource]_test.go` - Acceptance tests

2. **Register in `exports.go`**: Add your new factory functions to the **existing** `exports.go` file in the service directory:

   ```go
   func Resources(ctx context.Context) []func() resource.Resource {
       return []func() resource.Resource{
           newAliasResource,
           newCategoryResource,
           newYourNewResource,  // Add your new resource here
       }
   }

   func DataSources(ctx context.Context) []func() datasource.DataSource {
       return []func() datasource.DataSource{
           newAliasDataSource,
           newCategoryDataSource,
           newYourNewDataSource,  // Add your new data source here
       }
   }
   ```

3. **Implement all required interfaces** for your resource/data source
4. **Add comprehensive acceptance tests** (using the [Terraform Provider Framework](https://developer.hashicorp.com/terraform/plugin/framework/acctests))

#### Creating a New Service

If you're implementing an OPNsense API endpoint that doesn't fit into any existing service (e.g., adding support for a new OPNsense module):

1. **Create the service directory**: `internal/service/[newservice]/`

2. **Create the resource files** as described above

3. **Create `exports.go`**: Add a new `exports.go` file that exports the `Resources()` and `DataSources()` functions:

   ```go
   package newservice

   import (
       "context"

       "github.com/hashicorp/terraform-plugin-framework/datasource"
       "github.com/hashicorp/terraform-plugin-framework/resource"
   )

   func Resources(ctx context.Context) []func() resource.Resource {
       return []func() resource.Resource{
           newYourResource,
       }
   }

   func DataSources(ctx context.Context) []func() datasource.DataSource {
       return []func() datasource.DataSource{
           newYourDataSource,
       }
   }
   ```

4. **Register in `provider.go`**: Import your new service and add it to the provider's `Resources()` and `DataSources()` functions in `internal/provider/provider.go`:

   ```go
   import (
       // ... existing imports ...
       "github.com/browningluke/terraform-provider-opnsense/internal/service/newservice"
   )

   func (p *opnsenseProvider) Resources(ctx context.Context) []func() resource.Resource {
       controllers := [][]func() resource.Resource{
           diagnostics.Resources(ctx),
           firewall.Resources(ctx),
           newservice.Resources(ctx),  // Add your new service here
           // ... other services ...
       }

       var resources []func() resource.Resource
       for _, s := range controllers {
           resources = append(resources, s...)
       }
       return resources
   }

   func (p *opnsenseProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
       controllers := [][]func() datasource.DataSource{
           diagnostics.DataSources(ctx),
           firewall.DataSources(ctx),
           newservice.DataSources(ctx),  // Add your new service here
           // ... other services ...
       }

       var dataSources []func() datasource.DataSource
       for _, s := range controllers {
           dataSources = append(dataSources, s...)
       }
       return dataSources
   }
   ```

### Generating Documentation

Documentation is auto-generated from code. After making schema changes, regenerate docs:

```bash
$ make docs
# or
$ go generate ./...
```

Always commit the generated documentation with your changes.

### Code Formatting

Format your code before committing:

```bash
$ make fmt
# or
$ gofmt -w .
```

The CI pipeline will check formatting, so ensure all code is properly formatted.

## Testing

This project uses two types of tests:

### Unit Tests

Standard Go tests for utility functions and isolated logic:

```bash
$ make test
# or
$ go test ./...
```

### Acceptance Tests

Acceptance tests are handled by the Terraform Provider Framework and validate resources against a live OPNsense instance. These tests:
- Create real resources in OPNsense
- Verify resource state
- Test updates and modifications
- Clean up resources after testing

**Environment Variables:**

Set these before running acceptance tests:

```bash
export OPNSENSE_URI="https://your-opnsense-host.example.com"
export OPNSENSE_API_KEY="your-api-key"
export OPNSENSE_API_SECRET="your-api-secret"
export OPNSENSE_ALLOW_INSECURE="true"  # For self-signed certificates
```

**Running Acceptance Tests:**

```bash
# Run all acceptance tests
$ make testacc

# Run tests for a specific package
$ make testacc PKG=firewall

# Run a specific test in a package
$ make testacc PKG=firewall TEST=TestAccFirewallFilterResource

# Alternatively, use the full go test command (IMPORTANT: use -p 1 to run serially)
$ TF_ACC=1 go test -v -p 1 ./...
$ TF_ACC=1 go test -v -p 1 ./internal/service/firewall/...
$ TF_ACC=1 go test -v -p 1 ./internal/service/firewall/ -run TestAccFirewallFilterResource
```

> **Note:** The `-p 1` flag is **required** to run tests serially, avoiding conflicts from concurrent access to the shared OPNsense instance.

### Writing Tests

Acceptance tests should follow this pattern:

```go
func TestAccResourceName(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { acctest.AccPreCheck(t) },
        ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            // Create and Read testing
            {
                Config: testConfig(),
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr("opnsense_resource.test", "attribute", "value"),
                    resource.TestCheckResourceAttrSet("opnsense_resource.test", "id"),
                ),
            },
            // ImportState testing
            {
                ResourceName:      "opnsense_resource.test",
                ImportState:       true,
                ImportStateVerify: true,
            },
            // Update and Read testing
            {
                Config: testConfigUpdated(),
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr("opnsense_resource.test", "attribute", "updated_value"),
                ),
            },
            // Delete testing automatically occurs in TestCase
        },
    })
}
```

All tests should cover:
1. **Create** - Resource creation with valid configuration
2. **Read** - Verify attributes are correctly read
3. **Update** - Modify the resource and verify changes
4. **Delete** - Automatic cleanup verification
5. **Import** - Test `terraform import` functionality

> [!IMPORTANT]
> The OPNsense VM configured in the automated test pipeline only has a **single interface configured (`wan`)**.
>
> This may affect your test expectations or cause differences between your local tests (which may have multiple interfaces) and CI tests. When writing acceptance tests that reference interfaces, be aware of this limitation and ensure your tests work with the minimal `wan`-only configuration.

## Code Standards

### Go Best Practices

- Follow idiomatic Go conventions
- Use meaningful variable and function names
- Write clear comments for exported functions and complex logic
- Handle errors appropriately - avoid panics
- Use the standard library where possible

### Provider-Specific Standards

- **Schema Definitions**: Define schemas in `*_schema.go` files with clear `MarkdownDescription` fields
- **Type Conversions**: Use helper functions from `internal/tools/type_utils.go` for converting OPNsense API types (e.g., string "0"/"1" to bool)
- **Validators**: Use built-in validators or custom validators from `internal/validators/`
- **Resource Naming**: Follow the pattern `opnsense_[service]_[resource]`
- **Import Support**: Always implement `ImportState` for resources where applicable

#### Schema Documentation Standards

When writing `MarkdownDescription` fields for schema attributes:

1. **Source from OPNsense UI**: Use the "help" text from the relevant section in the OPNsense UI as your primary source
2. **Reword for Cohesion**: The help text may need slight rewording to be clear and grammatically correct
3. **Boolean Attributes**: Should start with "When enabled, ..." or "Whether to enable, ..." if the existing help text doesn't already make sense as-is
4. **Default Values**:
   - For **resources**: Append "Defaults to `value`." at the end of the description (use backticks to encapsulate the value)
   - For **data sources**: Do NOT append default values to attribute descriptions

**Example:**
```go
"enabled": schema.BoolAttribute{
    MarkdownDescription: "Enable this firewall filter rule. Defaults to `true`.",
    Optional:            true,
    Computed:            true,
},
"enable_prefetch_key": schema.BoolAttribute{
    MarkdownDescription: "When enabled, DNSKEYs are fetched earlier in the validation process. Defaults to `true`.",
    Optional:            true,
    Computed:            true,
},
"send_certificate_request": schema.StringAttribute{
    MarkdownDescription: "Whether to send a certificate request.",
    Required:            true,
},
```

#### Schema Versioning Requirements

**ALWAYS include `Version: 1` in your resource schema.** This is critical for future-proofing.

When you make changes to a resource schema:
1. Increment the `Version` field in the schema
2. Implement the `resource.ResourceWithUpgradeState` interface
3. Write an `UpgradeState()` function to handle the migration from the old version to the new version

**Example from `filter_resource.go`:**

```go
// Declare the interface implementation
var _ resource.ResourceWithUpgradeState = &filterResource{}

// Implement the UpgradeState function
func (r *filterResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
    schemaV0 := filterResourceSchemaV0()
    return map[int64]resource.StateUpgrader{
        // Upgrade from version 0 to version 1
        0: {
            PriorSchema:   &schemaV0,
            StateUpgrader: upgradeFilterStateV0toV1,
        },
    }
}

// Migration function
func upgradeFilterStateV0toV1(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
    tflog.Info(ctx, "Upgrading filter resource state from v0 to v1")

    var oldState filterResourceModelV0
    resp.Diagnostics.Append(req.State.Get(ctx, &oldState)...)
    // ... migration logic ...
}
```

> [!NOTE]
> **Exception**: Since this provider is v0 (pre-1.0), this rule can sometimes be violated for very small changes, at the discretion of the maintainer(s). However, always err on the side of caution and implement proper versioning when in doubt.

#### Number Attribute Handling

When converting and storing `number` attributes in the state, use **`-1` as the null value** instead of the literal `null`.

This pattern is used because the OPNsense API often represents "unset" numeric values as empty strings, and `-1` provides a clear sentinel value that can be reliably converted.

**Helper Functions:**
- `tools.StringToInt64(s string) int64` - Converts string to int64, returns 0 for empty strings
- `tools.Int64ToStringNegative(i int64) string` - Converts int64 to string, returns "" for -1

**Example:**
```go
// Converting from API to Terraform state
data.Timeout = types.Int64Value(tools.StringToInt64(apiData.Timeout))

// Converting from Terraform state to API
apiData.Timeout = tools.Int64ToStringNegative(data.Timeout.ValueInt64())
```

### Documentation Requirements

- Add clear descriptions to all schema attributes
- Include examples in the `examples/` directory
- Update the README if adding significant features
- Use proper Terraform HCL formatting in examples

### Error Messages

- Provide clear, actionable error messages
- Include context about what operation failed
- Suggest solutions when possible

## Pull Request Process

### Before Submitting

Ensure your PR meets these requirements:

- [ ] Code is formatted with `gofmt`
- [ ] All tests pass locally (`go test ./...`)
- [ ] Acceptance tests pass (if applicable)
- [ ] Documentation is generated (`go generate ./...`)
- [ ] Commit messages are clear and descriptive
- [ ] No unnecessary files are included (check `.gitignore`)

### PR Guidelines

1. **Keep PRs focused** - One feature or fix per PR. Smaller PRs are reviewed faster.
2. **Write a clear title** - Use format: `resource_name - Brief description` (e.g., `firewall_filter - Add support for tcp_flags attribute`)
3. **Describe your changes** - Explain what you changed and why
4. **Link related issues** - Use keywords like "Fixes #123" or "Closes #456"
5. **Update CHANGELOG** - Add an entry describing your change (if applicable)

### Review Process

1. **Automated Checks** - CI runs tests and linters automatically
2. **Maintainer Review** - Project maintainers will review your code
3. **Feedback** - Address any requested changes
4. **Approval** - Once approved, maintainers will merge your PR

Be patient - maintainers may take time to review. Feel free to politely ping if there's no response after a week.

## Project Structure

```
terraform-provider-opnsense/
├── internal/
│   ├── provider/          # Provider implementation
│   ├── service/           # Service modules (firewall, ipsec, etc.)
│   │   └── [service]/
│   │       ├── exports.go
│   │       ├── *_resource.go
│   │       ├── *_data_source.go
│   │       ├── *_schema.go
│   │       └── *_test.go
│   ├── tools/             # Utility functions (type conversions, etc.)
│   ├── validators/        # Custom validators
│   └── acctest/           # Acceptance test helpers
├── docs/                  # Generated documentation
├── examples/              # Terraform configuration examples
├── .github/               # GitHub workflows and templates
├── CONTRIBUTING.md        # This file
├── README.md              # Project overview
├── go.mod                 # Go module definition
└── main.go                # Provider entry point
```

### Key Concepts

- **Resources** - Manage OPNsense configuration (create, read, update, delete)
- **Data Sources** - Read-only access to OPNsense data
- **Schemas** - Define the structure of resources and data sources
- **Conversion Functions** - Transform between Terraform types and OPNsense API types
- **Service Modules** - Logical grouping of related resources

## Additional Resources

- **[Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)** - Official framework documentation
- **[OPNsense API Documentation](https://docs.opnsense.org/development/api.html)** - Complete API reference
- **[Terraform Provider Design Principles](https://developer.hashicorp.com/terraform/plugin/best-practices)** - Best practices for provider development

## Getting Help

- **Issues** - [GitHub Issues](https://github.com/browningluke/terraform-provider-opnsense/issues) for bug reports and feature requests
- **Documentation** - Check existing docs and examples first

When asking for help, please provide:
- Clear description of the problem
- Steps to reproduce (if applicable)
- Relevant code snippets or configuration
- OPNsense version and provider version
- Error messages (full output)

---

Thank you for contributing to the Terraform Provider for OPNsense! Your contributions help make infrastructure-as-code management of OPNsense better for everyone.
