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

// Package token generates a Go token package
package token

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/goccmack/gogll/ast"
	"github.com/goccmack/gogll/cfg"
	"github.com/goccmack/gogll/gen/golang/utils"
	"github.com/goccmack/gogll/symbols"
	"github.com/goccmack/goutil/ioutil"
)

type Data struct {
	Types        []*TypeDef
	TypeToString []string
}

type TypeDef struct {
	Name, Comment string
	Suppress      bool
}

func Gen(g *ast.GoGLL) {
	tmpl, err := template.New("Token").Parse(tmplSrc)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, getData())
	if err != nil {
		panic(err)
	}
	if err = ioutil.WriteFile(tokenFile(g.Package.GetString()), buf.Bytes()); err != nil {
		panic(err)
	}
}

func getData() *Data {
	return &Data{
		Types:        getTypes(),
		TypeToString: symbols.GetTerminalTypeStrings(),
	}
}

func getTypes() (types []*TypeDef) {
	for _, t := range symbols.GetTerminals() {
		types = append(types,
			&TypeDef{
				Name:     t.TypeString(),
				Comment:  utils.Escape(t.Literal()),
				Suppress: t.Suppress(),
			})
	}
	return
}

func tokenFile(pkg string) string {
	return filepath.Join(cfg.BaseDir, "token", "token.go")
}

const tmplSrc = `
// Package token is generated by GoGLL. Do not edit
package token

import(
    "fmt"
)

// Token is returned by the lexer for every scanned lexical token
type Token struct {
    typ        Type
    lext, rext int
    input      []rune
}

/*
New returns a new token.
lext is the left extent and rext the right extent of the token in the input.
input is the input slice scanned by the lexer.
*/
func New(t Type, lext, rext int, input []rune) *Token {
    return &Token{
        typ:   t,
        lext:  lext,
        rext:  rext,
        input: input,
    }
}

// GetLineColumn returns the line and column of the left extent of t
func (t *Token) GetLineColumn() (line, col int) {
    line, col = 1, 1
    for j := 0; j < t.lext; j++ {
        switch t.input[j] {
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

// GetInput returns the input from which t was parsed.
func (t *Token) GetInput() []rune {
    return t.input
}

// Lext returns the left extent of t
func (t *Token) Lext() int {
    return t.lext
}

// Literal returns the literal runes of t scanned by the lexer
func (t *Token) Literal() []rune {
    return t.input[t.lext:t.rext]
}

// LiteralString returns string(t.Literal())
func (t *Token) LiteralString() string {
    return string(t.Literal())
}

// LiteralStripEscape returns the literal runes of t scanned by the lexer
func (t *Token) LiteralStripEscape() []rune {
	lit := t.Literal()
	strip := make([]rune, 0, len(lit))
	for i := 0; i < len(lit); i++ {
		if lit[i] == '\\' {
			i++
			switch lit[i] {
			case 't':
				strip = append(strip, '\t')
			case 'r':
				strip = append(strip, '\r')
			case 'n':
				strip = append(strip, '\r')
			default:
				strip = append(strip, lit[i])
			}
		} else {
			strip = append(strip, lit[i])
		}
	}
	return strip
}

// LiteralStringStripEscape returns string(t.LiteralStripEscape())
func (t *Token) LiteralStringStripEscape() string {
	return string(t.LiteralStripEscape())
}

// Rext returns the right extent of t in the input
func (t *Token) Rext() int {
    return t.rext
}

func (t *Token) String() string {
    return fmt.Sprintf("%s (%d,%d) %s",
        t.TypeID(), t.lext, t.rext, t.LiteralString())
}

// Suppress returns true iff t is suppressed by the lexer
func (t *Token) Suppress() bool {
	return Suppress[t.typ]
}

// Type returns the token Type of t
func (t *Token) Type() Type {
    return t.typ
}

// TypeID returns the token Type ID of t. 
// This may be different from the literal of token t.
func (t *Token) TypeID() string {
    return t.Type().ID()
}

// Type is the token type
type Type int

func (t Type) String() string {
    return TypeToString[t]
}

// ID returns the token type ID of token Type t
func (t Type) ID() string {
    return TypeToID[t]
}


const({{range $i, $typ := .Types}}
    {{$typ.Name}} {{if eq $i 0}} Type = iota {{end}} // {{$typ.Comment}} {{end}}
)

var TypeToString = []string{ {{range $str := .TypeToString}}
    "{{$str}}",{{end}}
}

var StringToType = map[string] Type { {{range $typ := .TypeToString}}
    "{{$typ}}" : {{$typ}}, {{end}}
}

var TypeToID = []string { {{range $typ := .Types}}
    "{{$typ.Comment}}", {{end}}
}

var Suppress = []bool { {{range $typ := .Types}}
    {{$typ.Suppress}}, {{end}}
}

`
