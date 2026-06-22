# The Marshal format

Ruby's **Marshal** is a binary serialization format. `Marshal.dump` writes a
value as bytes; `Marshal.load` reads them back. `go-ruby-marshal/marshal`
implements **version 4.8** of that format and matches MRI's output byte-for-byte
for the supported value subset. This page describes the parts of the format the
module implements; for the Go surface that exposes them, see the [API](api.md).

!!! info "Version header"
    Every stream begins with the two-byte version `0x04 0x08` (major 4, minor 8).
    `Load` rejects a stream whose version it does not support.

## Integers: Fixnum and Bignum

Marshal splits integers into two encodings:

- **Fixnum** — the compact `i`-tagged form, used for values in
  **`[-2³⁰, 2³⁰-1]`**. This boundary is MRI's marshal Fixnum range, **not** the
  platform word size, so the split is the same on 32- and 64-bit hosts.
- **Bignum** — the `l`-tagged arbitrary-precision form, used for every value
  outside that range, stored as sign plus little-endian magnitude.

`Int` holds a `*big.Int`, so any magnitude round-trips; `Dump` picks the form
automatically.

## Floats: shortest round-trip

Floats are marshalled as a decimal **string** — the shortest decimal that
round-trips back to the exact `float64`, formatted **exactly as MRI does**
(including its exponent-notation choices and the special values `±inf`, `nan`,
and `-0`). This is what keeps `Dump(Float(x))` byte-equal to Ruby for every
double.

## Symbols and the symbol table

A `Symbol` is written by name (without the leading colon). The format keeps a
**symbol table**: the first occurrence of a symbol is written in full, and every
later occurrence of the same symbol is written as a compact reference to its
table index. The module implements both sides, so symbol-heavy structures encode
as compactly as MRI's and decode back to the same symbols.

## The object-link table: shared objects and cycles

Mutable objects (`String`, `Array`, `Hash`) participate in an **object-link
table**. The first time a given object is encountered it is written out in full
and assigned an id; any later reference to the **same object** is written as a
link (`@`) to that id.

This is why the composite `Value` types are pointers: identity is what the
format records. Two consequences follow, both matching MRI exactly:

- **Shared objects** are encoded once. If the same `*Str` appears twice in an
  array, the bytes contain the string once and a link the second time.
- **Cycles** are representable. `a = []; a << a` — a self-referential array —
  round-trips, because the container is registered before its contents are
  written (and reserved before its contents are read).

## Strings and encodings

A `String` is raw bytes plus an encoding, marshalled with an instance-variable
wrapper that carries the encoding marker:

| Encoding | How it is marshalled |
| --- | --- |
| **UTF-8** (default) | instance variable `E => true` |
| **US-ASCII** | `E => false` |
| **ASCII-8BIT** (BINARY) | a bare `String` with no encoding ivar |
| **Named** (any other) | `encoding => <name>` |

The `Str.Enc` field selects which of these is emitted, and `Str.Name` supplies
the name for the `Named` case.

## Hashes, including a default

A `Hash` preserves **insertion order** (parallel `Keys` / `Vals` slices, as Ruby
hashes do). A hash created with a default value (`Hash.new(d)`) is marshalled
with the dedicated `TYPE_HASH_DEF` tag, and the `Default` field on `Hash`
carries that value through a round-trip.

## Supported type summary

| Ruby type | `Value` | Notes |
| --- | --- | --- |
| `nil` | `Nil` | |
| `true` / `false` | `Bool` | |
| `Integer` | `Int` | Fixnum in `[-2³⁰, 2³⁰-1]`, else Bignum |
| `Float` | `Float` | shortest round-tripping decimal |
| `Symbol` | `Symbol` | symbol table (repeated symbols link) |
| `String` | `*Str` | UTF-8 / US-ASCII / ASCII-8BIT / named |
| `Array` | `*Array` | object links |
| `Hash` | `*Hash` | insertion order; optional default |

Everything above is verified equal to MRI Ruby's `Marshal.dump` through
differential tests against the reference interpreter (Ruby 4.0.5).
