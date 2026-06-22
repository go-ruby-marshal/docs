# API

The public API lives at the module root, `github.com/go-ruby-marshal/marshal`.
It is intentionally small: two functions (`Dump` / `Load`), one error sentinel
(`ErrShort`), and a typed `Value` model that the format is encoded from and
decoded into.

## Functions

### `Dump(v Value) []byte`

Returns the Ruby Marshal (version 4.8) encoding of `v` — the same bytes MRI's
`Marshal.dump` produces. The output begins with the two-byte version header
`0x04 0x08`.

```go
b := marshal.Dump(marshal.NewInt(42))
// b is Marshal.dump(42) byte-for-byte
```

### `Load(b []byte) (Value, error)`

Decodes a Ruby Marshal (version 4.8) byte stream into a `Value`. It returns an
error for a **truncated** stream, an **unsupported version**, or an **unknown
type tag**.

```go
v, err := marshal.Load(b)
if errors.Is(err, marshal.ErrShort) {
    // the input ended in the middle of a value
}
```

### `var ErrShort error`

`marshal: unexpected end of input` — reported when the input ends in the middle
of a value. The error returned by `Load` for a truncated stream matches it with
`errors.Is`, so callers can distinguish truncation from a malformed-but-complete
stream.

## The `Value` model

```go
type Value interface{ RubyClass() string }
```

`Value` is a Ruby value in the subset this package serializes. The concrete
types are `Nil`, `Bool`, `Int`, `Float`, `Symbol`, `*Str`, `*Array`, and
`*Hash`. Every `Value` reports the name of its corresponding Ruby class via
`RubyClass()`.

The composite, mutable types (`*Str`, `*Array`, `*Hash`) are **pointers** so
that identity is observable: the same pointer appearing more than once in a
structure is encoded once and thereafter as an **object link**, exactly as MRI
does — which is also what makes cyclic structures representable.

| Type | Go declaration | `RubyClass()` |
| --- | --- | --- |
| `Nil` | `type Nil struct{}` | `NilClass` |
| `Bool` | `type Bool bool` | `TrueClass` / `FalseClass` |
| `Int` | `type Int struct{ I *big.Int }` | `Integer` |
| `Float` | `type Float float64` | `Float` |
| `Symbol` | `type Symbol string` | `Symbol` |
| `*Str` | `type Str struct{ Bytes []byte; Enc Encoding; Name string }` | `String` |
| `*Array` | `type Array struct{ Elems []Value }` | `Array` |
| `*Hash` | `type Hash struct{ Keys, Vals []Value; Default Value }` | `Hash` |

### Scalars

- **`Nil`** — Ruby `nil`.
- **`Bool`** — Ruby `true` / `false`.
- **`Int`** — a Ruby `Integer` of any magnitude, held as a `*big.Int`. `Dump`
  emits the compact **Fixnum** form for values in `[-2³⁰, 2³⁰-1]` (MRI's marshal
  Fixnum range, not the platform word size) and the **Bignum** form otherwise.
- **`Float`** — a Ruby `Float`, formatted as the shortest decimal that
  round-trips, exactly as MRI does.
- **`Symbol`** — a Ruby `Symbol` (its name, without the leading colon). Repeated
  symbols become links via the format's symbol table.

### Strings and encodings

```go
type Str struct {
    Bytes []byte
    Enc   Encoding
    Name  string // only used when Enc == Named
}
```

`Str` is a Ruby `String`: raw bytes plus an encoding. The zero value is an empty
UTF-8 string. `Encoding` selects how the encoding is marshalled:

| Constant | Meaning | Marshalled as |
| --- | --- | --- |
| `UTF8` | the default | instance variable `E => true` |
| `USASCII` | US-ASCII | `E => false` |
| `ASCII8BIT` | BINARY | a bare `String` with no encoding ivar |
| `Named` | any other encoding | `encoding => <Name>` (uses `Str.Name`) |

### Collections

```go
type Array struct{ Elems []Value }

type Hash struct {
    Keys    []Value
    Vals    []Value
    Default Value
}
```

- **`Array`** — a Ruby `Array`.
- **`Hash`** — a Ruby `Hash`. `Keys` and `Vals` are parallel slices preserving
  insertion order, as Ruby hashes do. `Default`, when non-nil, is the hash's
  default value (`Hash.new(default)`); it is marshalled with the `TYPE_HASH_DEF`
  tag.

## Constructors

- **`NewInt(n int64) Int`** — returns an `Int` holding `n`.
- **`NewString(s string) *Str`** — returns a UTF-8 `*Str` holding `s`.

For magnitudes beyond `int64`, build `Int{I: someBigInt}` directly; for other
encodings, set `Str.Enc` (and `Str.Name` for `Named`) yourself.

## Examples

### Round-trip a nested structure

```go
in := &marshal.Hash{
    Keys: []marshal.Value{marshal.Symbol("name"), marshal.Symbol("tags")},
    Vals: []marshal.Value{
        marshal.NewString("go-ruby-marshal"),
        &marshal.Array{Elems: []marshal.Value{
            marshal.NewString("ruby"), marshal.NewString("marshal"),
        }},
    },
}

b := marshal.Dump(in)        // == Ruby's Marshal.dump of the same hash
out, err := marshal.Load(b)  // out is a *marshal.Hash equal to in
_ = err
```

### Shared objects and cycles

Because the composite types are pointers, identity is preserved through the
object-link table:

```go
shared := marshal.NewString("x")
a := &marshal.Array{Elems: []marshal.Value{shared, shared}}
b := marshal.Dump(a) // the string is written once, then linked — exactly as MRI

// A self-referential array (a = []; a << a):
cyc := &marshal.Array{}
cyc.Elems = append(cyc.Elems, cyc)
_ = marshal.Dump(cyc) // round-trips the cycle
```

### A big integer

```go
n, _ := new(big.Int).SetString("123456789012345678901234567890", 10)
b := marshal.Dump(marshal.Int{I: n}) // Bignum form (outside the Fixnum range)
```

## Relationship to Ruby

The value model maps onto Ruby's core types, but the surface follows Go
conventions: an explicit `error`, byte slices, value/pointer types, and
`*big.Int` for arbitrary precision. An embedded interpreter such as
[go-embedded-ruby](https://github.com/go-embedded-ruby) bridges its own objects
to and from this model at the boundary; the encoding/decoding work — and the
byte-exact compatibility with MRI — lives here.
