package parser

import (
	"bufio"
	"io"
)

//Parser estructura del traductors
type Parser struct {
}

//New Regresa la variable parser
func New() *Parser {
	return &Parser{}
}

//Parse hace la traduccion del codigo de la sic estandar
func (p *Parser) Parse(r io.Reader, flag bool, isXE bool) string {
	l := newLexer(bufio.NewReader(r))
	l.firstparse = flag
	l.isXE = isXE
	yyParse(l)
	return l.str
}
