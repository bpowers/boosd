include $(GOROOT)/src/Make.inc

TARG = boosd

GOFILES = \
	lex.go\
	parse.go\

CLEANFILES += \
	parse.go\
	y.output\

include $(GOROOT)/src/Make.cmd

.PHONY: gofmt
gofmt:
	gofmt -w $(GOFILES)

parse.go: calc.y
	goyacc -o $@ -p boosd $< && gofmt -w $@
