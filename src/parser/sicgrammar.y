%{
package parser

import (

	"util"

)

var pc int64
var bRegister int64
var format4 bool
type numbers struct {
	tam int64
	ascii []byte
	hex bool

}

%}

%union{
	str string
	token int
	number int64
	numbers numbers
}

%token <str> NEMONICO, F1NEMONICO, F2NEMONICO, RSUB, NEMONICOXE
%token <token> START, END, BYTE, WORD, RESB, RESW, EQU, BASE
%token <token> PLUS, MINUS, DIVIDE, HASH, AT, NL, COMMA, STAR
%token <str> IDENTIFIER
%token <number> HEXNUMBER, NUMBER
%token <numbers> HEXCONST, CHARCONST
%token <str> REG_X, REG

%type <numbers>  constantes
%type <number> indexado numerosDir numeros
%type <str> numIdent identificador

%start prod_principal

%%
	prod_principal:
		instr_start instrucciones
	;

	instr_start:
		nlinea_op IDENTIFIER START numerosDir nlinea_op {
			if yylex.(*yylexer).firstparse == true	{
				pc = $4
				format4 = false
				util.GetAssembler().Progname = $2
				util.GetAssembler().StartAddrs = $4
				util.GetTabSim().AddProgpc(int64(yylex.(*yylexer).tlineidx) - 1,pc)
				util.GetTabSim().AddProgpc(int64(yylex.(*yylexer).tlineidx) ,pc)
			}
		}
	;
	instrucciones:
		/*epsilon*/
	|	etiqueta_op formato {
			if yylex.(*yylexer).firstparse == true {
				if yylex.(*yylexer).yyerrok == false {
					util.GetTabSim().AddProgpc(int64(yylex.(*yylexer).tlineidx) ,pc)
				}else{
					yylex.(*yylexer).yyerrok = false
				}
			}
		}
	instrucciones

	|	etiqueta_op directiva {
			if yylex.(*yylexer).firstparse == true {
				if yylex.(*yylexer).yyerrok == false {
					util.GetTabSim().AddProgpc(int64(yylex.(*yylexer).tlineidx) ,pc)
				}else{
					yylex.(*yylexer).yyerrok = false
				}
			}
		}
		instrucciones
	;
	instr_end:
		END nlinea_op
	|	END IDENTIFIER nlinea_op {
			util.GetAssembler().Endinstr($2)
		}
	;
	formato:
		error {
			yylex.(*yylexer).Errorf("aaaa")
			goto yydefault
		}
	|	f1 {
			if yylex.(*yylexer).isXE {
				pc += 1
			}else{
				if yylex.(*yylexer).isXE == false {
					yylex.(*yylexer).CompatibilityErr("de formato 1")
					goto yydefault
				}
			}
		}
	|	f2 {
			if yylex.(*yylexer).isXE {
				pc += 2
			}else{
				if yylex.(*yylexer).isXE == false {
					yylex.(*yylexer).CompatibilityErr("de formato 2")
					goto yydefault
				}
			}
					}
	|	f3 {
			pc += 3
		}
	| 	f4 {
			if yylex.(*yylexer).isXE {
				pc += 4
			}else{
				if yylex.(*yylexer).isXE == false {
						yylex.(*yylexer).CompatibilityErr("de formato 4")
						goto yydefault
				}
			}
		}
	|	instr_end
	;
	f1:
		F1NEMONICO nlinea_op {
			if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true {
				util.GetAssembler().Assemblef1(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1))
			}
		}
	;
	f2:
			F2NEMONICO numeros COMMA numeros nlinea_op	{
							if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true {
								if(util.GetAssembler().IsValidRegister($2) && util.GetAssembler().IsValidRegister($4)){
									util.GetAssembler().Assemblef2(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), $2, $4)
							}else{
								yylex.(*yylexer).Error("error de registro")
								pc -= 2
								goto yydefault
							}
						}
			}
		|	F2NEMONICO REG nlinea_op	{
						if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true {
							util.GetAssembler().Assemblef2(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), util.GetAssembler().GetRegEntry($2), 0)
						}
		}
		|	F2NEMONICO REG COMMA REG nlinea_op	{
						if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true {
							util.GetAssembler().Assemblef2(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), util.GetAssembler().GetRegEntry($2), util.GetAssembler().GetRegEntry($4))
						}
		}
		|	F2NEMONICO REG COMMA numeros nlinea_op{
						if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true {
								if(util.GetAssembler().IsValidRegister($4)){
										util.GetAssembler().Assemblef2(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), util.GetAssembler().GetRegEntry($2), $4)
								}else{
										yylex.(*yylexer).Error("error de registro")
										pc -= 2
										goto yydefault
								}
						}
		}
		|	F2NEMONICO REG_X nlinea_op{
						if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true {
								util.GetAssembler().Assemblef2(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), util.GetAssembler().GetRegEntry($2), 0)
						}
		}
		|	F2NEMONICO REG_X COMMA REG_X nlinea_op{
					if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true {
						util.GetAssembler().Assemblef2(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), util.GetAssembler().GetRegEntry($2), util.GetAssembler().GetRegEntry($4))
					}
		}
		|	F2NEMONICO REG_X COMMA numeros nlinea_op{
					if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true {
							if(util.GetAssembler().IsValidRegister($4)){
									util.GetAssembler().Assemblef2(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), util.GetAssembler().GetRegEntry($2), $4)
							}else{
								yylex.(*yylexer).Error("error de registro")
								pc -= 2
								goto yydefault
							}
					}
		}
		| F2NEMONICO REG_X COMMA REG nlinea_op{
			if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true {
				util.GetAssembler().Assemblef2(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), util.GetAssembler().GetRegEntry($2), util.GetAssembler().GetRegEntry($4))
			}
		}
		|	F2NEMONICO numeros COMMA REG_X nlinea_op{
			if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true {
					if(util.GetAssembler().IsValidRegister($2)){
							util.GetAssembler().Assemblef2(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), $2, util.GetAssembler().GetRegEntry($4))
					}else{
						yylex.(*yylexer).Error("error de registro")
						pc -= 2
						goto yydefault
					}
			}
		}
		|	F2NEMONICO REG COMMA REG_X nlinea_op{
			if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true {
				util.GetAssembler().Assemblef2(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), util.GetAssembler().GetRegEntry($2), util.GetAssembler().GetRegEntry($4))
			}
		}

		;
		f3:
			simple3
		|	indirecto3
		|	inmediato3
		;
	f4:
			PLUS {format4 = true}f3 {format4 = false}
		;


	directiva:
			instr_byte
		|	instr_word
		|	instr_resb
		|	instr_resw
		|	instr_base
		;
	simple3:
			NEMONICO identificador indexado nlinea_op {
					if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == false	{
						util.GetAssembler().Assemble(yylex.(*yylexer).tlineidx,util.GetAssembler().GetCodeEntry($1), $3 , util.GetTabSim().GetSymbolAddrs($2))
					} else{
							if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true	{
								if format4 == false{
									util.GetAssembler().AssembleXE(yylex.(*yylexer).tlineidx, $1, $3, $2, bRegister, 3)
								}else{
									util.GetAssembler().Assemblef4(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), 3, $3, 0, 0, 1, util.GetTabSim().GetSymbolAddrs($2))

								}
						}
				}
		}
		| NEMONICO numeros indexado nlinea_op {
					if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == false	{
						//*¿Error?*/
						//util.GetAssembler().Assemble(yylex.(*yylexer).tlineidx,util.GetAssembler().GetCodeEntry($1), $3 , util.GetTabSim().GetSymbolAddrs($2))
					} else{
							if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true	{
								if format4 == false{
									if $2>=0 && $2<=4095 {
										util.GetAssembler().Assemblef3(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), 3, $3, 0, 0, 0, $2)
									}
							}else {
								if $2 >= 4095  {
									util.GetAssembler().Assemblef4(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), 3, $3, 0, 0, 1, $2)

								}
							}
						}
				}
		}
		|	NEMONICOXE identificador indexado nlinea_op{
				if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true {
					if format4 == false{
						util.GetAssembler().AssembleXE(yylex.(*yylexer).tlineidx, $1, $3, $2, bRegister, 3)
					}else{
						util.GetAssembler().Assemblef4(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), 3, $3, 0, 0, 1, util.GetTabSim().GetSymbolAddrs($2))
					}
				}	else{
					if yylex.(*yylexer).isXE == false	{
						yylex.(*yylexer).CompatibilityErr("de formato 3")
						goto yydefault
					}
				}
			}
			|	NEMONICOXE numeros indexado nlinea_op{
					if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true {
							if format4 == false{
								if $2>=0 && $2<=4095 {
									util.GetAssembler().Assemblef3(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), 3, $3, 0, 0, 0, $2)
								}
						}else {
							if $2 >= 4095  {
								util.GetAssembler().Assemblef4(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), 3, $3, 0, 0, 1, $2)
							}
						}
					}	else{
						if yylex.(*yylexer).isXE == false	{
							yylex.(*yylexer).CompatibilityErr("de formato 3")
							goto yydefault
						}
					}
				}
		|	RSUB nlinea_op	{
			if yylex.(*yylexer).firstparse == false	&& yylex.(*yylexer).isXE == false {
				util.GetAssembler().Assemble(yylex.(*yylexer).tlineidx,util.GetAssembler().GetCodeEntry($1), 0, 0)
			}else{
				if yylex.(*yylexer).firstparse == false	&& yylex.(*yylexer).isXE == true{
						if format4 == false{
							util.GetAssembler().Assemblef3(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1),3,0,0,0,0,0)
						}else{
							util.GetAssembler().Assemblef4(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1),3,0,0,0,0,0)
						}

					}
			}
			}
		;
	indirecto3:
			NEMONICO AT identificador nlinea_op{
				if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true	{
					if format4 == false{
						util.GetAssembler().AssembleXE(yylex.(*yylexer).tlineidx, $1, 0, $3, bRegister, 2)
					}else{
						util.GetAssembler().Assemblef4(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), 2, 0, 0, 0, 1, util.GetTabSim().GetSymbolAddrs($3))
					}
			}
		}
		|	NEMONICOXE AT identificador nlinea_op{
			if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true	{
				if format4 == false{
					util.GetAssembler().AssembleXE(yylex.(*yylexer).tlineidx, $1, 0, $3, bRegister, 2)
				}else{
					util.GetAssembler().Assemblef4(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), 2, 0, 0, 0, 1, util.GetTabSim().GetSymbolAddrs($3))
				}
		}
		}
		|	NEMONICO AT numeros nlinea_op{
			if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true {
					if format4 == false{
						if $3>=0 && $3<=4095 {
							util.GetAssembler().Assemblef3(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), 2, 0, 0, 0, 0, $3)
						}
				}else {
					if $3 >= 4095 {
						util.GetAssembler().Assemblef4(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), 2, 0, 0, 0, 1, $3)
					}
				}
			}
		}
		|	NEMONICOXE AT numeros nlinea_op{
			if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true {
					if format4 == false{
						if $3>=0 && $3<=4095 {
							util.GetAssembler().Assemblef3(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), 2, 0, 0, 0, 0, $3)
						}
				}else {
					if $3 >= 4095 {
						util.GetAssembler().Assemblef4(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), 2, 0, 0, 0, 1, $3)
					}
				}
			}
		}
		;
	inmediato3:
			NEMONICO HASH identificador nlinea_op{
				if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true	{
					if format4 == false{
						util.GetAssembler().AssembleXE(yylex.(*yylexer).tlineidx, $1, 0, $3, bRegister, 1)
					}else{
						util.GetAssembler().Assemblef4(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), 1, 0, 0, 0, 1, util.GetTabSim().GetSymbolAddrs($3))
					}
			}
			}
		|	NEMONICOXE HASH identificador nlinea_op{
			if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true	{
				if format4 == false{
					util.GetAssembler().AssembleXE(yylex.(*yylexer).tlineidx, $1, 0, $3, bRegister, 1)
				}else{
					util.GetAssembler().Assemblef4(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), 1, 0, 0, 0, 1, util.GetTabSim().GetSymbolAddrs($3))
				}
		}
		}
		|	NEMONICO HASH numeros nlinea_op{
			if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true {

					if format4 == false{
						if $3>=0 && $3<=4095 {
							util.GetAssembler().Assemblef3(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), 1, 0, 0, 0, 0, $3)
						}
				}else {
					if $3 >= 4095   {
						util.GetAssembler().Assemblef4(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), 1, 0, 0, 0, 1, $3)
					}
				}
			}
		}
		|	NEMONICOXE HASH numeros nlinea_op{
			if yylex.(*yylexer).firstparse == false && yylex.(*yylexer).isXE == true {

					if format4 == false{
						if $3>=0 && $3<=4095 {
							util.GetAssembler().Assemblef3(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), 1, 0, 0, 0, 0, $3)
						}
				}else {
					if $3 >= 4095 {
						util.GetAssembler().Assemblef4(yylex.(*yylexer).tlineidx, util.GetAssembler().GetCodeEntry($1), 1, 0, 0, 0, 1, $3)
					}
				}
			}
		}
		;
	instr_base:
			BASE identificador nlinea_op	{
						if yylex.(*yylexer).isXE == false	{
							yylex.(*yylexer).CompatibilityErr("BASE")
							goto yydefault
						}else{
							bRegister =  util.GetTabSim().GetSymbolAddrs($2)
							util.GetAssembler().SetNoCodeFlag(yylex.(*yylexer).tlineidx)
						}
					}
		;
	instr_byte:
			BYTE constantes nlinea_op {
				pc = pc + $2.tam
				if yylex.(*yylexer).firstparse == false	{
						if yylex.(*yylexer).isXE == false{
							if $2.hex == true{
								util.GetAssembler().AssembleByteHex(yylex.(*yylexer).tlineidx,$2.ascii)
							}else{
								util.GetAssembler().AssembleByteASCII(yylex.(*yylexer).tlineidx,$2.ascii)
							}
							$2.ascii = nil
					}else{
						var s string
						if $2.hex == true{
							s = util.GetAssembler().AssembleByteHex(yylex.(*yylexer).tlineidx,$2.ascii)
						}else{
							s = util.GetAssembler().AssembleByteASCII(yylex.(*yylexer).tlineidx,$2.ascii)
						}
						util.GetAssembler().Hexcode = append(util.GetAssembler().Hexcode, s)
					}
				}
			}
		;
	instr_word:
			WORD numerosDir nlinea_op {
			pc = pc + 3
				if yylex.(*yylexer).firstparse == false	{
					if yylex.(*yylexer).isXE == false{
						util.GetAssembler().AssembleWord(yylex.(*yylexer).tlineidx,$2)
					}else{
						s := util.GetAssembler().AssembleWord(yylex.(*yylexer).tlineidx,$2)
						util.GetAssembler().Hexcode = append(util.GetAssembler().Hexcode, s)
					}
				}
			}
		;
	instr_resb:
			RESB numerosDir nlinea_op {
			pc = pc + $2
				if yylex.(*yylexer).firstparse == false	{
						if yylex.(*yylexer).isXE == false{
							util.GetAssembler().SetRegFlag(yylex.(*yylexer).tlineidx)
						}else{
							util.GetAssembler().SetCutRegFlag(yylex.(*yylexer).tlineidx)
						}
				}
			}
		;
	instr_resw:
			RESW numerosDir nlinea_op {
			pc = pc + $2*3
				if yylex.(*yylexer).isXE == false{
					util.GetAssembler().SetRegFlag(yylex.(*yylexer).tlineidx)
				}else{
					util.GetAssembler().SetCutRegFlag(yylex.(*yylexer).tlineidx)
				}
			}
		;
	etiqueta_op:
			/*epsilon*/
		|	IDENTIFIER  {
										if yylex.(*yylexer).firstparse == true	{
												util.GetTabSim().AddSymbol($1,pc)
										}
									}
		;
	nlinea_op:
			/*epsilon*/
		|	NL nlinea_op
		;
	nlinea_nop:
			NL
		|	NL nlinea_nop
		;
	indexado:
			/*epsilon*/ {	$$ = 0 }
		|	COMMA REG_X { $$ = 1 }
		;
	numIdent:
			identificador { $$ = $1 }
		|	numeros				{ }
		;
	/*numConst:
			numeros
		|	constantes
		;*/
	constantes:
			CHARCONST
		|	HEXCONST
		;
	numeros:
			NUMBER {

						 }
		|	HEXNUMBER {

						 		}
		;
	numerosDir:
			NUMBER //{($1)}
		|	HEXNUMBER //{($1)}
		;
	identificador:
			IDENTIFIER
		;
%%
