package build //nolint: testpackage // Need to test private parseBuildInfo

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"runtime/debug"
	"testing"
	"time"

	"github.com/FollowTheProcess/test"
)

var update = flag.Bool("update", false, "Update golden files")

func TestString(t *testing.T) {
	when := time.Date(2024, time.October, 6, 15, 39, 24, 100, time.UTC)

	tests := []struct {
		name string    // The name of the test case
		file string    // Name of the file containing the expected output (relative to testdata)
		info BuildInfo // The build info under test
	}{
		{
			name: "simple",
			info: BuildInfo{
				Main: Module{
					Path:    "github.com/SomeGuy/project",
					Version: "v1.2.3",
					Sum:     "WwdigHlEGoXEzt8n/VGpqrNkD3j5gHsqBjYduqTqRE0=",
				},
				Time: when,
				Settings: map[string]string{
					"-buildmode":  "exe",
					"-compiler":   "gc",
					"-ldflags":    "-X main.version=dev",
					"CGO_ENABLED": "0",
					"GOAMD64":     "v4",
				},
				Go:      "go1.23.2",
				Path:    "github.com/SomeGuy/project",
				OS:      "darwin",
				Arch:    "amd64",
				VCS:     "git",
				Commit:  "5e8b8a68867eff5f754bfecdbc8baeb2c14c711c",
				Version: "v1.2.3",
				Dirty:   true,
			},
			file: "info.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.info.String()
			file := filepath.Join(test.Data(t), tt.file)

			if *update {
				t.Logf("Updating %s\n", file)
				err := os.WriteFile(file, []byte(got), os.ModePerm)
				test.Ok(t, err)
			}

			test.File(t, got, file)
		})
	}
}

func TestJSON(t *testing.T) {
	when := time.Date(2024, time.October, 6, 15, 39, 24, 100, time.UTC)

	tests := []struct {
		name string    // The name of the test case
		file string    // Name of the file containing the expected output (relative to testdata)
		info BuildInfo // The build info under test
	}{
		{
			name: "simple",
			info: BuildInfo{
				Main: Module{
					Path:    "github.com/SomeGuy/project",
					Version: "v1.2.3",
					Sum:     "WwdigHlEGoXEzt8n/VGpqrNkD3j5gHsqBjYduqTqRE0=",
				},
				Time: when,
				Settings: map[string]string{
					"-buildmode":  "exe",
					"-compiler":   "gc",
					"-ldflags":    "-X main.version=dev",
					"CGO_ENABLED": "0",
					"GOAMD64":     "v4",
				},
				Go:      "go1.23.2",
				Path:    "github.com/SomeGuy/project",
				OS:      "darwin",
				Arch:    "amd64",
				VCS:     "git",
				Commit:  "5e8b8a68867eff5f754bfecdbc8baeb2c14c711c",
				Version: "v1.2.3",
				Dirty:   true,
			},
			file: "info.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := json.NewEncoder(buf).Encode(tt.info)
			test.Ok(t, err)

			got := buf.String()
			file := filepath.Join(test.Data(t), tt.file)

			if *update {
				t.Logf("Updating %s\n", file)
				err := os.WriteFile(file, []byte(got), os.ModePerm)
				test.Ok(t, err)
			}

			test.File(t, got, file)
		})
	}
}

