package config

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/go-ozzo/ozzo-validation/v4"
)

type GRPC struct {
	LogRequest        bool   `yaml:"logRequests" default:"false"`
	Reflection        bool   `yaml:"reflection" default:"false"`
	RequireClientCert bool   `yaml:"requireClientCert"`
	ListenAddress     string `yaml:"listenAddress" default:"0.0.0.0:7777"`
}

func (g *GRPC) Validate() (err error) {
	return validation.ValidateStruct(g,
		validation.Field(&g.ListenAddress, validation.By(validateListenAddress)),
	)
}

func validateListenAddress(value interface{}) error {
	connectionUri := value.(string)
	re := regexp.MustCompile("^[aA-zZ]+://")
	splits := re.Split(connectionUri, -1)
	if len(splits) > 1 {
		connectionUri = splits[1]
	}

	s := strings.Split(connectionUri, ":")

	//nolint:gomnd
	if len(s) != 2 {
		return fmt.Errorf("invalid listen address '%s'", s)
	}

	p, e := strconv.Atoi(s[1])
	if e != nil {
		return fmt.Errorf("invalid port '%s': %s", s[1], e.Error())
	}

	port := uint16(p)
	if port <= 0 || port >= 65535 {
		return fmt.Errorf("invalid port '%d'", port)
	}

	if !govalidator.IsHost(s[0]) {
		return fmt.Errorf("invalid address '%s'", s[0])
	}

	return nil
}
