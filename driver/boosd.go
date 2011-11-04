package main

import (
	"boosd/parser"
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	flag.Parse()

	var fi *bufio.Reader
	// use the file if there is an argument, otherwise use stdin
	if flag.NArg() == 0 {
		fi = bufio.NewReader(os.NewFile(0, "stdin"))
	} else {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal("Open:", err)
		}
		fi = bufio.NewReader(f)
	}

	// dump in the string
	units, err := ioutil.ReadAll(fi)
	if err != nil {
		log.Fatal("ReadAll:", err)
	}

	// and parse
	f := parser.Parse(string(units))
	log.Printf("compilationUnit: %#v\n", f)
	indent := ""
	parser.Inspect(f, func(n parser.Node) bool {
		if n == nil {
			indent = indent[:len(indent) - 2]
		} else {
			indent += "  "
			fmt.Printf("%s%#v\n", indent, n)
		}
		return true
	})
}
