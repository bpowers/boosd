package main

import (
	"boosd/parser"
	"bufio"
	"flag"
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
}
