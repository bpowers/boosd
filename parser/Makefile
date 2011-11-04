include $(GOROOT)/src/Make.inc

TARG = boosd/parser

GOFILES = \
	lex.go\
	parse.go\
	ast.go\
	scope.go\
	resolve.go\
	errors.go\

CLEANFILES += \
	parse.go\
	y.output\

include $(GOROOT)/src/Make.pkg

.PHONY: gofmt
gofmt:
	gofmt -w $(GOFILES)

parse.go: parse.y
	goyacc -o $@ -p boosd $< && gofmt -w $@ && sed -i 's/package main/package parser/' $@
