/*
Copyright 2020 Marius Ackerman

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package lexer generates a Go lexer
package lexer

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/goccmack/gogll/v3/ast"
	"github.com/goccmack/gogll/v3/cfg"
	"github.com/goccmack/gogll/v3/lex/items"
	"github.com/goccmack/gogll/v3/symbols"
	"github.com/goccmack/goutil/ioutil"
	"github.com/goccmack/goutil/stringset"
)

type Data struct {
	Package string
	Accept  []string
	// A slice of transitions for each set
	Transitions [][]*Transition
	Tick        string
}

type Transition struct {
	Condition string
	NextState int
}

func Gen(g *ast.GoGLL, ls *items.Sets) {
	tmpl, err := template.New("lexer").Parse(tmplSrc)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	if err = tmpl.Execute(buf, getData(g, ls)); err != nil {
		panic(err)
	}
	lexFile := filepath.Join(cfg.BaseDir, "lexer", "lexer.go")
	if err = ioutil.WriteFile(lexFile, buf.Bytes()); err != nil {
		panic(err)
	}
}

// slits is the set of StringLiterals from the AST
func getAccept(ls *items.Sets, slits *stringset.StringSet) (tokTypes []string) {
	for _, s := range ls.Sets() {
		tok := s.Accept(slits)
		tokTypes = append(tokTypes, symbols.TerminalLiteralToType(tok).TypeString())
	}
	return
}

func getData(g *ast.GoGLL, ls *items.Sets) *Data {
	return &Data{
		Package:     g.Package.GetString(),
		Accept:      getAccept(ls, g.GetStringLiteralsSet()),
		Transitions: getTransitions(ls),
		Tick:        "`",
	}
}

func getTransitions(ls *items.Sets) [][]*Transition {
	trans := make([][]*Transition, len(ls.Sets()))
	for i, set := range ls.Sets() {
		trans[i] = getSetTransitions(set)
	}
	return trans
}

func getSetTransitions(set *items.Set) []*Transition {
	trans := make([]*Transition, len(set.Transitions))
	for i, t := range set.Transitions {
		trans[i] = getTransition(t)
	}
	return trans
}

func getTransition(t *items.Transition) *Transition {
	return &Transition{
		Condition: getCondition(t.Event),
		NextState: t.To.No,
	}
}

func getCondition(event ast.LexBase) string {
	switch e := event.(type) {
	case *ast.Any:
		return "true"
	case *ast.AnyOf:
		return fmt.Sprintf("any(r, %s)", e.Set)
	case *ast.CharLiteral:
		return fmt.Sprintf("r == %s", string(e.Literal))
	case *ast.Not:
		return fmt.Sprintf("not(r, %s)", e.Set)
	case *ast.UnicodeClass:
		switch e.Type {
		case ast.Letter:
			return "unicode.IsLetter(r)"
		case ast.Upcase:
			return "unicode.IsUpper(r)"
		case ast.Lowcase:
			return "unicode.IsLower(r)"
		case ast.Number:
			return "unicode.IsNumber(r)"
		case ast.Space:
			return "unicode.IsSpace(r)"
		}
		panic(fmt.Sprintf("Invalid type %d", e.Type))
	}
	panic(fmt.Sprintf("Invalid event %T", event))
}

const tmplSrc = `
// Package lexer is generated by GoGLL. Do not edit.
package lexer

import (
	// "fmt"
	"io/ioutil"
	"strings"
	"unicode"

	"{{.Package}}/token"
)

type state int

const nullState state = -1


// Lexer contains both the input slice of runes and the slice of tokens
// parsed from the input
type Lexer struct {
	// I is the input slice of runes
	I      []rune

	// Tokens is the slice of tokens constructed by the lexer from I
	Tokens []*token.Token
}

/*
NewFile constructs a Lexer created from the input file, fname. 

If the input file is a markdown file NewFile process treats all text outside
code blocks as whitespace. All text inside code blocks are treated as input text.

If the input file is a normal text file NewFile treats all text in the inputfile
as input text.
*/
func NewFile(fname string) *Lexer {
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	input := []rune(string(buf))
	if strings.HasSuffix(fname, ".md") {
		loadMd(input)
	}
	return New(input)
}

