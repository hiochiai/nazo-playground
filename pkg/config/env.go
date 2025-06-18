package config

import (
	"fmt"
	"os"
	"strconv"
)

func (c *Config) ConfigureWithEnv() error {
	if v := os.Getenv("NAZO_LOG"); len(v) > 0 {
		c.LogLevel = v
	}

	if v := os.Getenv("NAZO_PORT"); len(v) > 0 {
		port, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return fmt.Errorf("invalid value \"%v\" for environment value NAZO_PORT", v)
		}
		c.Server.Port = uint(port)
	}

	if v := os.Getenv("NAZO_CONFDIR"); len(v) > 0 {
		confFilePath := MakeConfYamlPath(v)
		pages, err := parsePagesWithFile(confFilePath)
		if err != nil {
			return fmt.Errorf("invalid value \"%v\" for environment value NAZO_CONF: %w", v, err)
		}
		c.ConfDirPath = v
		c.Pages = pages
	}

	return nil
}
