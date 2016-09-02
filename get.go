package gvmn

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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

// build builds the specified Go version.
func build(version string) *doubleError {
	var env []string
	if goroot, err := exec.Command("go", "env", "GOROOT").Output(); err == nil {
		env = append(os.Environ(), "GOROOT_BOOTSTRAP="+string(bytes.TrimSuffix(goroot, []byte("\n"))))
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
		return &doubleError{errors.Wrap(err, "git clone failed"), fmt.Errorf("%s", out)}
	}

	cmd = exec.Command("git", "reset", "--hard", version)
	cmd.Dir = versionsDir
	out, err = cmd.CombinedOutput()
	if err != nil {
		return &doubleError{errors.Wrap(err, "git reset failed"), fmt.Errorf("%s", out)}
	}
	return nil
}

// update updates the go repository.
func update() *doubleError {
	cmd := exec.Command("git", "fetch")
	cmd.Dir = gvmnrootRepo
	if out, err := cmd.CombinedOutput(); err != nil {
		return &doubleError{errors.Wrap(err, "failed to fetch"), fmt.Errorf("%s", out)}
	}
	return nil
}

// mirror mirrors the go repository.
func mirror() *doubleError {
	out, err := exec.Command("git", "clone", "--mirror", RepoURL, gvmnrootRepo).CombinedOutput()
	if err != nil {
		return &doubleError{errors.Wrap(err, "cloning repository failed"), fmt.Errorf("%s", out)}
	}
	return nil
}

// download fetches the go repository.
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
func Download(version string) error {
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
	if err := Download(version); err != nil {
		return err
	}
	if err := Install(version); err != nil {
		return err
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
