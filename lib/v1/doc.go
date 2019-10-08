/*
Package `v1` (`of/lib/v1`) is the first named-version of the
Observability Framework library (aka. "OF lib").

Status: **WIP**

This package defines and exposes all the public
[domain types and interfaces](https://www.youtube.com/watch?v=LMSbsW1Xpwg)
for OF lib `v1`.

In this package we MUST NOT:
- Import any code from directories: `/cmd`, `/wrap`.
  - This rule is necessary to keep our domain types and interfaces decoupled
    from their implementations and avoid circular package dependencies as
    these are not supported by Go.

Where there is **no external dependency**, this package MAY:
- Implement its own interfaces.
*/

package v1
