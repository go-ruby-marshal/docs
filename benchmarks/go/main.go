// SPDX-License-Identifier: BSD-3-Clause
package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"os"

	"github.com/go-ruby-marshal/marshal"
)

// buildGraph constructs the representative object graph used by both the Go and
// the Ruby side of the benchmark. It mixes every core Marshal type — nested
// Hash / Array / String / Integer (Fixnum and Bignum) / Float / Symbol — and
// deliberately shares two mutable objects (a String and an Array) so the
// object-link table and the symbol table are exercised: `shared` and `arr`
// each appear twice, and the symbol :alpha appears twice. The Ruby program in
// ruby/marshal.rb builds the byte-for-byte identical structure.
func buildGraph() marshal.Value {
	shared := marshal.NewString("shared-payload")
	arr := &marshal.Array{Elems: []marshal.Value{
		marshal.NewInt(1), marshal.NewInt(2), marshal.NewInt(3),
	}}
	big70 := new(big.Int).Lsh(big.NewInt(1), 70) // 2**70, forces the Bignum form

	nested := &marshal.Hash{
		Keys: []marshal.Value{marshal.Symbol("a"), marshal.Symbol("b"), marshal.Symbol("c")},
		Vals: []marshal.Value{
			&marshal.Array{Elems: []marshal.Value{
				marshal.NewInt(1), marshal.NewInt(2), marshal.NewInt(3),
				marshal.NewString("hello"), marshal.Float(4.5),
			}},
			shared, // first reference to the shared string
			arr,    // first reference to the shared array
		},
	}

	return &marshal.Hash{
		Keys: []marshal.Value{
			marshal.Symbol("name"), marshal.Symbol("version"), marshal.Symbol("pi"),
			marshal.Symbol("tags"), marshal.Symbol("nested"), marshal.Symbol("dup_str"),
			marshal.Symbol("dup_arr"), marshal.Symbol("flags"), marshal.Symbol("big"),
		},
		Vals: []marshal.Value{
			marshal.NewString("go-ruby-marshal"),
			marshal.NewInt(48),
			marshal.Float(3.14159),
			&marshal.Array{Elems: []marshal.Value{
				marshal.Symbol("alpha"), marshal.Symbol("beta"), marshal.Symbol("alpha"),
			}},
			nested,
			shared, // second reference -> object link
			arr,    // second reference -> object link
			&marshal.Array{Elems: []marshal.Value{
				marshal.Bool(true), marshal.Bool(false), marshal.Nil{},
				marshal.NewInt(100), marshal.NewInt(-42), marshal.NewInt(1000000),
			}},
			marshal.Int{I: big70},
		},
	}
}

func main() {
	graph := buildGraph()
	data := marshal.Dump(graph)

	// Emit the marshaled-byte digest (ignored by run.sh, which filters RESULT
	// lines) so the runner can prove Go and MRI produce identical bytes.
	sum := sha256.Sum256(data)
	fmt.Printf("DUMP\t%x\t%d\n", sum, len(data))
	if os.Getenv("VERIFY") != "" {
		return
	}

	bench("dump-graph", 2000, func() { sink = marshal.Dump(graph) })
	bench("load-graph", 2000, func() { v, _ := marshal.Load(data); sink = v })
}
