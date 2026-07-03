# frozen_string_literal: true
# SPDX-License-Identifier: BSD-3-Clause
require "digest"
require_relative "_harness"

# Build the representative object graph. It must be byte-for-byte identical to
# the graph built by ../go/main.go: nested Hash/Array/String/Integer(Fixnum and
# Bignum)/Float/Symbol, with a shared String and a shared Array (each referenced
# twice) plus a repeated Symbol, so Marshal's object-link and symbol tables are
# exercised the same way on both sides.
def build_graph
  shared = "shared-payload"
  arr    = [1, 2, 3]
  nested = { a: [1, 2, 3, "hello", 4.5], b: shared, c: arr }
  {
    name:    "go-ruby-marshal",
    version: 48,
    pi:      3.14159,
    tags:    [:alpha, :beta, :alpha],
    nested:  nested,
    dup_str: shared, # second reference -> object link
    dup_arr: arr,    # second reference -> object link
    flags:   [true, false, nil, 100, -42, 1_000_000],
    big:     2**70,
  }
end

graph = build_graph
data  = Marshal.dump(graph)

# Emit the marshaled-byte digest (ignored by run.sh) so the runner can prove
# MRI and the Go library produce identical bytes.
printf("DUMP\t%s\t%d\n", Digest::SHA256.hexdigest(data), data.bytesize)
exit if ENV["VERIFY"]

bench("dump-graph", 2000) { Marshal.dump(graph) }
bench("load-graph", 2000) { Marshal.load(data) }
