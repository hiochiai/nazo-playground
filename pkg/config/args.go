package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func (c *Config) ConfigureWithArgs() error {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `usage: %s [options]
options:
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	port := flag.Uint("p", c.Server.Port, "Port number")
	confDirPath := flag.String("c", c.ConfDirPath, "Configuration directory path")
	flag.Parse()

	if len(*confDirPath) > 0 {
		confFilePath := MakeConfYamlPath(*confDirPath)
		pages, err := parsePagesWithFile(confFilePath)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		c.ConfDirPath = *confDirPath
		c.Pages = pages
	}

	c.Server.Port = *port

	return nil
}
