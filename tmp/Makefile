include $(GOROOT)/src/Make.inc

TARG = cascade

GOFILES = \
	cascade.go\

include $(GOROOT)/src/Make.cmd

.PHONY: gofmt
gofmt:
	gofmt -w $(GOFILES)
