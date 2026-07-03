<!-- SPDX-License-Identifier: BSD-3-Clause -->
# `go-ruby-marshal` library-level benchmark harness

Reproducible, cross-runtime benchmark of the **pure-Go `go-ruby-marshal` library**
against the reference Ruby runtimes (MRI, MRI + YJIT, JRuby, TruffleRuby). It
measures the **library primitive** through its Go API, isolated from the rbgo
interpreter, so the numbers answer: *is the pure-Go implementation as fast as the
reference runtime's own `Marshal`?*

## Layout

- `go/`             — self-contained Go driver; `go.mod` pins the published library.
- `ruby/marshal.rb` — the equivalent workload; `ruby/_harness.rb` is the shared timer.
- `run.sh`          — runs every available runtime and prints one Markdown table per
  sub-benchmark (ns/op + ratio vs MRI).

## Run

```sh
bash benchmarks/run.sh
```

Environment knobs: `OUTER` (timed passes, default 25), `WARM` (untimed warm-up
passes, default 3), and `RUBY`/`JRUBY`/`TRUFFLERUBY` to select runtime binaries.

## Method

Each process runs `WARM` untimed passes (to let the JVM/GraalVM JITs warm up),
then `OUTER` timed passes of a fixed inner loop, timed with a monotonic clock;
the **best** pass is reported as **ns/op**. Interpreter start-up is outside the
timed region.

The Go driver and the Ruby script build a **byte-for-byte identical** object
graph — nested Hash/Array/String/Integer (Fixnum and Bignum)/Float/Symbol, with a
shared String and a shared Array (each referenced twice) plus a repeated Symbol,
so Marshal's object-link and symbol tables are exercised the same way on both
sides. Because Marshal is a byte-exact format, correctness is checked directly:
both sides print the SHA-256 of their `Marshal.dump` output and the runner
asserts they are **identical to MRI** before any timing (current graph:
`87a5f8f8…`, 214 bytes). Results are published, dated, in `../docs/performance.md`.
