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

	tempdir, err := ioutil.TempDir("", "")
	setGvmnroot(tempdir)
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

	mustRemoveAll(tempdir)
	os.Exit(r)
}

func TestCmdGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	if got := runGet([]string{"go1.7"}); got != 0 {
		t.Fatalf("got %v, want 0", got)
	}
	bin := filepath.Join(gvmnrootGo, "go1.7", "bin", "go")
	out, err := exec.Command(bin, "version").Output()
	if err != nil {
		t.Fatalf("go version: %v", err)
	}
	if !bytes.Contains(out, []byte("1.7")) {
		t.Fatalf("%v must contain 1.7", out)
	}
}

func TestGetUntilBuild(t *testing.T) {
	if !testing.Short() {
		t.Skip("skipping in non-short mode")
	}
	if err := download(); err != nil {
		t.Fatalf("download() failed: %v", err)
	}
	if err := checkout("go1.7"); err != nil {
		t.Fatalf(`checkout("go1.7") failed: %v`, err)
	}
}

func BenchmarkCheckout(b *testing.B) {
	if err := download(); err != nil {
		b.Fatalf("download() failed: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := checkout("go1.7"); err != nil {
			b.Fatalf(`checkout("go1.7") failed: %v`, err)
		}
		b.StopTimer()
		mustRemoveAll(gvmnrootGo)
		b.StartTimer()
	}
}
