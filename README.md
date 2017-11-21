# Gvmn

Go Version Manager Next

## Install

To install, use `go get`:

```bash
$ go get github.com/elpinal/gvmn/cmd/gvmn
```

## Command gvmn

Gvmn is a tool for managing Go versions.

Usage:

```bash
gvmn command [arguments]
```

The commands are:

```bash
get         download and install Go
list        list installed Go versions
use         select a Go version to use
remove      remove Go versions
run         execute the specified Go version
version     print gvmn version
```

Use "gvmn help [command]" for more information about a command.

### Download and install Go

Get specific Go versions:

```bash
$ gvmn get go1.7
```

Get the latest tagged Go, such as go1.7 and go1.7rc6:

```bash
$ gvmn get stable
```

If you want to get Go without building (which takes a few minutes), get binaries:

```bash
$ gvmn get -b go1.7
```

#### Install the latest developing version

```bash
$ gvmn get tip
```

### List installed Go versions

Know what Go versions is installed or downloaded:

```bash
$ gvmn list
  go1.6
* go1.7
  go1.7beta2
  go1.7rc6
```

### Select a Go version to use

Select a Go version to use:

```bash
$ gvmn use go1.6
```

### Remove Go versions

Say goodbye to particular Go versions:

```bash
$ gvmn remove go1.7
```

### Execute the specified Go version

Execute another Go version:

```bash
$ gvmn run go1.5 get github.com/elpinal/gvmn/cmd/gvmn
```

### Print gvmn version

```bash
$ gvmn version
```

## Contribution

1. Fork ([https://github.com/elpinal/gvmn/fork](https://github.com/elpinal/gvmn/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[elpinal](https://github.com/elpinal)
