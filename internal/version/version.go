package version

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

var (
	// Version is the current version of the API
	Version = "0.1.0"
	// BuildTime is the time the binary was built
	BuildTime = "unknown"
	// GitCommit is the git commit hash
	GitCommit = "unknown"
)

// Info contains version information
type Info struct {
	Version       string `json:"version"`
	BuildTime     string `json:"build_time"`
	GitCommit     string `json:"git_commit"`
	GoVersion     string `json:"go_version"`
	OS            string `json:"os"`
	Arch          string `json:"arch"`
	NucleiVersion string `json:"nuclei_version"`
}

// Get returns the version information
func Get() Info {
	nucleiVersion := "unknown"
	if cmd, err := exec.Command("nuclei", "-version").Output(); err == nil {
		nucleiVersion = strings.TrimSpace(string(cmd))
	}

	return Info{
		Version:       Version,
		BuildTime:     BuildTime,
		GitCommit:     GitCommit,
		GoVersion:     runtime.Version(),
		OS:            runtime.GOOS,
		Arch:          runtime.GOARCH,
		NucleiVersion: nucleiVersion,
	}
}

// String returns a string representation of the version information
func String() string {
	info := Get()
	return fmt.Sprintf(`Nuclei API Server
----------------
Version: %s
Build Time: %s
Git Commit: %s
Go Version: %s
OS: %s
Arch: %s
Nuclei Version: %s`,
		info.Version,
		info.BuildTime,
		info.GitCommit,
		info.GoVersion,
		info.OS,
		info.Arch,
		info.NucleiVersion)
}
