package main

import (
	"boosd/parser"
	"boosd/token"
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	flag.Parse()

	var fs *token.FileSet = token.NewFileSet()
	var filename string
	var fi *bufio.Reader
	// use the file if there is an argument, otherwise use stdin
	if flag.NArg() == 0 {
		filename = "stdin"
		fi = bufio.NewReader(os.NewFile(0, "stdin"))
	} else {
		filename = flag.Arg(0)
		f, err := os.Open(filename)
		if err != nil {
			log.Fatal("Open:", err)
		}
		fi = bufio.NewReader(f)
	}

	// dump in the string
	mdl, err := ioutil.ReadAll(fi)
	if err != nil {
		log.Fatal("ReadAll:", err)
	}
	file := fs.AddFile(filename, fs.Base(), len(mdl))

	// and parse
	f := parser.Parse(file, string(mdl))
	//log.Printf("compilationUnit: %#v\n", f)
	//indent := ""
	parser.Inspect(f, func(node parser.Node) bool {
		if node == nil {
			//indent = indent[:len(indent) - 2]
		} else {
			switch n := node.(type) {
			case *parser.ModelDecl:
				fmt.Println("model", n.Name.Name)
			}
			//fmt.Printf("%s%#v\n", indent, node)
		}
		return true
	})
}
