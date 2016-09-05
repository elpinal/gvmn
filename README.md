# gvmn

Go Version Manager Next

## Description

gvmn is a tool for managing Go versions.

## Usage

```bash
$ gvmn
gvmn is a tool for managing Go versions.

Usage:

        gvmn command [arguments]

The commands are:

        get         download and install Go
        list        list installed Go versions
        use         select a Go version to use
        remove      remove Go versions
        version     print gvmn version

Use "gvmn help [command]" for more information about a command.
```

### Example

Get specific Go versions:

```bash
$ gvmn get go1.7
```

Get the latest tagged Go, such as go1.7 and go1.7rc6:

```bash
$ gvmn get latest
```

Know what Go versions is installed or downloaded:

```bash
$ # Now, go1.7 is selected to use.
$ # And, go1.4 is just downloaded.
$ gvmn list
  go1.4 (just downloaded; not installed)
* go1.6
  go1.7
  go1.7beta2
  go1.7rc1
  go1.7rc6
```

Say goodbye to particular Go versions:

```bash
$ gvmn remove go1.7
```

## Install

To install, use `go get`:

```bash
$ go get github.com/susp/gvmn/cmd/gvmn
```

## Contribution

1. Fork ([https://github.com/susp/gvmn/fork](https://github.com/susp/gvmn/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[susp](https://github.com/susp)
