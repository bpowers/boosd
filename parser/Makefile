parse.go: parse.y
	go tool yacc -o $@ -p boosd $< && gofmt -w $@ && sed -i 's/package main/package parser/' $@