func loadMd(input []rune) {
	i := 0
	text := true
	for i < len(input) {
		if i <= len(input)-3 && input[i] == '{{.Tick}}' && input[i+1] == '{{.Tick}}' && input[i+2] == '{{.Tick}}' {
			text = !text
			for j := 0; j < 3; j++ {
				input[i+j] = ' '
			}
			i += 3
		}
		if i < len(input) {
			if text {
				if input[i] == '\n' {
					input[i] = '\n'
				} else {
					input[i] = ' '
				}
			}
			i += 1
		}
	}
}

/*
New constructs a Lexer from a slice of runes. 

All contents of the input slice are treated as input text.
*/
func New(input []rune) *Lexer {
	lex := &Lexer{
		I:      input,
		Tokens: make([]*token.Token, 0, 2048),
	}
	lext := 0
	for lext < len(lex.I) {
		for lext < len(lex.I) && unicode.IsSpace(lex.I[lext]) {
			lext++
		}
		if lext < len(lex.I) {
			tok := lex.scan(lext)
			lext = tok.Rext()
			if !tok.Suppress() {
				lex.addToken(tok)
			}
		}
	}
	lex.add(token.EOF, len(input), len(input))
	return lex
}

func (l *Lexer) scan(i int) *token.Token {
	// fmt.Printf("lexer.scan(%d)\n", i)
	s, typ, rext := nullState, token.Error, i+1
	if i < len(l.I) {
		// fmt.Printf("  rext %d, i %d\n", rext, i)
		s = nextState[0](l.I[i])
	}
	for s != nullState {
		if rext >= len(l.I) {
			typ = accept[s]
			s = nullState
		} else {
			typ = accept[s]
			s = nextState[s](l.I[rext])
			if s != nullState || typ == token.Error {
				rext++
			}
		}
	}
	tok := token.New(typ, i, rext, l.I)
	// fmt.Printf("  %s\n", tok)
	return tok
}

func escape(r rune) string {
	switch r {
	case '"':
		return "\""
	case '\\':
		return "\\\\"
	case '\r':
		return "\\r"
	case '\n':
		return "\\n"
	case '\t':
		return "\\t"
	}
	return string(r)
}

// GetLineColumn returns the line and column of rune[i] in the input
func (l *Lexer) GetLineColumn(i int) (line, col int) {
	line, col = 1, 1
	for j := 0; j < i; j++ {
		switch l.I[j] {
		case '\n':
			line++
			col = 1
		case '\t':
			col += 4
		default:
			col++
		}
	}
	return
}

// GetLineColumnOfToken returns the line and column of token[i] in the imput
func (l *Lexer) GetLineColumnOfToken(i int) (line, col int) {
	return l.GetLineColumn(l.Tokens[i].Lext())
}

// GetString returns the input string from the left extent of Token[lext] to
// the right extent of Token[rext]
func (l *Lexer) GetString(lext, rext int) string {
	return string(l.I[l.Tokens[lext].Lext():l.Tokens[rext].Rext()])
}

func (l *Lexer) add(t token.Type, lext, rext int) {
	l.addToken(token.New(t, lext, rext, l.I))
}

func (l *Lexer) addToken(tok *token.Token) {
	l.Tokens = append(l.Tokens, tok)
}

func any(r rune, set []rune) bool {
	for _, r1 := range set {
		if r == r1 {
			return true
		}
	}
	return false
}

func not(r rune, set []rune) bool {
	for _, r1 := range set {
		if r == r1 {
			return false
		}
	}
	return true
}

var accept = []token.Type{ {{range $tok := .Accept}}
	token.{{$tok}}, {{end}}
}

var nextState = []func(r rune) state{ {{range $i, $set := .Transitions}}
	// Set{{$i}}
	func(r rune) state {
		switch { {{range $cond := $set}}
		case {{$cond.Condition}}:
			return {{$cond.NextState}} {{end}}
		}
		return nullState
	}, {{end}}
}
`
