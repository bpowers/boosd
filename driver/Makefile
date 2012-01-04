include $(GOROOT)/src/Make.inc

TARG = boosd

GOFILES = \
	main.go \
	pass_*.go \

include $(GOROOT)/src/Make.cmd

.PHONY: gofmt
gofmt:
	gofmt -w $(GOFILES)