func TestParse(t *testing.T) {
	when := time.Date(2024, time.October, 6, 10, 39, 12, 0, time.UTC)

	tests := []struct {
		name string           // Name of the test case
		dbg  *debug.BuildInfo // The debug info to parse
		want BuildInfo        // The expected result
	}{
		{
			name: "full",
			dbg: &debug.BuildInfo{
				GoVersion: "go1.23.2",
				Path:      "github.com/SomeGuy/project",
				Main: debug.Module{
					Path:    "github.com/SomeGuy/project",
					Version: "v1.2.3",
					Sum:     "WwdigHlEGoXEzt8n/VGpqrNkD3j5gHsqBjYduqTqRE0=",
					Replace: nil,
				},
				Settings: []debug.BuildSetting{
					{Key: "GOOS", Value: "darwin"},
					{Key: "GOARCH", Value: "amd64"},
					{Key: "vcs", Value: "git"},
					{Key: "vcs.revision", Value: "5e8b8a68867eff5f754bfecdbc8baeb2c14c711c"},
					{Key: "vcs.time", Value: "2024-10-06T10:39:12Z"},
					{Key: "vcs.modified", Value: "true"},
					{Key: "-buildmode", Value: "exe"},
					{Key: "-compiler", Value: "gc"},
					{Key: "CGO_ENABLED", Value: "0"},
					{Key: "GOAMD64", Value: "v4"},
					{Key: "-ldflags", Value: "-X main.version=dev"},
				},
			},
			want: BuildInfo{
				Main: Module{
					Path:    "github.com/SomeGuy/project",
					Version: "v1.2.3",
					Sum:     "WwdigHlEGoXEzt8n/VGpqrNkD3j5gHsqBjYduqTqRE0=",
				},
				Time: when,
				Settings: map[string]string{
					"-buildmode":  "exe",
					"-compiler":   "gc",
					"-ldflags":    "-X main.version=dev",
					"CGO_ENABLED": "0",
					"GOAMD64":     "v4",
				},
				Go:      "go1.23.2",
				Path:    "github.com/SomeGuy/project",
				OS:      "darwin",
				Arch:    "amd64",
				VCS:     "git",
				Commit:  "5e8b8a68867eff5f754bfecdbc8baeb2c14c711c",
				Version: "v1.2.3",
				Dirty:   true,
			},
		},
		{
			name: "bad vcs.time",
			dbg: &debug.BuildInfo{
				GoVersion: "go1.23.2",
				Path:      "github.com/SomeGuy/project",
				Main: debug.Module{
					Path:    "github.com/SomeGuy/project",
					Version: "v1.2.3",
					Sum:     "WwdigHlEGoXEzt8n/VGpqrNkD3j5gHsqBjYduqTqRE0=",
					Replace: nil,
				},
				Settings: []debug.BuildSetting{
					{Key: "GOOS", Value: "darwin"},
					{Key: "GOARCH", Value: "amd64"},
					{Key: "vcs", Value: "git"},
					{Key: "vcs.revision", Value: "5e8b8a68867eff5f754bfecdbc8baeb2c14c711c"},
					{Key: "vcs.time", Value: "not a time"}, // <- This bit here is wrong
					{Key: "vcs.modified", Value: "true"},
					{Key: "-buildmode", Value: "exe"},
					{Key: "-compiler", Value: "gc"},
					{Key: "CGO_ENABLED", Value: "0"},
					{Key: "GOAMD64", Value: "v4"},
					{Key: "-ldflags", Value: "-X main.version=dev"},
				},
			},
			want: BuildInfo{
				Main: Module{
					Path:    "github.com/SomeGuy/project",
					Version: "v1.2.3",
					Sum:     "WwdigHlEGoXEzt8n/VGpqrNkD3j5gHsqBjYduqTqRE0=",
				},
				Settings: map[string]string{
					"-buildmode":  "exe",
					"-compiler":   "gc",
					"-ldflags":    "-X main.version=dev",
					"CGO_ENABLED": "0",
					"GOAMD64":     "v4",
				},
				Go:      "go1.23.2",
				Path:    "github.com/SomeGuy/project",
				OS:      "darwin",
				Arch:    "amd64",
				VCS:     "git",
				Commit:  "5e8b8a68867eff5f754bfecdbc8baeb2c14c711c",
				Version: "v1.2.3",
				Dirty:   true,
			},
		},
		{
			name: "bad vcs.modified",
			dbg: &debug.BuildInfo{
				GoVersion: "go1.23.2",
				Path:      "github.com/SomeGuy/project",
				Main: debug.Module{
					Path:    "github.com/SomeGuy/project",
					Version: "v1.2.3",
					Sum:     "WwdigHlEGoXEzt8n/VGpqrNkD3j5gHsqBjYduqTqRE0=",
					Replace: nil,
				},
				Settings: []debug.BuildSetting{
					{Key: "GOOS", Value: "darwin"},
					{Key: "GOARCH", Value: "amd64"},
					{Key: "vcs", Value: "git"},
					{Key: "vcs.revision", Value: "5e8b8a68867eff5f754bfecdbc8baeb2c14c711c"},
					{Key: "vcs.time", Value: "2024-10-06T10:39:12Z"},
					{Key: "vcs.modified", Value: "notabool"}, // <- This bit here is wrong
					{Key: "-buildmode", Value: "exe"},
					{Key: "-compiler", Value: "gc"},
					{Key: "CGO_ENABLED", Value: "0"},
					{Key: "GOAMD64", Value: "v4"},
					{Key: "-ldflags", Value: "-X main.version=dev"},
				},
			},
			want: BuildInfo{
				Main: Module{
					Path:    "github.com/SomeGuy/project",
					Version: "v1.2.3",
					Sum:     "WwdigHlEGoXEzt8n/VGpqrNkD3j5gHsqBjYduqTqRE0=",
				},
				Settings: map[string]string{
					"-buildmode":  "exe",
					"-compiler":   "gc",
					"-ldflags":    "-X main.version=dev",
					"CGO_ENABLED": "0",
					"GOAMD64":     "v4",
				},
				Time:    when,
				Go:      "go1.23.2",
				Path:    "github.com/SomeGuy/project",
				OS:      "darwin",
				Arch:    "amd64",
				VCS:     "git",
				Commit:  "5e8b8a68867eff5f754bfecdbc8baeb2c14c711c",
				Version: "v1.2.3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseBuildInfo(tt.dbg)
			test.Diff(t, got, tt.want)
		})
	}
}
