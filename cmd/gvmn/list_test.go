package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/elpinal/gvmn"
)

func TestMain(m *testing.M) {
	dir, err := ioutil.TempDir("", "gvmn")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	gvmn.SetRoot(dir)

	subdir := filepath.Join(dir, "go/go1.7/bin")
	err = os.MkdirAll(subdir, os.ModePerm)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	f, err := os.Create(filepath.Join(subdir, "go"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = f.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	code := m.Run()

	_ = os.RemoveAll(dir)

	os.Exit(code)
}

func TestList(t *testing.T) {
	var buf bytes.Buffer
	l := lister{out: &buf, err: &buf}
	code := l.listMain()
	s := buf.String()
	if code != 0 {
		t.Fatalf("listMain = %d; output is: %s", code, s)
	}
	if !strings.Contains(s, "Installed:") {
		t.Fatalf("listMain: output does not contain %q in %q", "Installed:", s)
	}
}
