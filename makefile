Executable=main
parserSrc=sicgrammar.go
lexerSrc=siclexer.go
sources=main.go
parserDir=src/parser/
binDir=bin/
srcDir=src/
REPORTE=reporte6

LexSources=$(addprefix $(parserDir),$(lexerSrc))
ParserSources=$(addprefix $(parserDir),$(parserSrc))
GoExecutable=$(addprefix $(binDir),$(Executable))


all:	dparser $(GoExecutable)

dparser:	
	cd src/parser && $(MAKE)

$(GoExecutable): $(sources)
		go build $^ 
practica:
	./main $(EJEMPLOS)/*.[sx]

reporte: $(REPORTE).tex
	latex $(REPORTE).tex
	dvipdf $(REPORTE).dvi
	evince $(REPORTE).pdf&
compara:
	diff copy.os $(EJEMPLOS)/copy.os
	diff ejemplo1.os $(EJEMPLOS)/ejemplo1.os
	diff ejemplo2.os $(EJEMPLOS)/ejemplo2.os
	diff ejemplo3.os $(EJEMPLOS)/ejemplo3.os
	diff copy.ox $(EJEMPLOS)/copy.ox
	diff ejemplo1.ox $(EJEMPLOS)/ejemplo1.ox
	diff ejercicio.ox $(EJEMPLOS)/ejercicio.ox
	diff exa2.ox $(EJEMPLOS)/exa2.ox
clean:
	rm -rf $(Executable)
	rm -rf $(LexSources)
	rm -rf $(ParserSources)
	rm -rf $(REPORTE).aux 
	rm -rf $(REPORTE).log 
	rm -rf $(REPORTE).dvi 
	rm -rf $(REPORTE).out
	rm *.os
	rm *.tx
	rm *.ts
	rm *.ox
