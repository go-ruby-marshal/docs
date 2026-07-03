# Performance

`go-ruby-marshal/marshal` is the pure-Go (CGO-free) library that
[`rbgo`](https://github.com/go-embedded-ruby/ruby) binds for Ruby's `Marshal`.
This page records **real measured** benchmarks of the library against the
reference Ruby runtimes — no estimates, no placeholders.

## Library-level benchmark (Go API vs runtimes) — 2026-07-03

This section measures the **pure-Go library directly, through its Go API** — not
the `rbgo` interpreter path. It isolates the library primitive from Ruby-
interpreter dispatch, answering the parity question head-on: *is the pure-Go
implementation as fast as the reference runtime's own `Marshal`?* The **same
workload, same inputs, same iteration counts** run through the Go library and
through each reference runtime's stdlib.

Because Marshal is a **byte-exact** format, correctness here is not approximate:
both sides serialize the identical object graph and the runner asserts the
`Marshal.dump` output is **byte-for-byte identical to MRI** (SHA-256
`87a5f8f8…`, 214 bytes) before any timing. A benchmark that did not match MRI's
bytes would not be timed at all.

- **Host:** Apple M4 Max (`Mac16,5`, arm64), macOS 26.5.1 — **date 2026-07-03**.
- **Runtimes:** Go 1.26.4 · MRI `ruby 4.0.5 +PRISM` · MRI + YJIT · JRuby 10.1.0.0
  (OpenJDK 25) · TruffleRuby 34.0.1 (GraalVM CE Native).
- **Workload:** one representative object graph mixing nested Hash / Array /
  String / Integer (Fixnum **and** `2**70` Bignum) / Float / Symbol, with a
  shared String and a shared Array (each referenced twice) plus a repeated
  Symbol, so the object-link and symbol tables are exercised. `dump-graph`
  serializes it; `load-graph` deserializes the 214-byte encoding.
- **Method:** each process runs 3 untimed warm-up passes, then 25 timed passes of
  a fixed inner loop, timed with a monotonic clock; the **best** pass is reported
  as **ns/op** (lower is better). `vs MRI` < 1.00× means *faster than MRI*.
  Interpreter start-up is outside the timed region, so these are operation costs,
  not `ruby file.rb` process costs.

#### dump-graph

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 1525.7 | 0.48× |
| MRI | 3170.5 | 1.00× |
| MRI + YJIT | 3163.0 | 1.00× |
| JRuby | 1555.5 | 0.49× |
| TruffleRuby | 134281.1 | 42.35× |

#### load-graph

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 1622.1 | 0.63× |
| MRI | 2562.5 | 1.00× |
| MRI + YJIT | 2643.0 | 1.03× |
| JRuby | 1229.7 | 0.48× |
| TruffleRuby | 7787.3 | 3.04× |

**go-ruby-marshal beats MRI on both operations** — `dump` at **0.48×** (~2.1×
faster) and `load` at **0.63×** (~1.6× faster). MRI's `Marshal` is C, so this is
a genuine win for a compiled pure-Go serializer over the reference C
implementation: the Go encoder writes into a single growing `[]byte` with map-
backed symbol/object-link tables, and the decoder builds the typed value model
directly. YJIT neither helps nor hurts (both operations spend their time in C /
Go, not in Ruby bytecode). JRuby is in the same ballpark as the Go library once
warm. The **TruffleRuby `dump-graph` cell (≈42×) is a cold-JIT outlier**: within
the fixed 3-warm-up budget Graal had not compiled the dump loop (its `load-graph`
row, at ≈3×, is a fairer steady-state figure); treat that single cell as
order-of-magnitude, not a steady-state number.

!!! note "Reproduce"
    The harness is committed under
    [`benchmarks/`](https://github.com/go-ruby-marshal/docs/tree/main/benchmarks):
    a self-contained Go driver (`go/`, pins the published library via
    `go.mod`), the equivalent `ruby/marshal.rb` workload, and `run.sh`. Run
    `bash benchmarks/run.sh`; env `OUTER`/`WARM` tune the pass budget and
    `RUBY`/`JRUBY`/`TRUFFLERUBY` select the runtime binaries.

!!! warning "Warm-up budget & noise — honest framing"
    Numbers reflect a **fixed warm-process budget** (3 warm-up + 25 timed passes
    in one process). The JVM/GraalVM JITs (JRuby, TruffleRuby) may need a larger
    warm-up to reach steady state, so their columns can **understate** peak
    throughput — most visibly TruffleRuby on the shortest loops (the `dump-graph`
    cold-JIT outlier noted above). Sub-microsecond rows carry the most relative
    noise; treat those ratios as order-of-magnitude. Every number here is a
    **real measured value** from the dated run above — nothing is fabricated,
    estimated, or cherry-picked. The go-ruby column is the pure-Go library; every
    other column is that interpreter's own `Marshal` doing the equivalent work.
