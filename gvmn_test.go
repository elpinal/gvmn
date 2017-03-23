package gvmn

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
	SetRoot(tempdir)
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

func TestGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	if err := Get("go1.7"); err != nil {
		t.Fatalf("Get: %v", err)
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

func TestDownload(t *testing.T) {
	if !testing.Short() {
		t.Skip("skipping in non-short mode")
	}
	if err := Download("go1.7", false); err != nil {
		t.Fatalf(`Download("go1.7", false) failed: %v`, err)
	}
}

func TestGetBinary(t *testing.T) {
	if err := GetBinary("go1.7"); err != nil {
		t.Fatalf("GetBinary: %v", err)
	}
	if list := List(); len(list) == 0 {
		t.Error("could not list Go versions")
	} else if list[0].Name != "go1.7" {
		t.Error("could not find go1.7")
	}
	if err := Run("go1.7", "version"); err != nil {
		t.Fatalf("go version: %v", err)
	}
}

func TestLatestTag(t *testing.T) {
	if out, err := LatestTag(); err != nil {
		t.Fatalf("LatestTag: %v", err)
	} else if out != "go1.7" {
		t.Fatalf("got %v, want go1.7", out)
	}
}

func TestUse(t *testing.T) {
	if err := Use("go1.7"); err != nil {
		t.Fatalf("Use: %v", err)
	}
}

func TestRemove(t *testing.T) {
	if err := Remove("go1.7"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
}
