package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

var cmdInstall = &Command{
	Run:       runInstall,
	UsageLine: "install ",
	Short:     "Install Go",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdInstall.Flag.BoolVar(&flagA, "a", false, "")
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
	cmd.Dir = filepath.Join(GvmnDir, "versions", version, "src")
	cmd.Env = env
	var buf bytes.Buffer
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return &doubleError{errors.Wrap(err, "./make.bash failed"), fmt.Errorf(buf.String())}
	}
	return nil
}

// writeVersion writes the version to $GOROOT/VERSION
// to enable Go to determine the version to use in the version string.
func writeVersion(version string) *doubleError {
	var ver string
	if strings.HasPrefix(version, "go") {
		ver = version
	} else {
		cmd := exec.Command("git", "log", "-n", "1", "--format=format: +%h %cd", version)
		cmd.Dir = filepath.Join(GvmnDir, "repo")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		tag, err := cmd.Output()
		if err != nil {
			return &doubleError{err, fmt.Errorf(stderr.String())}
		}
		ver = "devel" + strings.TrimSpace(string(tag))
	}
	if err := ioutil.WriteFile(filepath.Join(GvmnDir, "versions", version, "VERSION"), []byte(ver), 0666); err != nil {
		return &doubleError{errors.Wrap(err, "failed to write the version to VERSION"), nil}
	}
	return nil
}

// checkout checkouts specified version of the Go repository.
func checkout(version string) *doubleError {
	cmd := exec.Command("git", "archive", "--prefix="+version+"/", version)
	cmd.Dir = filepath.Join(GvmnDir, "repo")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return &doubleError{errors.Wrap(err, "git archive "+version+" failed"), fmt.Errorf(stderr.String())}
	}

	versionsDir := filepath.Join(GvmnDir, "versions")
	if !exist(versionsDir) {
		if err := os.MkdirAll(versionsDir, 0777); err != nil {
			return &doubleError{err, nil}
		}
	}
	cmd = exec.Command("tar", "xf", "-")
	cmd.Dir = versionsDir
	cmd.Stdin = bytes.NewReader(out)
	if stdout, err := cmd.Output(); err != nil {
		return &doubleError{errors.Wrap(err, "tar failed"), fmt.Errorf("%s", stdout)}
	}
	return nil
}

func checkout2(version string) *doubleError {
	versionsDir := filepath.Join(GvmnDir, "versions", version)
	cmd := exec.Command("git", "clone", "--branch="+version, filepath.Join(GvmnDir, "repo"), versionsDir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return &doubleError{errors.Wrap(err, "checkout "+version+" failed"), fmt.Errorf("%s", out)}
	}
	return nil
}

// latestTag reports the latest tag of the Go repository.
func latestTag() (string, error) {
	cmd := exec.Command("git", "rev-list", "--tags", "--max-count=1")
	cmd.Dir = filepath.Join(GvmnDir, "repo")
	out, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "git rev-list failed")
	}
	sha := string(bytes.TrimSuffix(out, []byte("\n")))
	cmd = exec.Command("git", "describe", "--tags", sha)
	cmd.Dir = filepath.Join(GvmnDir, "repo")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	tag, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, stderr.String())
	}
	return string(bytes.TrimSuffix(tag, []byte("\n"))), nil
}

// install installs the specified version of Go.
func install(version string) error {
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

	if err := writeVersion(version); err != nil {
		return err
	}

	if err := build(version); err != nil {
		return err
	}
	return nil
}

// download fetches repository from RepoURL.
func download() *doubleError {
	dir := filepath.Join(GvmnDir, "repo")
	if !exist(dir) {
		out, err := exec.Command("git", "clone", "--bare", RepoURL, dir).CombinedOutput()
		if err != nil {
			return &doubleError{errors.Wrap(err, "cloning repository failed"), fmt.Errorf("%s", out)}
		}
	}

	cmd := exec.Command("git", "fetch", "--tags")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return &doubleError{errors.Wrap(err, "failed to fetch"), nil}
	}
	return nil
}

// runInstall executes install command and return exit code.
func runInstall(args []string) int {
	if len(args) == 0 {
		log.Print("gvmn install: no Go version specified")
		return 1
	}

	if err := download(); err != nil {
		log.Print(err)
		return 1
	}

	if err := install(args[0]); err != nil {
		log.Print(err)
		return 1
	}
	return 0
}
