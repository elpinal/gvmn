package gvmn

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
)

// doubleError is a type which has two error.
type doubleError struct {
	a, b error
}

func (e *doubleError) Error() string {
	if e == nil {
		return ""
	}
	if e.a == nil {
		return e.b.Error()
	}
	if e.b == nil {
		return e.a.Error()
	}
	return fmt.Sprintf("%v\n%v", e.a, e.b)
}

// build builds the specified Go version. It accept GOROOT_BOOTSTRAP if it is
// set. Otherwise, it uses GOROOT as GOROOT_BOOTSTRAP if it is set.  If neither
// GOROOT_BOOTSTRAP nor GOROOT is set, the result of the execution of `go env
// GOROOT`.
func build(version string) *doubleError {
	env := os.Environ()
	if gorootBootstrap := os.Getenv("GOROOT_BOOTSTRAP"); gorootBootstrap != "" {
		// nothing to do
	} else if goroot := os.Getenv("GOROOT"); goroot != "" {
		env = append(env, "GOROOT_BOOTSTRAP="+goroot)
	} else if goroot, err := exec.Command("go", "env", "GOROOT").Output(); err == nil {
		env = append(env, "GOROOT_BOOTSTRAP="+string(bytes.TrimSuffix(goroot, []byte("\n"))))
	}
	cmd := exec.Command("./make.bash")
	cmd.Dir = filepath.Join(gvmnrootGo, version, "src")
	cmd.Env = env
	var buf bytes.Buffer
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return &doubleError{errors.Wrap(err, "./make.bash failed"), fmt.Errorf(buf.String())}
	}
	return nil
}

// Install installs Go version.
func Install(version string) error {
	if exist(filepath.Join(gvmnrootGo, version, "bin", "go")) {
		return nil
	}
	if err := build(version); err != nil {
		return err
	}
	return nil
}

// checkout checkouts the specified version of the Go repository.
func checkout(version string) *doubleError {
	versionsDir := filepath.Join(gvmnrootGo, version)
	cmd := exec.Command("git", "clone", gvmnrootRepo, versionsDir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return &doubleError{errors.Wrapf(err, "cloning (%s) from cached repository", version), fmt.Errorf("%s", out)}
	}

	cmd = exec.Command("git", "reset", "--hard", version)
	cmd.Dir = versionsDir
	out, err = cmd.CombinedOutput()
	if err != nil {
		_ = os.RemoveAll(versionsDir)
		return &doubleError{errors.Wrapf(err, "checking out (%s)", version), fmt.Errorf("%s", out)}
	}

	return nil
}

// update updates the Go repository.
func update() *doubleError {
	cmd := exec.Command("git", "fetch")
	cmd.Dir = gvmnrootRepo
	if out, err := cmd.CombinedOutput(); err != nil {
		return &doubleError{errors.Wrap(err, "fetching updates from remote repository"), fmt.Errorf("%s", out)}
	}
	return nil
}

// mirror mirrors the Go repository.
func mirror() *doubleError {
	out, err := exec.Command("git", "clone", "--mirror", RepoURL, gvmnrootRepo).CombinedOutput()
	if err != nil {
		return &doubleError{errors.Wrap(err, "cloning from remote repository"), fmt.Errorf("%s", out)}
	}
	return nil
}

// download fetches the Go repository.
func download() error {
	if !exist(gvmnrootRepo) {
		if err := mirror(); err != nil {
			return err
		}
	}
	if err := update(); err != nil {
		return err
	}
	return nil
}

// Download fetches the Go repository and check out version.
func Download(version string, update bool) error {
	if exist(filepath.Join(gvmnrootGo, version)) {
		if !update {
			return nil
		}
		os.RemoveAll(filepath.Join(gvmnrootGo, version))
	}
	if err := download(); err != nil {
		return err
	}
	if err := checkout(version); err != nil {
		return err
	}
	return nil
}

// Get downloads and installs Go.
func Get(version string) error {
	if err := Download(version, false); err != nil {
		return err
	}
	if err := Install(version); err != nil {
		return err
	}
	return nil
}

// writeFile writes files.
func writeFile(dest string, r io.Reader, mode os.FileMode) error {
	f, err := os.Create(dest)
	if err != nil {
		return errors.Wrap(err, "os.Create")
	}
	defer f.Close()
	if err := f.Chmod(mode); err != nil {
		return errors.Wrap(err, "Chmod")
	}
	if _, err := io.Copy(f, r); err != nil {
		return errors.Wrap(err, "io.Copy")
	}
	return nil
}

// unTarGz extract content from tar.gz.
func unTarGz(src io.Reader, version string) error {
	r, err := gzip.NewReader(src)
	if err != nil {
		return err
	}
	defer r.Close()
	rd := tar.NewReader(r)

	for {
		hdr, err := rd.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Trim prefix, "go".
		if len(hdr.Name) < 2 {
			return fmt.Errorf("failed to extract %v; too short name", hdr.Name)
		}
		path := filepath.Join(gvmnrootGo, version, hdr.Name[2:])

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, os.FileMode(hdr.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := writeFile(path, rd, hdr.FileInfo().Mode()); err != nil {
				return err
			}
		default:
			return fmt.Errorf("cannot: %c, %s", hdr.Typeflag, path)
		}
	}

	return nil
}

var script = []byte(`#!/bin/bash

path="$(dirname $0)"

GOROOT="$path" "$path"/go-org "$@"
`)

// GetBinary is like Get but gets binaries instead.
func GetBinary(version string) error {
	if exist(filepath.Join(gvmnrootGo, version, "bin", "go")) {
		return nil
	}

	goos := runtime.GOOS
	goarch := runtime.GOARCH
	suffix := "tar.gz"
	if goos == "windows" {
		suffix = "zip"
	}
	file := fmt.Sprintf("%s.%s-%s.%s", version, goos, goarch, suffix)
	u := url.URL{
		Scheme: "https",
		Host:   "storage.googleapis.com",
		Path:   path.Join("golang", url.PathEscape(file)),
	}
	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch goos {
	case "windows":
		// FIXME: unzip...
		return fmt.Errorf("Windows is not supported yet")
	default:
		if err := unTarGz(resp.Body, version); err != nil {
			return err
		}
	}

	goBin := filepath.Join(gvmnrootGo, version, "bin", "go")
	if err := os.Rename(goBin, goBin+"-org"); err != nil {
		return err
	}
	if err := ioutil.WriteFile(goBin, script, 0777); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

// latestTag reports the latest tag of the Go repository.
func latestTag() (string, error) {
	cmd := exec.Command("git", "rev-list", "--tags", "--max-count=1")
	cmd.Dir = gvmnrootRepo
	out, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "git rev-list failed")
	}
	sha := string(bytes.TrimSuffix(out, []byte("\n")))
	cmd = exec.Command("git", "describe", "--tags", sha)
	cmd.Dir = gvmnrootRepo
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	tag, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, stderr.String())
	}
	return string(bytes.TrimSuffix(tag, []byte("\n"))), nil
}

// LatestTag downloads the updated Go repository and reports
// the latest tag of it.
func LatestTag() (string, error) {
	if err := download(); err != nil {
		return "", err
	}
	return latestTag()
}
