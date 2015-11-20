all: sicgrammar.go siclexer.go

sicgrammar.go: sicgrammar.y
	go tool yacc -o sicgrammar.go sicgrammar.y
siclexer.go: siclexer.l
	golex -t $< | gofmt > $@

