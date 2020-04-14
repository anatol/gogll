//  Copyright 2019 Marius Ackerman
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package parser

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/goccmack/gogll/ast"
	"github.com/goccmack/gogll/frstflw"
	"github.com/goccmack/gogll/gen/golang/parser/bsr"
	"github.com/goccmack/gogll/gen/golang/parser/slots"
	"github.com/goccmack/gogll/gen/golang/parser/symbols"
	"github.com/goccmack/gogll/gslot"
	"github.com/goccmack/goutil/ioutil"
)

/*** Main parser section ***/

type gen struct {
	g  *ast.GoGLL
	gs *gslot.GSlot
	ff *frstflw.FF
}

func Gen(parserDir string, g *ast.GoGLL, gs *gslot.GSlot, ff *frstflw.FF) {
	gn := &gen{g, gs, ff}
	gn.genParser(parserDir)
	bsr.Gen(filepath.Join(parserDir, "bsr", "bsr.go"), g.Package.GetString())
	slots.Gen(filepath.Join(parserDir, "slot", "slot.go"), g, gs, ff)
	symbols.Gen(filepath.Join(parserDir, "symbols", "symbols.go"), g)
}

func (g *gen) genParser(parserDir string) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("Parser Main Template").Parse(mainTemplate)
	if err != nil {
		parseErrorError(err)
	}
	data := g.getData(parserDir)
	if err = tmpl.Execute(buf, data); err != nil {
		parseErrorError(err)
	}
	fmtSrc, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Printf("Error formatting generated parsers: %s\n", err)
		fmtSrc = buf.Bytes()
	}
	fname := path.Join(parserDir, "parser.go")
	if err := ioutil.WriteFile(fname, []byte(fmtSrc)); err != nil {
		parseErrorError(err)
	}
}

type Data struct {
	Package     string
	StartSymbol string
	CodeX       string
	TestSelect  string
}

func (g *gen) getData(baseDir string) *Data {
	data := &Data{
		Package:     g.g.Package.GetString(),
		StartSymbol: g.g.StartSymbol(),
		CodeX:       g.genAlternatesCode(),
		TestSelect:  g.genTestSelect(),
	}
	return data
}

func parseErrorError(err error) {
	fmt.Printf("Error generating parser: %s\n", err)
	panic("fix me")
	os.Exit(1)
}

