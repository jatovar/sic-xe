FIBO	START	2000H
	LDX		ZERO
 	LDA     ONE
	STCH	RT,X
	TIX     LENGTH
	COMP	LENGTH
	JGT		FIN
	JEQ		FIN
FOR	LDA		RI
	ADD		RJ
	STCH	RT,X
	LDA		RJ
	STA		RI
	LDCH	RT,X
	STA		RJ
	TIX		LENGTH
	JLT		FOR
	J		FIN
LENGTH	WORD	10
RI		WORD	0
RJ		WORD	1
ONE		WORD	1
ZERO	WORD	0
RT		RESB	30
FIN		BYTE	C'EOF'	
		END
