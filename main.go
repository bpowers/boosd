package main

import (
	"boosd/parser"
	"bufio"
	"bytes"
	"flag"
	"go/ast"
	"go/format"
	"go/token"
	"io/ioutil"
	"log"
	"os"
)

func gofmtFile(f *ast.File, goFset *token.FileSet) ([]byte, error) {
	var buf bytes.Buffer
	if err := format.Node(&buf, goFset, f); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func init() {
	flag.Parse()
}

func main() {
	var fset, goFset *token.FileSet = token.NewFileSet(), token.NewFileSet()
	var filename string
	var fi *bufio.Reader
	var f *os.File
	var err error

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
	passScopeChain(pkg)

	/*
		mainMdl := pkg.GetModel("main")

		if mainMdl == nil {
			log.Fatal("No main model")
		} else if mainMdl.Virtual {
			log.Fatal("Main model can't have undefined variables")
		}
	*/

	goFile, err := genGoFile(pkg)
	if err != nil {
		log.Fatal("genGoAST(%v): %s", pkg, err)
	}

	src, err := gofmtFile(goFile, goFset)
	if err != nil {
		log.Fatal("gofmtFile(%v): %s", goFile, err)
	}

	os.Stdout.Write(src)
}
