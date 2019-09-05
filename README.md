# of
Observability Framework

## Commands

### `cmd/of-handler-apic`

Status: **UNDER CONSTRUCTION**

The [of-handler-apic](https://github.com/cisco-cx/of/blob/master/cmd/of-handler-apic) command scrapes APIC servers in Cisco ACI clusters for their fault lists and then notifies Prometheus Alertmanager to fire and resolve Alertmanager alerts.

## Project Layout

### Overview

Here are some general rules that apply regardless of the directory or Go package:

- SHOULD include unit tests.
  - Regardless of the directory, please include per-`.go` file test files as practical. For example, if you have a file called `foo.go`, then you SHOULD probably also have a `foo_test.go` to hold matching unit tests.
- MUST use go modules for all dependencies.
- MUST support the building of ALL `cmd` executables in the shared `/Dockerfile`. Each command MUST compile to a static binary and not require Go on the build host. See: https://hub.docker.com/_/golang

### Directories / Go Packages

#### `/`

(`package of`)

This root package contains [domain types and public interfaces](https://www.youtube.com/watch?v=LMSbsW1Xpwg) for the Observability Framework. The code here should have no external dependencies, with the only exception being unavoidable dependencies on the Go standard lib.

**Until at least of 1.x.x, you MUST consider `package of` to be unstable.** If it happens that multiple named versions of a domain type or public interface need to exist, it MAY be good to start an `/external/$NAMED_VERSION` package to follow suit with the design below. However, we won't do this unless it becomes absolutely and obviously the right choice.

#### `/internal/$DEPENDENCY/$NAMED_VERISON`

(e.g. various `package v1alpha1` packages in directories like `/internal/postgres/v1alpha1`)

The `/internal` directory contains [internal packages](https://golang.org/doc/go1.4#internalpackages).

Each internal package SHOULD:
- Implement domain types and domain interfaces as imported from `package of` for no more than [one external dependency](https://www.youtube.com/watch?v=LMSbsW1Xpwg).
- Be named for its dependency and version.
  - For example, if the package we depend on is called `postgres`, the first internal package for that should be `package v1alpha1` in `/internal/postgres/v1alpha1`.

Over in `/cmd` packages, any versioned `/internal` packages SHOULD be import-aliased back to their dependency name (e.g. `postgres`) -- and they may be combined as needed to form executable commands (e.g. `of-handler-apic` vs. `of-handler-snmp`).

#### `/cmd/$EXECUTABLE`

(various `package main` packages in directories like `/cmd/of-handler-apic`)

This directory's packages contain source code and Dockerfile content for executable commands (e.g. `/cmd/of-handler-apic`). In these "command" packages, we wire together code from `/` and `/internal` to arrive at statically compiled binaries that will be deployed inside Docker images.

#### `/mock`

This directory will contain mocks for all of the above. The form that it takes is flexible -- and this is left open until we figure out what works.

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
