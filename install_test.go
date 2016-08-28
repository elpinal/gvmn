package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
	var err error
	GvmnDir, err = ioutil.TempDir("", "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "TempDir: %v", err)
		os.Exit(2)
	}

	r := m.Run()

	mustRemoveAll(GvmnDir)
	os.Exit(r)
}

func TestCmdInstall(t *testing.T) {
	t.Log("GvmnDir for test:", GvmnDir)
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
