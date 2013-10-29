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
	"github.com/davecheney/gogo/build"
	"github.com/davecheney/gogo/project"
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
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

		if fi, err := f.Stat(); err != nil || fi.IsDir() {
			log.Fatalf("err(%s) or %s IsDir", err, filename)
		}

		in = bufio.NewReader(f)
	}

	goSource, err := transliterate(filename, in)
	if err != nil {
		log.Fatalf("%s", err)
	}

	err = compileAndLink(goSource, outPath)
	if err != nil {
		log.Fatalf("compileAndLink('%s')", goSource, err)
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
func compileAndLink(src []byte, outPath string) error {
	workDir, err := ioutil.TempDir("", "boosd_temp")
	if err != nil {
		return fmt.Errorf("ioutil.TempDir: %s", err)
	}

	if err = mkdir(0700, workDir, ".gogo"); err != nil {
		return err
	}

	if err = mkdir(0700, workDir, "bin"); err != nil {
		return err
	}

	if err = mkdir(0700, workDir, "src"); err != nil {
		return err
	}

	if err = mkdir(0700, workDir, "src", "main"); err != nil {
		return err
	}

	proj, err := project.NewProject(workDir)
	if err != nil {
		return fmt.Errorf("NewProject(%s): %s", workDir, err)
	}

	ctx, err := build.NewDefaultContext(proj)
	if err != nil {
		return fmt.Errorf("NewDefaultContext(): %s", err)
	}
	defer ctx.Destroy()

	if err = os.Symlink(GithubDir, path.Join(workDir, "src", "github.com")); err != nil {
		return fmt.Errorf("symlink: %s", err)
	}

	srcPath := path.Join(workDir, "src", "main", "main.go")
	f, err := os.Create(srcPath)
	if err != nil {
		return fmt.Errorf("Create(%s): %s", srcPath, err)
	}
	f.Write(src)
	// this Close is not deferred so that we're sure the contents are flushed to
	// the kernel before buid.Build is called.
	f.Close()

	goPkg, err := ctx.ResolvePackage(runtime.GOOS, runtime.GOARCH, "main").Result()
	if err != nil {
		return fmt.Errorf("ResolvePackage(main): %s", err)
	}

	if err = build.Build(ctx, goPkg).Result(); err != nil {
		return fmt.Errorf("Build: %s", err)
	}

	exePath := path.Join(ctx.Bindir(), "main")
	if err = copyFile(exePath, outPath); err != nil {
		return fmt.Errorf("copyFile('%s', '%s'): %s", exePath, outPath, err)
	}
	return nil
}
