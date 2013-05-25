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

const (
	// FIXME: this is $GOPATH/src/github.com - non-hardcode this
	GithubDir = "/var/unsecure/src/github.com"
)

var (
	outPath string
)

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
	var filename string
	var in *bufio.Reader
	var err error

	// use the file if there is an argument, otherwise use stdin
	if flag.NArg() == 0 {
		filename = "stdin"
		in = bufio.NewReader(os.NewFile(0, "stdin"))
	} else {
		filename = flag.Arg(0)
		f, err := os.Open(filename)
		if err != nil {
			log.Fatal("Open:", err)
		}
		defer f.Close()
		in = bufio.NewReader(f)
	}

	goSource, err := transliterate(filename, in)
	if err != nil {
		log.Fatalf("%s", err)
	}

	exePath, err := compileAndLink(goSource)
	if err != nil {
		log.Fatalf("compileAndLink('%s')", goSource, err)
	}

	if err = copyFile(exePath, outPath); err != nil {
		log.Fatalf("copyFile('%s', '%s'): %s", exePath, outPath, err)
	}
}

// copyFile copies the file at path 'from' to path 'to', overwriting
// the file at 'to' if it already exists.
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

// mustGetwd returns the current working directory, panicing on error.
func mustGetwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}

// gofmt takes the given, valid, Go AST and returns a
// canonically-formatted go program in a byte-array, or an error.
func gofmt(f *ast.File) ([]byte, error) {
	fset := token.NewFileSet()
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// transliterate takes an input stream and a name and returns a byte
// buffer containing valid & gofmt'ed source code, or an error.  The
// name is used purely for diagnostic purposes
func transliterate(name string, in io.Reader) ([]byte, error) {
	fset := token.NewFileSet()

	// dump in the file
	mdlSrc, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, fmt.Errorf("ReadAll(%v): %s", in, err)
	}

	fsetFile := fset.AddFile(name, fset.Base(), len(mdlSrc))

	// and parse
	pkg, err := boosd.Parse(fsetFile, string(mdlSrc))
	if err != nil {
		return nil, fmt.Errorf("Parse(%v): %s", name, err)
	}
	if pkg.NErrors > 0 {
		return nil, fmt.Errorf("There were errors parsing the file")
	}

	goSource, err := boosd.GenGo(pkg)
	if err != nil {
		log.Fatalf("GenGo(%v): %s", pkg, err)
	}

	src, err := gofmt(goSource)
	if err != nil {
		log.Fatalf("gofmtFile(%v): %s", goSource, err)
	}
	return src, nil
}

// mkdir joins the components into a path and creates that path with
// the given octal permissions.  Parent directories must already
// exist.
func mkdir(perm os.FileMode, components ...string) error {
	p := path.Join(components...)
	if err := os.Mkdir(p, perm); err != nil {
		return fmt.Errorf("os.Mkdir(%s, %d): %s", p, perm, err)
	}
	return nil
}

// compileAndLink compiles and links the given source, returning a
// path to the given binary, or an error.
func compileAndLink(src []byte) (string, error) {
	workDir, err := ioutil.TempDir("", "boost_temp")
	if err != nil {
		return "", fmt.Errorf("ioutil.TempDir: %s", err)
	}

	if err = mkdir(0700, workDir, ".gogo"); err != nil {
		return "", err
	}

	if err = mkdir(0700, workDir, "src"); err != nil {
		return "", err
	}

	if err = mkdir(0700, workDir, "src", "model.out"); err != nil {
		return "", err
	}

	proj, err := gogo.NewProject(workDir)
	if err != nil {
		return "", fmt.Errorf("NewProject(%s): %s", workDir, err)
	}

	ctx, err := gogo.NewDefaultContext(proj)
	if err != nil {
		return "", fmt.Errorf("NewDefaultContext(): %s", err)
	}
	defer ctx.Destroy()

	if err = os.Symlink(GithubDir, path.Join(workDir, "src", "github.com")); err != nil {
		return "", fmt.Errorf("symlink: %s", err)
	}

	srcPath := path.Join(workDir, "src", "model.out", "main.go")
	f, err := os.Create(srcPath)
	if err != nil {
		return "", fmt.Errorf("Create(%s): %s", srcPath, err)
	}
	f.Write(src)
	// this Close is not deferred so that we're sure the contents are flushed to
	// the kernel before buid.Build is called.
	f.Close()

	goPkg, err := ctx.ResolvePackage("model.out")
	if err != nil {
		return "", fmt.Errorf("ResolvePackage(model.out): %s", err)
	}

	if err = build.Build(ctx, goPkg).Result(); err != nil {
		return "", fmt.Errorf("Build: %s", err)
	}

	return path.Join(workDir, "bin", "linux", "amd64", "model.out"), nil
}