const mainTemplate = `
// Package parser is generated by gogll. Do not edit.
//
//  Copyright 2019 Marius Ackerman
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
package parser

import(
	"bytes"
	"fmt"
	"os"
	"sort"

	"{{.Package}}/parser/bsr"
	"{{.Package}}/lexer"
	"{{.Package}}/parser/bsr"
    "{{.Package}}/parser/slot"
    "{{.Package}}/parser/symbols"
    "{{.Package}}/token"
)

var (
	cI    = 0

	R *descriptors
	U *descriptors

	popped 		map[poppedNode]bool
	crf			map[clusterNode][]*crfNode
	crfNodes	map[crfNode]*crfNode

	lex         *lexer.Lexer
	parseErrors []*ParseError
)

func initParser(l *lexer.Lexer) {
	lex = l
	cI = 0
	R, U = &descriptors{}, &descriptors{}
	popped = make(map[poppedNode]bool)
	crf = map[clusterNode][]*crfNode{
		{symbols.NT_{{.StartSymbol}}, 0}:{},
	}
	crfNodes = map[crfNode]*crfNode{}
	bsr.Init(symbols.NT_{{.StartSymbol}}, lex.I)
	parseErrors = nil
}

func Parse(l *lexer.Lexer) (error, []*ParseError) {
	initParser(l)
	var L slot.Label
	m, cU := len(l.Tokens), 0
	ntAdd(symbols.NT_{{.StartSymbol}}, 0)
	// DumpDescriptors()
	for !R.empty() {
		L, cU, cI = R.remove()

		// fmt.Println()
		// fmt.Printf("L:%s, cI:%d, I[cI]:%s, cU:%d\n", L, cI, nextI, cU)
		// DumpDescriptors()

		switch L { 
{{.CodeX}}

		default:
			panic("This must not happen")
		}
	}
	if !bsr.Contain(symbols.NT_{{.StartSymbol}},0,m) {
		sortParseErrors()
		err := fmt.Errorf("Error: Parse Failed right extent=%d, m=%d", 
			bsr.GetRightExtent(), len(l.Tokens))
		return err, parseErrors
	}
	return nil, nil
}

func ntAdd(nt symbols.NT, j int) {
	// fmt.Printf("ntAdd(%s, %d)\n", nt, j)
	failed := true
	for _, l := range slot.GetAlternates(nt) {
		if testSelect(l) {
			dscAdd(l, j, j)
		} else {
			failed = false
		}
	}
	if failed {
		for _, l := range slot.GetAlternates(nt) {
			parseError(l, j)
		}
	}
}

/*** Call Return Forest ***/

type poppedNode struct {
	X    symbols.NT
	k, j int
}

type clusterNode struct {
	X symbols.NT
	k int
}

type crfNode struct {
	L slot.Label
	i int
}

/*
suppose that L is Y ::=αX ·β
if there is no CRF node labelled (L,i) 
	create one let u be the CRF node labelled (L,i)
if there is no CRF node labelled (X, j) { 
	create a CRF node v labelled (X, j) 
	create an edge from v to u 
	ntAdd(X, j) 
} else { 
	let v be the CRF node labelled (X, j) 
	if there is not an edge from v to u {
		create an edge from v to u 
		for all ((X, j,h)∈P) {
			dscAdd(L, i, h); 
			bsrAdd(L, i, j, h) 
		} 
	} 
}
*/
func call(L slot.Label, i, j int) {
	// fmt.Printf("call(%s,%d,%d)\n", L,i,j)
	u, exist := crfNodes[crfNode{L, i}]
	// fmt.Printf("  u exist=%t\n", exist)
	if !exist {
		u = &crfNode{L, i}
		crfNodes[*u] = u
	}
	X := L.Symbols()[L.Pos()-1].(symbols.NT)
	ndV := clusterNode{X, j}
	v, exist := crf[ndV]
	if !exist {
		// fmt.Println("  v !exist")
		crf[ndV] = []*crfNode{u}
		ntAdd(X, j)
	} else {
		// fmt.Println("  v exist")
		if !existEdge(v, u) {
			// fmt.Printf("  !existEdge(%v)\n", u)
			crf[ndV] = append(v, u)
			// fmt.Printf("|popped|=%d\n", len(popped))
			for pnd, _ := range popped {
				if pnd.X == X && pnd.k == j {
					dscAdd(L, i, pnd.j)
					bsr.Add(L, i, j, pnd.j)
				}
			}
		}
	}
}

func existEdge(nds []*crfNode, nd *crfNode) bool {
	for _, nd1 := range nds {
		if nd1 == nd {
			return true
		}
	}
	return false
}

func rtn(X symbols.NT, k, j int) {
	// fmt.Printf("rtn(%s,%d,%d)\n", X,k,j)
	p := poppedNode{X, k, j}
	if _, exist := popped[p]; !exist {
		popped[p] = true
		for _, nd := range crf[clusterNode{X, k}] {
			dscAdd(nd.L, nd.i, j)
			bsr.Add(nd.L, nd.i, k, j)
		}
	}
}

func CRFString() string {
	buf := new(bytes.Buffer)
	buf.WriteString("CRF: {")
	for cn, nds := range crf{
		for _, nd := range nds {
			fmt.Fprintf(buf, "%s->%s, ", cn, nd)
		}
	}
	buf.WriteString("}")
	return buf.String()
}

func (cn clusterNode) String() string {
	return fmt.Sprintf("(%s,%d)", cn.X, cn.k)
}

func (n crfNode) String() string {
	return fmt.Sprintf("(%s,%d)", n.L.String(), n.i)
}

func PoppedString() string {
	buf := new(bytes.Buffer)
	buf.WriteString("Popped: {")
	for p, _ := range popped {
		fmt.Fprintf(buf, "(%s,%d,%d) ", p.X, p.k, p.j)
	}
	buf.WriteString("}")
	return buf.String()
}

/*** descriptors ***/

type descriptors struct {
	set []*descriptor
}

func (ds *descriptors) contain(d *descriptor) bool {
	for _, d1 := range ds.set {
		if d1 == d {
			return true
		}
	}
	return false
}

func (ds *descriptors) empty() bool {
	return len(ds.set) == 0
}

func (ds *descriptors) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString("{")
	for i, d := range ds.set {
		if i > 0 {
			buf.WriteString("; ")
		}
		fmt.Fprintf(buf, "%s", d)
	}
	buf.WriteString("}")
	return buf.String()
}

type descriptor struct {
	L slot.Label
	k int
	i int
}

func (d *descriptor) String() string {
	return fmt.Sprintf("%s,%d,%d", d.L, d.k, d.i)
}

func dscAdd(L slot.Label, k, i int) {
	// fmt.Printf("dscAdd(%s,%d,%d)\n", L, k, i)
	d := &descriptor{L, k, i}
	if !U.contain(d) {
		R.set = append(R.set, d)
		U.set = append(U.set, d)
	}
}

func (ds *descriptors) remove() (L slot.Label, k, i int) {
	d := ds.set[len(ds.set)-1]
	ds.set = ds.set[:len(ds.set)-1]
	// fmt.Printf("remove: %s,%d,%d\n", d.L, d.k, d.i)
	return d.L, d.k, d.i
}

func DumpDescriptors() {
	DumpR()
	DumpU()
}

func DumpR() {
	fmt.Println("R:")
	for _, d := range R.set {
		fmt.Printf(" %s\n", d)
	}
}

func DumpU() {
	fmt.Println("U:")
	for _, d := range U.set {
		fmt.Printf(" %s\n", d)
	}
}

/*** TestSelect ***/

func follow(nt symbols.NT) bool {
    _, exist := followSets[nt][lex.Tokens[cI].Type]
    return exist
}

func testSelect(l slot.Label) bool {
    _, exist := first[l][lex.Tokens[cI].Type]
    return exist
}

{{.TestSelect}}

	
/*** Errors ***/

type ParseError struct {
	Slot         slot.Label
	Token        *token.Token
	Line, Column int
}

func (pe *ParseError) String() string {
	return fmt.Sprintf("Parse Error: %s I[cI]=%s at line %d col %d", 
		pe.Slot, pe.Token, pe.Line, pe.Column)
}

func parseError(slot slot.Label, i int) {
	pe := &ParseError{Slot: slot, Token: lex.Tokens[i]}
	parseErrors = append(parseErrors, pe)
}

func sortParseErrors() {
	sort.Slice(parseErrors,
		func(i, j int) bool {
			return parseErrors[j].Token.Lext < parseErrors[i].Token.Lext
		})
	for _, pe := range parseErrors {
		pe.Line, pe.Column = lex.GetLineColumn(pe.Token.Lext)
	}
}

func parseErrorError(err error) {
	fmt.Printf("Error: %s\n", err)
	os.Exit(1)
}
`
