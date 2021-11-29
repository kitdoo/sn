package path

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/kitdoo/sn/internal/version"
)

func prefix() string {
	if version.EnvPrefix != "" {
		return version.EnvPrefix + "_"
	}
	return ""
}

func getEnvVar(v string) string {
	p := os.Getenv(prefix() + v)
	if len(p) > 0 && p[0] == '~' {
		p = filepath.Join(homedir(), p[1:])
	}
	return p
}

func homedir() string {
	usr, _ := user.Current()
	return usr.HomeDir
}

func BinDir() string {
	if dir := getEnvVar("BIN_DIR"); dir != "" {
		return strings.TrimRight(dir, "/")
	}
	return "/opt/bin"
}

func VarDir() string {
	if dir := getEnvVar("VAR_DIR"); dir != "" {
		return strings.TrimRight(dir, "/")
	}
	return "/var"
}

func EtcDir() string {
	if dir := getEnvVar("ETC_DIR"); dir != "" {
		return strings.TrimRight(dir, "/")
	}
	return "/etc/" + strings.ToLower(version.AppName)
}

func LibDir() string {
	return VarDir() + "/lib/" + strings.ToLower(version.AppName)
}

func CertsCacheDir() string {
	return LibDir() + "/certs"
}
