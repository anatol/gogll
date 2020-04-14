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

import (
	"bytes"
	"fmt"
	"os"
	"sort"

	"github.com/goccmack/gogll/lexer"
	"github.com/goccmack/gogll/parser/bsr"
	"github.com/goccmack/gogll/parser/slot"
	"github.com/goccmack/gogll/parser/symbols"
	"github.com/goccmack/gogll/token"
)

var (
	cI = 0

	R *descriptors
	U *descriptors

	popped   map[poppedNode]bool
	crf      map[clusterNode][]*crfNode
	crfNodes map[crfNode]*crfNode

	lex         *lexer.Lexer
	parseErrors []*ParseError
)

func initParser(l *lexer.Lexer) {
	lex = l
	cI = 0
	R, U = &descriptors{}, &descriptors{}
	popped = make(map[poppedNode]bool)
	crf = map[clusterNode][]*crfNode{
		{symbols.NT_GoGLL, 0}: {},
	}
	crfNodes = map[crfNode]*crfNode{}
	bsr.Init(symbols.NT_GoGLL, lex.I)
	parseErrors = nil
}

func Parse(l *lexer.Lexer) (error, []*ParseError) {
	initParser(l)
	var L slot.Label
	m, cU := len(l.Tokens), 0
	ntAdd(symbols.NT_GoGLL, 0)
	// DumpDescriptors()
	for !R.empty() {
		L, cU, cI = R.remove()

		// fmt.Println()
		// fmt.Printf("L:%s, cI:%d, I[cI]:%s, cU:%d\n", L, cI, nextI, cU)
		// DumpDescriptors()

		switch L {
		case slot.GoGLL0R0: // GoGLL : ∙NT_Package NT_Rules

			call(slot.GoGLL0R1, cU, cI)
		case slot.GoGLL0R1: // GoGLL : NT_Package ∙NT_Rules

			if !testSelect(slot.GoGLL0R1) {
				parseError(slot.GoGLL0R1, cI)
				break
			}

			call(slot.GoGLL0R2, cU, cI)
		case slot.GoGLL0R2: // GoGLL : NT_Package NT_Rules ∙

			if follow(symbols.NT_GoGLL) {
				rtn(symbols.NT_GoGLL, cU, cI)
			}
		case slot.NT0R0: // NT : ∙T_3

			bsr.Add(slot.NT0R1, cU, cI, cI+1)
			cI++
			if follow(symbols.NT_NT) {
				rtn(symbols.NT_NT, cU, cI)
			}
		case slot.Symbol0R0: // Symbol : ∙NT_NT

			call(slot.Symbol0R1, cU, cI)
		case slot.Symbol0R1: // Symbol : NT_NT ∙

			if follow(symbols.NT_Symbol) {
				rtn(symbols.NT_Symbol, cU, cI)
			}
		case slot.Symbol1R0: // Symbol : ∙T_6

			bsr.Add(slot.Symbol1R1, cU, cI, cI+1)
			cI++
			if follow(symbols.NT_Symbol) {
				rtn(symbols.NT_Symbol, cU, cI)
			}
		case slot.Symbol2R0: // Symbol : ∙T_5

			bsr.Add(slot.Symbol2R1, cU, cI, cI+1)
			cI++
			if follow(symbols.NT_Symbol) {
				rtn(symbols.NT_Symbol, cU, cI)
			}
		case slot.Package0R0: // Package : ∙T_4 T_5

			bsr.Add(slot.Package0R1, cU, cI, cI+1)
			cI++
			if !testSelect(slot.Package0R1) {
				parseError(slot.Package0R1, cI)
				break
			}

			bsr.Add(slot.Package0R2, cU, cI, cI+1)
			cI++
			if follow(symbols.NT_Package) {
				rtn(symbols.NT_Package, cU, cI)
			}
		case slot.Rules0R0: // Rules : ∙NT_Rule

			call(slot.Rules0R1, cU, cI)
		case slot.Rules0R1: // Rules : NT_Rule ∙

			if follow(symbols.NT_Rules) {
				rtn(symbols.NT_Rules, cU, cI)
			}
		case slot.Rules1R0: // Rules : ∙NT_Rules NT_Rule

			call(slot.Rules1R1, cU, cI)
		case slot.Rules1R1: // Rules : NT_Rules ∙NT_Rule

			if !testSelect(slot.Rules1R1) {
				parseError(slot.Rules1R1, cI)
				break
			}

			call(slot.Rules1R2, cU, cI)
		case slot.Rules1R2: // Rules : NT_Rules NT_Rule ∙

			if follow(symbols.NT_Rules) {
				rtn(symbols.NT_Rules, cU, cI)
			}
		case slot.Rule0R0: // Rule : ∙NT_NT T_0 NT_Alternates T_1

			call(slot.Rule0R1, cU, cI)
		case slot.Rule0R1: // Rule : NT_NT ∙T_0 NT_Alternates T_1

			if !testSelect(slot.Rule0R1) {
				parseError(slot.Rule0R1, cI)
				break
			}

			bsr.Add(slot.Rule0R2, cU, cI, cI+1)
			cI++
			if !testSelect(slot.Rule0R2) {
				parseError(slot.Rule0R2, cI)
				break
			}

			call(slot.Rule0R3, cU, cI)
		case slot.Rule0R3: // Rule : NT_NT T_0 NT_Alternates ∙T_1

			if !testSelect(slot.Rule0R3) {
				parseError(slot.Rule0R3, cI)
				break
			}

			bsr.Add(slot.Rule0R4, cU, cI, cI+1)
			cI++
			if follow(symbols.NT_Rule) {
				rtn(symbols.NT_Rule, cU, cI)
			}
		case slot.Alternates0R0: // Alternates : ∙NT_Alternate

			call(slot.Alternates0R1, cU, cI)
		case slot.Alternates0R1: // Alternates : NT_Alternate ∙

			if follow(symbols.NT_Alternates) {
				rtn(symbols.NT_Alternates, cU, cI)
			}
		case slot.Alternates1R0: // Alternates : ∙NT_Alternates T_7 NT_Alternate

			call(slot.Alternates1R1, cU, cI)
		case slot.Alternates1R1: // Alternates : NT_Alternates ∙T_7 NT_Alternate

			if !testSelect(slot.Alternates1R1) {
				parseError(slot.Alternates1R1, cI)
				break
			}

			bsr.Add(slot.Alternates1R2, cU, cI, cI+1)
			cI++
			if !testSelect(slot.Alternates1R2) {
				parseError(slot.Alternates1R2, cI)
				break
			}

			call(slot.Alternates1R3, cU, cI)
		case slot.Alternates1R3: // Alternates : NT_Alternates T_7 NT_Alternate ∙

			if follow(symbols.NT_Alternates) {
				rtn(symbols.NT_Alternates, cU, cI)
			}
		case slot.Alternate0R0: // Alternate : ∙NT_Symbols

			call(slot.Alternate0R1, cU, cI)
		case slot.Alternate0R1: // Alternate : NT_Symbols ∙

			if follow(symbols.NT_Alternate) {
				rtn(symbols.NT_Alternate, cU, cI)
			}
		case slot.Alternate1R0: // Alternate : ∙T_2

			bsr.Add(slot.Alternate1R1, cU, cI, cI+1)
			cI++
			if follow(symbols.NT_Alternate) {
				rtn(symbols.NT_Alternate, cU, cI)
			}
		case slot.Symbols0R0: // Symbols : ∙NT_Symbol

			call(slot.Symbols0R1, cU, cI)
		case slot.Symbols0R1: // Symbols : NT_Symbol ∙

			if follow(symbols.NT_Symbols) {
				rtn(symbols.NT_Symbols, cU, cI)
			}
		case slot.Symbols1R0: // Symbols : ∙NT_Symbols NT_Symbol

			call(slot.Symbols1R1, cU, cI)
		case slot.Symbols1R1: // Symbols : NT_Symbols ∙NT_Symbol

			if !testSelect(slot.Symbols1R1) {
				parseError(slot.Symbols1R1, cI)
				break
			}

			call(slot.Symbols1R2, cU, cI)
		case slot.Symbols1R2: // Symbols : NT_Symbols NT_Symbol ∙

			if follow(symbols.NT_Symbols) {
				rtn(symbols.NT_Symbols, cU, cI)
			}

		default:
			panic("This must not happen")
		}
	}
	if !bsr.Contain(symbols.NT_GoGLL, 0, m) {
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
	for cn, nds := range crf {
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

var first = []map[token.Type]string{
	// Alternate : ∙NT_Symbols
	map[token.Type]string{},
	// Alternate : NT_Symbols ∙
	map[token.Type]string{
		token.Type7: "|",
		token.Type1: ";",
	},
	// Alternate : ∙T_2
	map[token.Type]string{},
	// Alternate : T_2 ∙
	map[token.Type]string{
		token.Type1: ";",
		token.Type7: "|",
	},
	// Alternates : ∙NT_Alternate
	map[token.Type]string{},
	// Alternates : NT_Alternate ∙
	map[token.Type]string{
		token.Type1: ";",
		token.Type7: "|",
	},
	// Alternates : ∙NT_Alternates T_7 NT_Alternate
	map[token.Type]string{},
	// Alternates : NT_Alternates ∙T_7 NT_Alternate
	map[token.Type]string{},
	// Alternates : NT_Alternates T_7 ∙NT_Alternate
	map[token.Type]string{},
	// Alternates : NT_Alternates T_7 NT_Alternate ∙
	map[token.Type]string{
		token.Type1: ";",
		token.Type7: "|",
	},
	// GoGLL : ∙NT_Package NT_Rules
	map[token.Type]string{},
	// GoGLL : NT_Package ∙NT_Rules
	map[token.Type]string{},
	// GoGLL : NT_Package NT_Rules ∙
	map[token.Type]string{
		token.EOF: "EOF",
	},
	// NT : ∙T_3
	map[token.Type]string{},
	// NT : T_3 ∙
	map[token.Type]string{
		token.Type0: ":",
		token.Type7: "|",
		token.Type1: ";",
		token.Type3: "nt",
		token.Type6: "tokid",
		token.Type5: "string_lit",
	},
	// Package : ∙T_4 T_5
	map[token.Type]string{},
	// Package : T_4 ∙T_5
	map[token.Type]string{},
	// Package : T_4 T_5 ∙
	map[token.Type]string{
		token.Type3: "nt",
	},
	// Rule : ∙NT_NT T_0 NT_Alternates T_1
	map[token.Type]string{},
	// Rule : NT_NT ∙T_0 NT_Alternates T_1
	map[token.Type]string{},
	// Rule : NT_NT T_0 ∙NT_Alternates T_1
	map[token.Type]string{},
	// Rule : NT_NT T_0 NT_Alternates ∙T_1
	map[token.Type]string{},
	// Rule : NT_NT T_0 NT_Alternates T_1 ∙
	map[token.Type]string{
		token.EOF:   "EOF",
		token.Type3: "nt",
	},
	// Rules : ∙NT_Rule
	map[token.Type]string{},
	// Rules : NT_Rule ∙
	map[token.Type]string{
		token.EOF:   "EOF",
		token.Type3: "nt",
	},
	// Rules : ∙NT_Rules NT_Rule
	map[token.Type]string{},
	// Rules : NT_Rules ∙NT_Rule
	map[token.Type]string{},
	// Rules : NT_Rules NT_Rule ∙
	map[token.Type]string{
		token.EOF:   "EOF",
		token.Type3: "nt",
	},
	// Symbol : ∙NT_NT
	map[token.Type]string{},
	// Symbol : NT_NT ∙
	map[token.Type]string{
		token.Type5: "string_lit",
		token.Type7: "|",
		token.Type1: ";",
		token.Type3: "nt",
		token.Type6: "tokid",
	},
	// Symbol : ∙T_6
	map[token.Type]string{},
	// Symbol : T_6 ∙
	map[token.Type]string{
		token.Type1: ";",
		token.Type3: "nt",
		token.Type6: "tokid",
		token.Type5: "string_lit",
		token.Type7: "|",
	},
	// Symbol : ∙T_5
	map[token.Type]string{},
	// Symbol : T_5 ∙
	map[token.Type]string{
		token.Type7: "|",
		token.Type1: ";",
		token.Type3: "nt",
		token.Type6: "tokid",
		token.Type5: "string_lit",
	},
	// Symbols : ∙NT_Symbol
	map[token.Type]string{},
	// Symbols : NT_Symbol ∙
	map[token.Type]string{
		token.Type6: "tokid",
		token.Type5: "string_lit",
		token.Type7: "|",
		token.Type1: ";",
		token.Type3: "nt",
	},
	// Symbols : ∙NT_Symbols NT_Symbol
	map[token.Type]string{},
	// Symbols : NT_Symbols ∙NT_Symbol
	map[token.Type]string{},
	// Symbols : NT_Symbols NT_Symbol ∙
	map[token.Type]string{
		token.Type3: "nt",
		token.Type6: "tokid",
		token.Type5: "string_lit",
		token.Type7: "|",
		token.Type1: ";",
	},
}

var followSets = []map[token.Type]string{
	// GoGLL
	map[token.Type]string{
		token.EOF: "EOF",
	},
	// NT
	map[token.Type]string{
		token.Type7: "|",
		token.Type1: ";",
		token.Type3: "nt",
		token.Type6: "tokid",
		token.Type5: "string_lit",
		token.Type0: ":",
	},
	// Rule
	map[token.Type]string{
		token.EOF:   "EOF",
		token.Type3: "nt",
	},
	// Alternates
	map[token.Type]string{
		token.Type1: ";",
		token.Type7: "|",
	},
	// Alternate
	map[token.Type]string{
		token.Type1: ";",
		token.Type7: "|",
	},
	// Symbols
	map[token.Type]string{
		token.Type7: "|",
		token.Type1: ";",
		token.Type3: "nt",
		token.Type6: "tokid",
		token.Type5: "string_lit",
	},
	// Symbol
	map[token.Type]string{
		token.Type7: "|",
		token.Type1: ";",
		token.Type3: "nt",
		token.Type6: "tokid",
		token.Type5: "string_lit",
	},
	// Package
	map[token.Type]string{
		token.Type3: "nt",
	},
	// Rules
	map[token.Type]string{
		token.EOF:   "EOF",
		token.Type3: "nt",
	},
}

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
