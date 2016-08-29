package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func mustRemoveAll(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	flag.Parse()

	var err error
	GvmnDir, err = ioutil.TempDir("", "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "TempDir: %v", err)
		os.Exit(2)
	}

	if !exist("testdata/go") {
		if testing.Verbose() {
			log.SetFlags(log.Lshortfile)
			log.Print("fetching for test...")
		}
		out, err := exec.Command("git", "clone", "--depth=1", "--bare", "--branch=go1.7", RepoURL, "testdata/go").CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to fetch for test: %v\n", err)
			fmt.Fprintln(os.Stderr, string(out))
			os.Exit(2)
		}
	}
	RepoURL = "testdata/go"

	r := m.Run()

	mustRemoveAll(GvmnDir)
	os.Exit(r)
}

func TestCmdInstall(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	if got := runInstall([]string{"go1.7"}); got != 0 {
		t.Fatalf("got %v, want 0", got)
	}
	bin := filepath.Join(GvmnDir, "versions/go1.7/bin/go")
	out, err := exec.Command(bin, "version").Output()
	if err != nil {
		t.Fatalf("go version: %v", err)
	}
	if !bytes.Contains(out, []byte("1.7")) {
		t.Fatalf("%v must contain 1.7", out)
	}
}

func TestInstallUntilBuild(t *testing.T) {
	if !testing.Short() {
		t.Skip("skipping in non-short mode")
	}
	if err := download(); err != nil {
		t.Fatalf("download() failed: %v", err)
	}
	if err := checkout("go1.7"); err != nil {
		t.Fatalf(`checkout("go1.7") failed: %v`, err)
	}
	if err := writeVersion("go1.7"); err != nil {
		t.Fatalf(`writeVersion("go1.7") failed: %v`, err)
	}
}
