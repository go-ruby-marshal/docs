# Contributing

Contributions are welcome. `go-ruby-marshal/marshal` is built to a small set of
non-negotiable rules — they are what keep the module pure-Go, correct, and
byte-compatible with MRI. Please read these before opening a pull request.

## Hard rules

- **Build from source — no vendoring.** Everything compiles from source. Do not
  reach for prebuilt binaries or vendored blobs as a shortcut; being able to
  compile from source is a guarantee of independence.
- **100% test coverage target, enforced in CI.** New code ships with tests, and
  coverage is a CI gate. Fill the error branches (truncated input, unknown tags,
  unsupported versions), not just the happy path.
- **All GitHub content in English.** Issues, pull requests, commits, comments,
  and discussions are English-only.
- **Differential testing against MRI.** Correctness is defined by reference
  Ruby. Values are run through both Ruby's `Marshal.dump` and this encoder and
  the bytes are compared **exactly** — and bytes Ruby produced are decoded back
  and compared. "Byte-for-byte identical to MRI" is a tested property, not an
  approximation from memory.
- **Pure Go, cgo disabled.** The whole point is a single static binary with no C
  toolchain. Code must build with `CGO_ENABLED=0`. If a feature seems to need C,
  it needs a pure-Go path instead.

## Workflow

1. Pick or open an issue describing the change.
2. Work test-first: add the differential / unit tests, then make them pass.
3. Run the full suite with coverage and confirm the gate is green.
4. Open a PR in English, referencing the issue.

## Where things live

The encoder/decoder and the `Value` model are in
[`github.com/go-ruby-marshal/marshal`](https://github.com/go-ruby-marshal/marshal)
(`dump.go`, `load.go`, `value.go`). This documentation site is in
[`github.com/go-ruby-marshal/docs`](https://github.com/go-ruby-marshal/docs).
Start from the [API](api.md) and [The Marshal format](format.md) to find the
right place for your change.
