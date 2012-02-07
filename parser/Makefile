parse.go: parse.y
	go tool yacc -o $@ -p boosd $< && gofmt -w $@
