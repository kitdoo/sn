package version

import (
	"strings"
	"sync"
)

const (
	BuildTypeProduction  BuildType = "production"
	BuildTypeDevelopment BuildType = "development"
)

type BuildType = string

var (
	Version   = "0.0.0"
	Commit    = ""
	Build     = BuildTypeDevelopment
	AppName   = "sn"
	EnvPrefix = "SN"
)

func IsDevelopment() bool {
	return strings.ToLower(Build) == BuildTypeDevelopment
}

var (
	once        sync.Once
	fullVersion string
)

func FullVersion() string {
	once.Do(func() {
		fv := &strings.Builder{}
		fv.WriteString(Version)
		if strings.ToLower(Build) == BuildTypeDevelopment {
			fv.WriteString("dev")
		}
		if Commit != "" {
			fv.WriteString("/")
			fv.WriteString(Commit)
		}
		fullVersion = fv.String()
	})
	return fullVersion
}
