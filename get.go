package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

var cmdGet = &Command{
	Run:       runGet,
	UsageLine: "get version",
	Short:     "download and install Go",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdGet.Flag.BoolVar(&flagA, "a", false, "")
}

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

// build builds the specified version of Go.
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

// checkout checkouts specified version of the Go repository.
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

// get downloads and installs Go.
func get(version string) error {
	if version == "latest" {
		var err error
		version, err = latestTag()
		if err != nil {
			return errors.Wrap(err, "failed to get the latest version")
		}
	}

	if err := checkout(version); err != nil {
		return err
	}

	if err := build(version); err != nil {
		return err
	}
	return nil
}

// download fetches repository from RepoURL.
func download() *doubleError {
	dir := gvmnrootRepo
	if !exist(dir) {
		out, err := exec.Command("git", "clone", "--mirror", RepoURL, dir).CombinedOutput()
		if err != nil {
			return &doubleError{errors.Wrap(err, "cloning repository failed"), fmt.Errorf("%s", out)}
		}
	}

	cmd := exec.Command("git", "fetch")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return &doubleError{errors.Wrap(err, "failed to fetch"), fmt.Errorf("%s", out)}
	}
	return nil
}

// runGet executes get command and return exit code.
func runGet(args []string) int {
	if len(args) == 0 {
		log.Print("gvmn get: no Go version specified")
		return 1
	}

	if err := download(); err != nil {
		log.Print(err)
		return 1
	}

	if err := get(args[0]); err != nil {
		log.Print(err)
		return 1
	}
	return 0
}
