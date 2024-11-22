# Build

[![License](https://img.shields.io/github/license/FollowTheProcess/build)](https://github.com/FollowTheProcess/build)
[![Go Reference](https://pkg.go.dev/badge/github.com/FollowTheProcess/build.svg)](https://pkg.go.dev/github.com/FollowTheProcess/build)
[![Go Report Card](https://goreportcard.com/badge/github.com/FollowTheProcess/build)](https://goreportcard.com/report/github.com/FollowTheProcess/build)
[![GitHub](https://img.shields.io/github/v/release/FollowTheProcess/build?logo=github&sort=semver)](https://github.com/FollowTheProcess/build)
[![CI](https://github.com/FollowTheProcess/build/workflows/CI/badge.svg)](https://github.com/FollowTheProcess/build/actions?query=workflow%3ACI)
[![codecov](https://codecov.io/gh/FollowTheProcess/build/branch/main/graph/badge.svg)](https://codecov.io/gh/FollowTheProcess/build)

## Project Description

Ridiculously simple Go build info 🛠️

## Installation

```shell
go get github.com/FollowTheProcess/build@latest
```

## Quickstart

```go
package main

import (
	"fmt"
	"os"

	"github.com/FollowTheProcess/build"
)

func main() {
	info, ok := build.Info()
	if !ok {
		fmt.Fprintf(os.Stderr, "could not get build info")
		os.Exit(1)
	}

	fmt.Printf("%s\n", info)
}
```

Gets you...

```shell
go:           go1.23.2
path:         github.com/FollowTheProcess/build/cmd/build
os:           darwin
arch:         amd64
vcs:          git
version:      (devel)
commit:       5e8b8a68867eff5f754bfecdbc8baeb2c14c711c
dirty:        true
time:         2024-10-06T10:39:12Z
main:         mod  github.com/FollowTheProcess/build  (devel)  
-buildmode:   exe
-compiler:    gc
-ldflags:     -X main.version=dev
CGO_ENABLED:  0
GOAMD64:      v1
```

> [!TIP]
> It's also JSON serialisable!

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/FollowTheProcess/build"
)

func main() {
	info, _ := build.Info()
	if err := json.NewEncoder(os.Stdout).Encode(info); err != nil {
		fmt.Fprintf(os.Stderr, "could not write JSON: %v\n", err)
		os.Exit(1)
	}
}
```

Gets you...

```json
{
  "main": {
    "path": "github.com/FollowTheProcess/build",
    "version": "(devel)"
  },
  "time": "2024-10-06T10:39:12Z",
  "settings": {
    "-buildmode": "exe",
    "-compiler": "gc",
    "-ldflags": "-X main.version=dev",
    "CGO_ENABLED": "0",
    "GOAMD64": "v1"
  },
  "go": "go1.23.2",
  "path": "github.com/FollowTheProcess/build/cmd/build",
  "os": "darwin",
  "arch": "amd64",
  "vcs": "git",
  "version": "(devel)",
  "dirty": true
}
```

`build.Info` returns a `BuildInfo` struct from which you can take any component of the build info:

```go
package main

import (
	"fmt"
	"os"

	"github.com/FollowTheProcess/build"
)

func main() {
	info, ok := build.Info()
	if !ok {
		fmt.Fprintf(os.Stderr, "could not get build info")
		os.Exit(1)
	}

	fmt.Printf("Version: %s\n", info.Version)
  fmt.Printf("Commit: %s\n", info.Commit)
}
```

### Credits

This package is wholly based on the Go internal implementation of `runtime/debug.BuildInfo`, this is just a slightly nicer wrapper that makes it easier to access common settings
