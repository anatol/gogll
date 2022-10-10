package sppf

import (
	"bytes"
	"text/template"

	"github.com/goccmack/goutil/ioutil"
)

func Gen(sppfFile string, pkg string) {
	tmpl, err := template.New("sppf").Parse(tmpl)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	if err = tmpl.Execute(buf, pkg); err != nil {
		panic(err)
	}
	if err = ioutil.WriteFile(sppfFile, buf.Bytes()); err != nil {
		panic(err)
	}
}

const tmpl = `// Package sppf is generated by gogll. Do not edit.

/*
Package sppf implements a Shared Packed Parse Forest as defined in:

	Elizabeth Scott, Adrian Johnstone
	GLL parse-tree generation
	Science of Computer Programming (2012), doi:10.1016/j.scico.2012.03.005
*/
package sppf

import (
	"fmt"
	"bytes"
	"github.com/goccmack/goutil/ioutil"

	"{{.}}/parser/symbols"
)

type Node interface {
	isNode()
	dot(*dotBuilder)
	Label() string
	String() string
}

type IntermediateNode struct {
	NT         symbols.NT
	Body       symbols.Symbols
	Pos        int
	Lext, Rext int
	Children   []*PackedNode
}

type SymbolNode struct {
	Symbol     symbols.Symbol
	Lext, Rext int
	Children   []*PackedNode
}

type PackedNode struct {
	NT                symbols.NT
	Body              symbols.Symbols
	Pos               int
	Lext, Pivot, Rext int

	LeftChild  Node // Either an intermediate or Symbol node
	RightChild *SymbolNode
}

func (*IntermediateNode) isNode() {}
func (*SymbolNode) isNode()       {}
func (*PackedNode) isNode()       {}

func slotString(nt symbols.NT, body symbols.Symbols, pos int) string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "%s:", nt)
	for i, sym := range body {
		fmt.Fprint(w, " ")
		if i == pos {
			fmt.Fprint(w, "•")
		}
		fmt.Fprint(w, sym)
	}
	if len(body) == pos {
		fmt.Fprint(w, "•")
	}
	return w.String()
}

func (n *IntermediateNode) Label() string {
	return fmt.Sprintf("\"%s:,%d,%d\"", slotString(n.NT, n.Body, n.Pos), n.Lext, n.Rext)
}

func (n *SymbolNode) Label() string {
	return fmt.Sprintf("\"%s,%d,%d\"", n.Symbol, n.Lext, n.Rext)
}

func (n *PackedNode) Label() string {
	return fmt.Sprintf("\"%s,%d\"", slotString(n.NT, n.Body, n.Pos), n.Pivot)
}

func (n *IntermediateNode) String() string {
	return "IN: " + n.Label()
}

func (n *SymbolNode) String() string {
	return "SN: " + n.Label()
}

func (n *PackedNode) String() string {
	return "PN: " + n.Label()
}

//---- Dot ----

type dotBuilder struct {
	iNodes map[*IntermediateNode]bool
	pNodes map[string]bool // index PN Label
	sNodes map[*SymbolNode]bool
	w      *bytes.Buffer
}

func (root *SymbolNode) DotFile(file string) {
	bld := &dotBuilder{
		iNodes: make(map[*IntermediateNode]bool),
		pNodes: make(map[string]bool),
		sNodes: make(map[*SymbolNode]bool),
		w:      new(bytes.Buffer),
	}
	fmt.Fprintln(bld.w, "digraph SPPF {")
	root.dot(bld)
	fmt.Fprintln(bld.w, "}")
	ioutil.WriteFile(file, bld.w.Bytes())
}

func (n *IntermediateNode) dot(bld *dotBuilder) {
	fmt.Println("in.dot", n.Label())

	fmt.Fprintf(bld.w, "%s [shape=box]\n", n.Label())

	for _, c := range n.Children {
		fmt.Fprintf(bld.w, "%s -> %s\n", n.Label(), c.Label())
		if !bld.pNodes[c.Label()] {
			c.dot(bld)
		}
	}
}

func (n *SymbolNode) dot(bld *dotBuilder) {
	fmt.Println("sn.dot", n.Label())

	fmt.Fprintln(bld.w, n.Label())
	for _, pn := range n.Children {
		fmt.Fprintf(bld.w, "%s -> %s\n", n.Label(), pn.Label())
		if !bld.pNodes[pn.Label()] {
			pn.dot(bld)
		}
	}
	for i, pn := range n.Children {
		if i > 0 {
			fmt.Fprint(bld.w, ";")
		}
		fmt.Fprintf(bld.w, "%s", pn.Label())
	}
	fmt.Fprintln(bld.w)

}

func (n *PackedNode) dot(bld *dotBuilder) {
	fmt.Println("pn.dot", n.Label(), "exist", bld.pNodes[n.Label()])

	fmt.Fprintf(bld.w, "%s [shape=box,style=rounded]\n", n.Label())
	if n.LeftChild != nil {
		n.LeftChild.dot(bld)
		fmt.Fprintf(bld.w, "%s -> %s\n", n.Label(), n.LeftChild.Label())
	}
	if n.RightChild != nil {
		n.RightChild.dot(bld)
		fmt.Fprintf(bld.w, "%s -> %s\n", n.Label(), n.RightChild.Label())
	}
	if n.LeftChild != nil && n.RightChild != nil {
		fmt.Fprintf(bld.w, "%s,%s\n", n.LeftChild.Label(), n.RightChild.Label())
	}
	bld.pNodes[n.Label()] = true
}

`
