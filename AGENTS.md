# Agent Instructions

## Build, Test, and Lint

```bash
# Build provider binary to current directory
make build

# Build and install to local Terraform plugins directory (dev.io/browningluke/opnsense v1.0.0)
make build-local

# Run unit tests
make test
go test ./...

# Run all acceptance tests (requires live OPNsense instance)
make testacc

# Run acceptance tests for a specific package
make testacc PKG=firewall

# Run a specific acceptance test
make testacc PKG=firewall TEST=TestAccFirewallFilterResource

# Regenerate documentation (after schema changes)
make docs

# Format code
make fmt
```

Acceptance tests require these environment variables:
```bash
export OPNSENSE_URI="https://your-opnsense-host.example.com"
export OPNSENSE_API_KEY="your-api-key"
export OPNSENSE_API_SECRET="your-api-secret"
export OPNSENSE_ALLOW_INSECURE="true"  # for self-signed certs
```

Always use `-p 1` when running `go test` directly for acceptance tests — the shared OPNsense instance cannot handle concurrent test runs.

The CI OPNsense VM only has a single `wan` interface. Tests referencing interfaces must work with this minimal setup.

## Architecture

This is a Terraform Plugin Framework (v6) provider for OPNsense. The API client lives in the separate [`opnsense-go`](https://github.com/browningluke/opnsense-go) module (`github.com/browningluke/opnsense-go`). To develop against a local copy, add to `go.mod`:
```go
replace github.com/browningluke/opnsense-go => ../opnsense-go
```

### Service Module Pattern

Each OPNsense module maps to a service package under `internal/service/<service>/`:

```
internal/service/<service>/
├── exports.go               # Resources() and DataSources() factory functions
├── <resource>_resource.go   # CRUD implementation
├── <resource>_data_source.go
├── <resource>_schema.go     # Schema definition + convert functions (SchemaToStruct / StructToSchema)
└── <resource>_test.go       # Acceptance tests
```

`internal/provider/provider.go` aggregates all services — each service's `Resources(ctx)` and `DataSources(ctx)` are appended into the provider's flat list.

### Adding a New Resource to an Existing Service

1. Create `<resource>_resource.go`, `<resource>_data_source.go`, `<resource>_schema.go`, `<resource>_test.go`
2. Register factory functions in the service's existing `exports.go`

### Adding a New Service

1. Create `internal/service/<newservice>/` with all resource files plus `exports.go`
2. Import the package and append its `Resources(ctx)` / `DataSources(ctx)` in `internal/provider/provider.go`

## Key Conventions

### File Roles

- `*_schema.go` — contains the resource model struct, `*ResourceSchema()` / `*DataSourceSchema()`, and the two conversion functions `convert*SchemaToStruct` / `convert*StructToSchema`. No CRUD logic here.
- `*_resource.go` — implements `resource.Resource`, `resource.ResourceWithConfigure`, `resource.ResourceWithImportState` (and optionally `resource.ResourceWithUpgradeState`, `resource.ResourceWithConfigValidators`). CRUD calls the conversion functions from `_schema.go`.
- `exports.go` — only exports `Resources()` and `DataSources()` slices; constructor functions (`new*Resource`, `new*DataSource`) are unexported.

### Resource Naming

Resources are named `opnsense_<service>_<resource>` (e.g., `opnsense_firewall_alias`).

### Schema Versioning

**Always include `Version: 1` in new resource schemas.** When a schema changes:
1. Increment `Version`
2. Implement `resource.ResourceWithUpgradeState`
3. Write an `UpgradeState()` map with a migration function per version transition

The `id` attribute is always `Computed` with `stringplanmodifier.UseStateForUnknown()`.

### Type Conversion Helpers (`internal/tools/type_utils.go`)

The OPNsense API represents booleans as `"0"`/`"1"` and unset numbers as empty strings. Use:

| Conversion | Function |
|---|---|
| `bool` → API string | `tools.BoolToString(b)` |
| API string → `bool` | `tools.StringToBool(s)` |
| `int64` → API string (`-1` → `""`) | `tools.Int64ToStringNegative(i)` |
| API string → `int64` (`""` → `-1`) | `tools.StringToInt64(s)` |
| `float64` → API string (`-1` → `""`) | `tools.Float64ToStringNegative(f)` |
| empty string → `types.String` null | `tools.StringOrNull(s)` |
| `[]string` → `types.Set` (deduped, skips `""`) | `tools.StringSliceToSet(s)` |

Use **`-1` as the sentinel null value** for numeric attributes — never store literal `null` for numbers.

### Renamed OPNsense API Fields (dual-key writes)
When an OPNsense release renames a model field (e.g. 26.1.10 renamed the route model's `disabled` to `enabled`), opnsense-go keeps **both** fields on the struct. In the conversion functions, set both on writes (OPNsense ignores the key it doesn't know, so this works on either side of the rename) and prefer the new field on reads, falling back to the old one when empty:

```go
// SchemaToStruct: set both
Enabled:  tools.BoolToString(d.Enabled.ValueBool()),
Disabled: tools.BoolToString(!d.Enabled.ValueBool()),

// StructToSchema: new key wins, old key is the fallback
enabled := tools.StringToBool(d.Enabled)
if d.Enabled == "" {
    enabled = !tools.StringToBool(d.Disabled)
}
```

`internal/service/routes/route_schema.go` is the reference example. Cover both response shapes with unit tests (see `route_schema_test.go`) — CI only runs against one OPNsense version.

### Schema Documentation

- Source `MarkdownDescription` text from the OPNsense UI help text for that field
- Boolean attributes: start with `"When enabled, ..."` or `"Whether to enable, ..."`
- Resources: append `"Defaults to \`value\`."` to optional/computed attribute descriptions
- Data sources: **omit** default value notes from descriptions

### `NotFoundError` Handling in Read

When `Read` receives a `*errs.NotFoundError`, call `resp.State.RemoveResource(ctx)` and return — do not add an error diagnostic. This lets Terraform detect and recreate deleted-out-of-band resources.

### Acceptance Test Pattern

```go
func TestAccResourceName(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { acctest.AccPreCheck(t) },
        ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {Config: testConfig(), Check: ...},
            {ResourceName: "opnsense_x.test", ImportState: true, ImportStateVerify: true},
            {Config: testConfigUpdated(), Check: ...},
        },
    })
}
```

Each test must cover Create, Read, Update, Delete (automatic), and ImportState.

### Documentation Generation

Docs in `docs/` are auto-generated from schema `MarkdownDescription` fields via `tfplugindocs` (triggered by `//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs` in `main.go`). Always run `make docs` and commit the result after schema changes.

## Keeping This File Up-to-Date

Update `AGENTS.md` whenever you make changes that affect how an agent should work in this repository:

- **New service or resource** — update the architecture section if the pattern changes, or note any exceptions
- **New utility functions** in `internal/tools/` — add them to the type conversion helpers table
- **New custom validators** in `internal/validators/` — document their purpose and usage
- **Changes to build, test, or lint commands** (e.g. new `Makefile` targets) — update the commands section
- **New conventions or patterns** established during a PR — capture them in Key Conventions
- **Schema versioning exceptions or new rules** — update the Schema Versioning section

When in doubt, err on the side of updating this file. A future agent reading stale instructions will make mistakes that are hard to catch.
