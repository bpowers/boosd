// Copyright 2013 Bobby Powers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/bpowers/boosd/boosd"
	"github.com/davecheney/gogo"
	"github.com/davecheney/gogo/build"
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
)

const usage = `Usage: %s [OPTION...]
Compile system dynamics models.

Options:
`

var (
	outPath string
)

func gofmtFile(f *ast.File, goFset *token.FileSet) ([]byte, error) {
	var buf bytes.Buffer
	if err := format.Node(&buf, goFset, f); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&outPath, "o", "model.out",
		"file name to use as output")

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
		defer f.Close()
		fi = bufio.NewReader(f)
	}

	// dump in the file
	mdl, err := ioutil.ReadAll(fi)
	if err != nil {
		log.Fatal("ReadAll:", err)
	}

	file := fset.AddFile(filename, fset.Base(), len(mdl))

	// and parse
	pkg, err := boosd.Parse(file, string(mdl))
	if err != nil {
		log.Fatalf("Parse(%v): %s", file, err)
	}
	if pkg.NErrors > 0 {
		log.Fatal("There were errors parsing the file")
	}
	// log.Printf("compilationUnit: %#v\n", f)
	boosd.PassScopeChain(pkg)

	/*
		mainMdl := pkg.GetModel("main")

		if mainMdl == nil {
			log.Fatal("No main model")
		} else if mainMdl.Virtual {
			log.Fatal("Main model can't have undefined variables")
		}
	*/

	goFile, err := boosd.GenGo(pkg)
	if err != nil {
		log.Fatalf("generateGoAST(%v): %s", pkg, err)
	}

	src, err := gofmtFile(goFile, goFset)
	if err != nil {
		log.Fatalf("gofmtFile(%v): %s", goFile, err)
	}

	workDir, err := ioutil.TempDir("", "boost_temp")
	if err != nil {
		log.Fatalf("ioutil.TempDir('', 'boost_temp'): %s", err)
	}

	if err = os.Mkdir(path.Join(workDir, ".gogo"), os.ModeDir|0700); err != nil {
		log.Fatalf("os.Mkdir(%s, mode|0700): %s", path.Join(workDir, ".gogo"), err)
	}

	if err = os.Mkdir(path.Join(workDir, "src"), os.ModeDir|0700); err != nil {
		log.Fatalf("os.Mkdir(%s, mode|0700): %s", path.Join(workDir, "src"), err)
	}

	proj := path.Join(workDir, "src", path.Base(outPath))
	if err = os.Mkdir(proj, os.ModeDir|0700); err != nil {
		log.Fatalf("os.Mkdir(%s, mode|0700): %s", proj, err)
	}

	log.Printf("proj: %s", proj)

	goProj, err := gogo.NewProject(workDir)
	if err != nil {
		log.Fatalf("NewProject(%s): %s", proj, err)
	}

	ctx, err := gogo.NewDefaultContext(goProj)
	if err != nil {
		log.Fatalf("NewDefaultContext: %s", err)
	}
	defer ctx.Destroy()

	if err = os.Symlink("/var/unsecure/src/github.com", path.Join(workDir, "src", "github.com")); err != nil {
		log.Fatalf("symlink: %s", err)
	}

	of, err := os.Create(path.Join(proj, "main.go"))
	if err != nil {
		log.Fatalf("Create: %s", err)
	}
	of.Write(src)
	of.Close()

	goPkg, err := ctx.ResolvePackage(path.Base(outPath))
	if err != nil {
		log.Fatalf("ResolvePackage: %s", err)
	}

	if err = build.Build(ctx, goPkg).Result(); err != nil {
		log.Fatalf("Build: %s", err)
	}

	if err = copyFile(path.Join(workDir, "bin", "linux", "amd64", path.Base(outPath)), path.Base(outPath)); err != nil {
		log.Fatalf("copyFile: %s", err)
	}
}

func copyFile(from, to string) error {
	fromF, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("Open(%s): %s", from, err)
	}
	defer fromF.Close()

	toF, err := os.OpenFile(to, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return fmt.Errorf("Open(%s): %s", to, err)
	}
	defer toF.Close()

	_, err = io.Copy(toF, fromF)
	return err
}

func mustGetwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}
