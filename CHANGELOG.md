# Changelog

## Unreleased

### Breaking Changes

* (generator) [#25](https://github.com/yihuang/go-abi/pull/25) Removed `-uint256` flag and `generator.UseUint256(...)`, uint256 backing is now selected by `-tags uint256` at compile time, and flag from `//go:generate` directives is dropped.

### Bug Fixes

### Improvements

* (generator) [#17](https://github.com/yihuang/go-abi/pull/17) Added `-lazy` code generation mode for on-demand decoding views (`*View`, `*SliceView`, `*ArrayView`) with `Get`, `Len`, `Raw`, and `Materialize` helpers.
