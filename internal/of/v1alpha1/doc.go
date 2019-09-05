// Package `of/internal/of/v1alpha1` is the first implementation of `package of`.
//
// This package is rather odd in that this is an internal implementation of the
// public `of` package's
// [domain types and interfaces](https://www.youtube.com/watch?v=LMSbsW1Xpwg).
//
// That is, in this package, we'll often import structs from `package of` and
// implement their interfaces. We're doing that under `of/internal/of/v1alpha1`
// so that if we need major changes in the **implementation details**, we have
// the ability to fork off to `v1alpha2` or whatever as necessary.
//
// Eventually, and if needs presents, we MAY promote `of/internal/of` packages
// to `of/external`. Until then though, the only consumer of this code SHOULD
// be packages in `of/cmd`.
package v1alpha1
