include $(GOROOT)/src/Make.inc

TARG = boosd

GOFILES = \
	boosd.go

include $(GOROOT)/src/Make.cmd

.PHONY: gofmt
gofmt:
	gofmt -w $(GOFILES)

