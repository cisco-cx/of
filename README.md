# of
Observability Framework

## Commands

### `of help`

Status: **TODO**

### `of version`

Status: **TODO**

### `of handler apic`

Status: **WIP**

This subcommand starts a daemon that scrapes APIC servers in Cisco ACI clusters for their fault lists and then notifies Prometheus Alertmanager to fire and resolve Alertmanager alerts.

## Repo Layout

### Overview

Here are some general rules that apply regardless of the directory or Go package:

- SHOULD include unit tests.
  - Regardless of the directory, please include per-`.go` file test files as practical. For example, if you have a file called `foo.go`, then you SHOULD probably also have a `foo_test.go` to hold matching unit tests.
- MUST use go modules for all dependencies.
- MUST support the building of ALL `cmd` executables in the shared `/Dockerfile`. Each command MUST compile to a static binary and not require Go on the build host. See: https://hub.docker.com/_/golang
- Where imported, named-version packages SHOULD be import-aliased back to their dependency name (e.g. `postgres`) -- and they may be combined as needed to form executable commands (e.g. `of-handler-apic` vs. `of-handler-snmp`).
- If two named versions of a single dependency's implementation `pkg` must be imported, you MUST alias each named-version like this: `postgresv1alpha1` (for `package v1alpha1`) and `postgresv1` (for `package v1`). Try to avoid that situation by forking widely-used named-version code to a new named-version package.

### Directories / Go Packages

#### `/`

This directory contains files like `LICENSE`, `NOTICE`, `README.md`, `go.mod`, `go.sum`.

Files with names ending in `.go` MUST not be committed in `/`.

#### `/lib/$named_version`

(e.g. `package v1` in a directory like `/lib/v1`)

Each package in this directory pattern SHOULD:
- Contain a version-named set (e.g. `v1alpha1`) of [domain types and interfaces](https://www.youtube.com/watch?v=LMSbsW1Xpwg) for the Observability Framework.
- Have NO external dependencies, with the only exception being unavoidable dependencies on the Go standard library.

Each package matching this pattern MUST NOT:
- Import any code from directories: `/cmd`, `/wrap`. This rule is necessary to keep our domain types and interfaces decoupled from their implementations and avoid circular package dependencies as these are not supported by Go.

Where **no external dependency exists**, packages in `/lib` MAY:
- Implement its own interfaces. Here's a somewhat contrived example of a compliant scenario:

```
package v999alpha1  // github.com/cisco-cx/of/lib/v999alpha1

import "fmt"

// Domain Types

// APIVersion represents a named version of an API (e.g. "v1alpha1", "v1").
type APIVersion string

// Implementations (with very trivial standard-library dependencies)

// String implements the APIVersionStringer interface.
func (v APIVersion) String() string {
    return fmt.Sprintf("%v", v)
}
```

#### `/wrap/$dependency_name`

Each directory in `/wrap` MUST contain Go code in named-version packages that wrap **one** external or standard library dependency (see below).

#### `/wrap/$dependency_name/$named_version`

(e.g. `package v2alpha1` package in a directory like `/wrap/postgres/v2alpha1`)

Each package matching this pattern SHOULD:
- **Implement not define** domain types and interfaces as imported from `/lib/$named_version` for no more than [one external dependency](https://www.youtube.com/watch?v=LMSbsW1Xpwg).
- Be named for its dependency and version.
  - For example, if the package we depend on is called `postgres`, the first-draft package for that would be `package v1alpha1` inside `/wrap/postgres/v1alpha1`.
- Be the definitive place we implement that external dependency. In this way, over in `cmd` we SHOULD only import our own wrappers of external dependencies. Exceptions (even on standard library's `http`) should be avoided if at all possible.

Each package matching this pattern MAY:
- Within reason, import any other package under the `/wrap` directory.

#### `/cmd/$EXECUTABLE`

(e.g. `package main` package for executable `foo` would be in `/cmd/foo`)

The `/cmd` directory's packages contains the Go packages for executable commands (e.g. `/cmd/foo`) that wire together named-version packages in `/wrap` with those in `/lib`.

In each of these "command" packages, we combine one or more packages from `/lib` (contains domain types and interfaces) with one or more "one external dependency" packages like `/wrap/$dependency_name/$named_version` (contains **implementations** of domain types and interfaces) to arrive at statically compiled binaries that MAY be deployed inside one or more Docker images.

Each package in this directory pattern MAY:
- Import a couple of external dependencies to simplify the building of a CLI. For example, you might choose to directly import [cobra](https://github.com/spf13/cobra) or [kingpin](https://github.com/alecthomas/kingpin) in your `/cmd` package.

Each package in this directory pattern SHOULD:
- NOT import any non-standard-lib external depedencies not related to the "MAY" list directly above this one.
- NOT try to avoid or skip wrapping your external dependency over in `/wrap`.

#### `/cmd/of`

The `of` command is to become the core Go command of the Observability Framework. That is, we plan to ship one combined static binary that can assume mutliple personalities (e.g. am-apic-client AND am-snmp-client), not unlike the [hashicorp/vault](https://github.com/hashicorp/vault)'s `vault` command has [subcommands](https://www.vaultproject.io/docs/commands/) `server` and `agent`.

This command SHOULD eventually support multiple named-version `/lib` and `/wrap` packages, but for now a tightly coupling to a single named-version for any `/lib` or `/wrap` packages is allowed.

#### `/mock`

This directory will contain mocks for all of the above. The form that this directory takes is flexible right now -- and so its design is left open until further notice.

## Inspiration

### for Project Layout

- https://www.youtube.com/watch?v=LMSbsW1Xpwg
- https://www.youtube.com/watch?v=MzTcsI6tn-0
- https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1
- https://medium.com/wtf-dial/wtf-dial-domain-model-9655cd523182
- https://github.com/benbjohnson/wtf
- https://medium.com/p/7cdbc8391fc1/responses/show
- https://www.youtube.com/watch?v=zzAdEt3xZ1M
- https://github.com/kubernetes/kubernetes
- https://github.com/benbjohnson/peapod

### for Error Handling

- https://dave.cheney.net/2016/04/07/constant-errors
- https://blog.golang.org/error-handling-and-go
