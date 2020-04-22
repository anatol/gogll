package lexer

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/goccmack/gogll/ast"
	"github.com/goccmack/gogll/im/tokens"
	"github.com/goccmack/gogll/lex/items"
	"github.com/goccmack/goutil/ioutil"
)

type Data struct {
	Package string
	Accept  []string
	// A slice of transitions for each set
	Transitions [][]*Transition
}

type Transition struct {
	Condition string
	NextState int
}

func Gen(lexDir string, g *ast.GoGLL, ls *items.Sets, ts *tokens.Tokens) {
	tmpl, err := template.New("lexer").Parse(tmplSrc)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	if err = tmpl.Execute(buf, getData(g, ls, ts)); err != nil {
		panic(err)
	}
	if err = ioutil.WriteFile(filepath.Join(lexDir, "lexer.go"), buf.Bytes()); err != nil {
		panic(err)
	}
}

func getAccept(ls *items.Sets, ts *tokens.Tokens) (tokTypes []string) {
	for _, s := range ls.Sets() {
		tok := s.Accept()
		tokTypes = append(tokTypes, ts.LiteralToString[tok])
	}
	return
}

func getData(g *ast.GoGLL, ls *items.Sets, ts *tokens.Tokens) *Data {
	return &Data{
		Package:     g.Package.GetString(),
		Accept:      getAccept(ls, ts),
		Transitions: getTransitions(ls),
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
		return fmt.Sprintf("r == %s", string(e.Literal()))
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
	"fmt"
	"io/ioutil"
	"strings"
	"unicode"

	"github.com/goccmack/goutil/md"

	"github.com/goccmack/gogll/token"
)

type state int

const nullState state = -1

type Lexer struct {
	I      []rune
	Tokens []*token.Token
}

func NewFile(fname string) *Lexer {
	if strings.HasSuffix(fname, ".md") {
		src, err := md.GetSource(fname)
		if err != nil {
			panic(err)
		}
		return New([]rune(src))
	}
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	return New([]rune(string(buf)))
}

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
		if lext >= len(lex.I) {
			lex.add(token.EOF, len(input), len(input))
		} else {
			tok := lex.scan(lext)
			lext = tok.Rext
			lex.addToken(tok)
		}
	}
	return lex
}

func (l *Lexer) scan(i int) *token.Token {
	s, tok := state(0), token.New(token.Error, i, i, nil)
	for s != nullState {
		if tok.Rext >= len(l.I) {
			s = nullState
		} else {
			s = nextState[s](l.I[tok.Rext])
			tok.Rext++
			tok.Type = accept[s]
		}
	}
	return tok
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

func (l *Lexer) GetLineColumnOfToken(i int) (line, col int) {
	return l.GetLineColumn(l.Tokens[i].Lext)
}

// GetString returns the input string from the left extent of Token[lext] to
// the right extent of Token[rext]
// func (l *Lexer) GetString(lext, rext int) string {
// 	return string(l.I[l.Tokens[lext].Lext:l.Tokens[rext].Rext])
// }

func (l *Lexer) add(t token.Type, lext, rext int) {
	l.addToken(token.New(t, lext, rext, l.I[lext:rext]))
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
		panic(fmt.Sprintf("Unexpected rune '%c' in state S{{$i}}", r))
	}, {{end}}
}
`
