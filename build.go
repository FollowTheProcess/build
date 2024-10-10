// Package build provides a simple mechanism for getting at debug.BuildInfo in a more ergonomic way.
package build

import (
	"fmt"
	"maps"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

// TableWriter config, used for showing build info.
const (
	minWidth = 1   // Min cell width
	tabWidth = 8   // Tab width in spaces
	padding  = 2   // Padding
	padChar  = ' ' // Char to pad with
	flags    = 0   // Config flags
)

// BuildInfo contains the build information of a Go binary.
type BuildInfo struct { //nolint:revive // Yes it stutters but having the function be build.Info is worth it
	Main     Module            `json:"main,omitempty"`     // The main module
	Time     time.Time         `json:"time,omitempty"`     // The modification time associated with Commit
	Settings map[string]string `json:"settings,omitempty"` // The remaining settings
	Go       string            `json:"go,omitempty"`       // The Go toolchain version used to build the binary
	Path     string            `json:"path,omitempty"`     // The package path of the main package
	OS       string            `json:"os,omitempty"`       // The value of $GOOS
	Arch     string            `json:"arch,omitempty"`     // The value of $GOARCH
	VCS      string            `json:"vcs,omitempty"`      // The version control system for the source tree where the build ran
	Commit   string            `json:"commit,omitempty"`   // The SHA1 of the current commit when the build ran
	Version  string            `json:"version,omitempty"`  // The module version
	Dirty    bool              `json:"dirty,omitempty"`    // Whether the source tree had local modifications at the time of the build
}

// String implements [fmt.Stringer] for [BuildInfo].
func (b BuildInfo) String() string {
	s := &strings.Builder{}
	tab := tabwriter.NewWriter(s, minWidth, tabWidth, padding, padChar, flags)
	fmt.Fprintf(tab, "go:\t%s\n", b.Go)
	fmt.Fprintf(tab, "path:\t%s\n", b.Path)
	fmt.Fprintf(tab, "os:\t%s\n", b.OS)
	fmt.Fprintf(tab, "arch:\t%s\n", b.Arch)
	fmt.Fprintf(tab, "vcs:\t%s\n", b.VCS)
	fmt.Fprintf(tab, "version:\t%s\n", b.Version)
	fmt.Fprintf(tab, "commit:\t%s\n", b.Commit)
	fmt.Fprintf(tab, "dirty:\t%v\n", b.Dirty)
	fmt.Fprintf(tab, "time:\t%s\n", b.Time.Format(time.RFC3339))
	fmt.Fprintf(tab, "main:\t%s\n", writeModule("mod", b.Main))

	// Sort the settings map for deterministic printing
	keys := slices.Sorted(maps.Keys(b.Settings))
	for _, key := range keys {
		fmt.Fprintf(tab, "%s:\t%s\n", key, b.Settings[key])
	}

	tab.Flush()
	return s.String()
}

// A Module describes a single module included in a build.
type Module struct {
	Path    string `json:"path,omitempty"`    // module path
	Version string `json:"version,omitempty"` // module version
	Sum     string `json:"sum,omitempty"`     // checksum
}

// Info returns the Go build info for the current binary.
func Info() (info BuildInfo, ok bool) {
	dbg, ok := debug.ReadBuildInfo()
	if !ok {
		return BuildInfo{}, false
	}

	return parseBuildInfo(dbg), true
}

func parseBuildInfo(dbg *debug.BuildInfo) BuildInfo {
	info := BuildInfo{
		Main: Module{
			Path:    dbg.Main.Path,
			Version: dbg.Main.Version,
			Sum:     dbg.Main.Sum,
		},
		Go:       dbg.GoVersion,
		Path:     dbg.Path,
		Version:  dbg.Main.Version,
		Settings: make(map[string]string, len(dbg.Settings)),
	}

	for _, setting := range dbg.Settings {
		switch setting.Key {
		case "GOOS":
			info.OS = setting.Value
		case "GOARCH":
			info.Arch = setting.Value
		case "vcs":
			info.VCS = setting.Value
		case "vcs.revision":
			info.Commit = setting.Value
		case "vcs.time":
			t, err := time.Parse(time.RFC3339, setting.Value)
			if err != nil {
				// Skip this setting
				continue
			}
			info.Time = t
		case "vcs.modified":
			modified, err := strconv.ParseBool(setting.Value)
			if err != nil {
				// Skip this setting
				continue
			}
			info.Dirty = modified
		default:
			// Add any remaining settings into our info
			info.Settings[setting.Key] = setting.Value
		}
	}

	return info
}

func writeModule(word string, m Module) string {
	buf := &strings.Builder{}
	buf.WriteString(word)
	buf.WriteByte('\t')
	buf.WriteString(m.Path)
	buf.WriteByte('\t')
	buf.WriteString(m.Version)
	buf.WriteByte('\t')
	buf.WriteString(m.Sum)

	return buf.String()
}
