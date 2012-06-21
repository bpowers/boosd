package main

import (
	"boosd/parser"
	"boosd/token"
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	var fset *token.FileSet = token.NewFileSet()
	var filename string
	var fi *bufio.Reader
	var f *os.File
	var err error

	flag.Parse()

	// use the file if there is an argument, otherwise use stdin
	if flag.NArg() == 0 {
		filename = "stdin"
		fi = bufio.NewReader(os.NewFile(0, "stdin"))
	} else {
		filename = flag.Arg(0)
		f, err = os.Open(filename)
		if err != nil {
			log.Fatal("Open:", err)
		}
		fi = bufio.NewReader(f)
	}

	// dump in the file
	mdl, err := ioutil.ReadAll(fi)
	if err != nil {
		log.Fatal("ReadAll:", err)
	}

	if f != nil && f.Close() != nil {
		log.Fatal("f.Close()")
	}

	file := fset.AddFile(filename, fset.Base(), len(mdl))

	// and parse
	pkg := parser.Parse(file, string(mdl))
	if pkg.NErrors > 0 {
		log.Fatal("There were errors parsing the file")
	}
	// log.Printf("compilationUnit: %#v\n", f)
	passTimespec(pkg)
	passScopeChain(pkg)
	passTypeResolution(pkg)

	mainMdl := pkg.GetModel("main")

	if mainMdl == nil {
		log.Fatal("No main model")
	} else if mainMdl.Virtual {
		log.Fatal("Main model can't have undefined variables")
	}
}
