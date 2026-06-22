# go-ruby-marshal documentation

**A pure-Go (`CGO_ENABLED=0`) implementation of Ruby's Marshal** — the binary
serialization format produced by `Marshal.dump` and consumed by `Marshal.load`,
version **4.8**. Its byte output is **byte-for-byte identical to MRI Ruby's**,
differential-tested against the reference interpreter (Ruby 4.0.5).

`go-ruby-marshal/marshal` encodes and decodes the Marshal wire format against its
own small, typed value model rather than any interpreter's objects, so it is
**standalone and reusable**: any Go program can import it. The module path is
`github.com/go-ruby-marshal/marshal`.

!!! success "Byte-exact with MRI"
    The encoder is verified equal to MRI's `Marshal.dump` for every supported
    type: the compact Fixnum form and arbitrary-precision Bignum split at MRI's
    own boundary (`[-2³⁰, 2³⁰-1]`), the shortest round-tripping float formatted
    exactly as MRI does, the symbol table, the object-link table (so shared
    mutable objects and cycles round-trip), and string encodings. Decoding the
    bytes Ruby produces yields the same values back.

## The family

It is part of the `go-ruby-*` family of standalone front-end / runtime
components that [go-embedded-ruby](https://github.com/go-embedded-ruby) builds
on — alongside [go-ruby-parser](https://github.com/go-ruby-parser/parser)
(lexer/parser/AST) and [go-ruby-regexp](https://github.com/go-ruby-regexp/regexp)
(the Onigmo-style regexp engine). It has **no dependency on any interpreter**:
it works against its own value model and is bridged into one by converting
values at the boundary.

## Repositories

| Repo | What it is |
| --- | --- |
| [`marshal`](https://github.com/go-ruby-marshal/marshal) | the encoder/decoder — `dump.go`, `load.go`, the `Value` model in `value.go`, and the public API |
| [`docs`](https://github.com/go-ruby-marshal/docs) | this documentation site (MkDocs Material, versioned with mike) |
| [`brand`](https://github.com/go-ruby-marshal/brand) | logo and brand assets |

## What it is

- A faithful encoder/decoder for Ruby's **Marshal 4.8** wire format — the same
  bytes MRI writes and reads.
- **Pure Go, `CGO_ENABLED=0`** — trivial cross-compilation, a single static
  binary, no C toolchain.
- A **standalone module** with a small, Go-idiomatic public surface
  (`Dump` / `Load` over a typed `Value`).
- **Differential-tested against MRI** for byte equality, not approximated.

## What it is not

- **Not** a general Ruby object marshaller: it covers the value subset listed
  below (nil/bool/Integer/Float/Symbol/String/Array/Hash). User-defined classes,
  `_dump`/`_load`, and `marshal_dump`/`marshal_load` hooks are out of scope of
  this core module.
- **Not** dependent on go-embedded-ruby — the dependency runs the other way.

## Install

```sh
go get github.com/go-ruby-marshal/marshal
```

## Quick start

```go
import "github.com/go-ruby-marshal/marshal"

// Encode  →  Ruby Marshal bytes (Marshal.dump compatible)
b := marshal.Dump(&marshal.Array{Elems: []marshal.Value{
    marshal.NewInt(1), marshal.NewString("hi"), marshal.Bool(true),
}})

// Decode  ←  bytes produced by Ruby's Marshal.dump
v, err := marshal.Load(b)
if err != nil {
    // err is (or wraps) marshal.ErrShort for a truncated stream
}
```

## Where to go next

- [API](api.md) — the `Dump` / `Load` functions, the `ErrShort` sentinel, and
  the full `Value` model with Go examples.
- [The Marshal format](format.md) — what the 4.8 wire format encodes and the
  parts this module implements (Fixnum/Bignum split, shortest float, symbol
  table, object-link table, string encodings, Hash default).
- [Contributing](contributing.md) — the hard rules (pure Go, 100% coverage,
  differential testing against MRI).

Source lives at
[github.com/go-ruby-marshal/marshal](https://github.com/go-ruby-marshal/marshal).
