package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
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

// runInstall executes install command and return exit code.
func runInstall(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "gvmn install: no Go version specified")
		return 1
	}
	dir := filepath.Join(GvmnDir, "repo")
	if !exist(dir) {
		_, err := exec.Command("git", "clone", "--bare", RepoURL, dir).CombinedOutput()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}

	cmd := exec.Command("git", "fetch", "--tags")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed fetch"))
		return 1
	}

	return install(args[0])
}

func install(version string) int {
	if version == "latest" {
		var err error
		version, err = latestTag()
		if err != nil {
			fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed get the latest version"))
			return 1
		}
	}
	cmd := exec.Command("git", "archive", "--prefix="+version+"/", version)
	cmd.Dir = filepath.Join(GvmnDir, "repo")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed git archive"))
		fmt.Fprintln(os.Stderr, stderr.String())
		return 1
	}

	versionsDir := filepath.Join(GvmnDir, "versions")
	if !exist(versionsDir) {
		if err := os.MkdirAll(versionsDir, 0777); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}
	cmd = exec.Command("tar", "xf", "-")
	cmd.Dir = versionsDir
	cmd.Stdin = bytes.NewReader(out)
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed tar"))
		return 1
	}

	var env []string
	if goroot, err := exec.Command("go", "env", "GOROOT").Output(); err == nil {
		env = append(os.Environ(), "GOROOT_BOOTSTRAP="+string(bytes.TrimSuffix(goroot, []byte("\n"))))
	}

	var ver string
	if strings.HasPrefix(version, "go") {
		ver = version
	} else {
		cmd = exec.Command("git", "log", "-n", "1", "--format=format: +%h %cd", version)
		cmd.Dir = filepath.Join(GvmnDir, "repo")
		tag, err := cmd.Output()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		ver = "devel" + strings.TrimSpace(string(tag))
	}
	if err := ioutil.WriteFile(filepath.Join(GvmnDir, "versions", version, "VERSION"), []byte(ver), 0666); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed write a version to VERSION"))
		return 1
	}

	cmd = exec.Command("./make.bash")
	cmd.Dir = filepath.Join(versionsDir, version, "src")
	cmd.Env = env
	var buf bytes.Buffer
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed install"))
		fmt.Fprintln(os.Stderr, buf.String())
		return 1
	}
	return 0
}

func latestTag() (string, error) {
	cmd := exec.Command("git", "rev-list", "--tags", "--max-count=1")
	cmd.Dir = filepath.Join(GvmnDir, "repo")
	out, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "failed git rev-list")
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
