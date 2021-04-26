package items

import (
	"testing"

	"github.com/goccmack/gogll/lexer"
	"github.com/goccmack/gogll/parser"
	"github.com/goccmack/gogll/parser/bsr"
	"github.com/goccmack/gogll/v3/ast"
)

const src = `package "names"
qualifiedName : letter {letter|number|'_'} <'.' <letter|number|'_'>> ;
`

func Test1(t *testing.T) {
	lex := lexer.New([]rune(src))
	parser.Parse(lex)
	g := ast.Build(bsr.GetRoot(), lex)

	New(g)
}
