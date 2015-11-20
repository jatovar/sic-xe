%{

package parser

import (
    "bufio"
    "fmt"
    "strconv"
    "strings"
)

type yylexer struct{
    src     *bufio.Reader
    buf     []byte
    empty   bool
    current byte
    str        string
    tlineidx   int
    lineidx    int
    yyerrok    bool
    err	       error
    firstparse  bool
    isXE	bool

}

func newLexer(src *bufio.Reader) (y *yylexer) {
    y = &yylexer{src: src}
    y.yyerrok = false
	y.current, _ = src.ReadByte()
	    return y
}

func (y *yylexer) getc() byte {
    var err error
    if y.current != 0 {
        y.buf = append(y.buf, y.current)
    }

    if y.current, err = y.src.ReadByte(); err == nil {

	if y.current == 10 {
		y.lineidx++
	}
    }

    return y.current
}

func (y *yylexer) Error(e string) {
	y.str += fmt.Sprintf("error linea %d: expresion inesperada" ,y.tlineidx+1)
  y.yyerrok = true

}
func (y *yylexer) Errorf(e string) {
     y.str += fmt.Sprintf(" cerca de ... %s \n",y.buf)
    y.yyerrok = true

}

func (y *yylexer) CompatibilityErr(e string) {
     y.str += fmt.Sprintf("Instruccion/Directiva incompatible %s cerca de ... %s \n",e,y.buf)
    y.yyerrok = true

}

func (y *yylexer) Lex(yylval *yySymType) int {

    c := y.current
    if y.empty {
        c, y.empty = y.getc(), false
    }
%}

%yyc c
%yyn c = y.getc()


%%
			y.tlineidx = y.lineidx
			 y.buf = y.buf[:0]
[ \t]+
(\n)+|(\r)+|(\r\n)+	return NL

"ADD"|"AND"|"COMP"|"DIV"|"J"|"JEQ"|"JGT"|"JLT"|"JSUB"|"LDA"|"LDCH"|"LDL"	return y.strtoken(yylval,NEMONICO);
"LDX"|"MUL"|"OR"|"RD"|"STA"|"STCH"|"STL"|"STSW"|"STX"|"SUB"|"TD"|"TIX"|"WD"	return y.strtoken(yylval,NEMONICO);
"ADDF"|"COMPF"|"DIVF"|"LDB"|"LDF"|"LDS"|"LDT"|"LPS"|"MULF"|"SSK"|"STB"|"STF"|"STI"|"STS"|"STT"|"SUBF" return y.strtoken(yylval,NEMONICOXE);
"RSUB"	return y.strtoken(yylval,RSUB);

"FIX"|"FLOAT"|"HIO"|"NORM"|"SIO"|"TIO" return y.strtoken(yylval,F1NEMONICO);

"ADDR"|"CLEAR"|"COMPR"|"DIVR"|"MULR"|"RMO"|"SHIFTL"|"SHIFTR"|"SUBR"|"SVC"|"TIXR" return y.strtoken(yylval,F2NEMONICO);

"RESW"	return y.token(yylval,RESW);
"RESB" 	return y.token(yylval,RESB);
"WORD" 	return y.token(yylval,WORD);
"BYTE" 	return y.token(yylval,BYTE);
"EQU"		return y.token(yylval,EQU);
"START" return y.token(yylval,START);
"BASE"	return y.token(yylval,BASE);
"END" 	return y.token(yylval,END);

"+"	return PLUS
"-"	return MINUS
"/"	return DIVIDE
"#"	return HASH
"@"	return AT
","	return COMMA
"*"	return STAR

"A"|"B"|"F"|"L"|"S"|"SW"|"T"|"PC"	return y.strtoken(yylval,REG);
"X"	return y.strtoken(yylval, REG_X);


[_a-zA-Z][a-zA-Z0-9_]*  return y.strtoken(yylval, IDENTIFIER);
[0-9]+|[1-9][0-9]+ 	return y.numtoken(yylval, NUMBER);
[0-9A-F]+(h|H) 		return y.hextoken(yylval, HEXNUMBER);
C\'[ _a-zA-Z0-9]+\' return y.charconstoken(yylval, CHARCONST);
X\'[0-9A-F]+\' 	 return y.hexconstoken(yylval, HEXCONST);



%%
    y.empty = true
    return int(c)
}
//metodo para igualar el token de la union al valor constante del token
func (y *yylexer) token(yylval *yySymType, tok int) int {
    yylval.token  = tok
    return tok
}
//metodo para tomar la cadena de las etiquetas
func (y *yylexer) strtoken(yylval *yySymType, tok int) int {
    yylval.str = string(y.buf)
    return y.token(yylval, tok)
}

//funcion que convierte un numero hexadecimal en decimal
func (y *yylexer) hextoken(yylval *yySymType, tok int) int {
  str := string(y.buf)
  str = strings.ToLower(str)
	strArray := strings.Split(str, "h")

  i, err := strconv.ParseInt(strArray[0], 16, 32)
  if(err == nil){
    yylval.number = i
  }else
    fmt.Println("Error en conversion hextoken",err)
    return y.token(yylval,tok)
}
//funcion que convierte a numero una cadena recibida
func (y *yylexer) numtoken(yylval *yySymType, tok int) int {
  str := string(y.buf)

  i, err := strconv.ParseInt(str, 10, 32)
  if(err == nil){
    yylval.number  = i
    //fmt.Println(i)

  }else
    fmt.Println("Error en conversion numtoken",err)
    return y.token(yylval,tok)
}


func (y *yylexer) charconstoken (yylval *yySymType, tok int) int{
  tam := int64(len(string(y.buf)))
  yylval.numbers.hex = false
  yylval.numbers.tam  = tam - 3
  yylval.numbers.ascii = nil
  //yylval.numbers.ascii = make([]byte,100)  //Solo 100 caracteres
  for i := 2; i < len(y.buf)-1; i++{
    yylval.numbers.ascii = append(yylval.numbers.ascii,y.buf[i])
  }
  return y.token(yylval,tok)
}

func (y *yylexer) hexconstoken (yylval *yySymType, tok int) int {
  tam := int64(len(string(y.buf)))
    yylval.numbers.hex = true
      yylval.numbers.ascii = nil
  for i := 2; i < len(y.buf)-1; i++{
    yylval.numbers.ascii = append(yylval.numbers.ascii,y.buf[i])
  }
  tam = tam - 3
  if(tam % 2 != 0){
    yylval.numbers.tam  = (tam/2 + 1);
  }else{
      yylval.numbers.tam  = tam/2;
  }

  return y.token(yylval,tok)
}
