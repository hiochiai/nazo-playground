package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	LogLevel    string
	Server      ServerConfig `yaml:"server"`
	ConfDirPath string
	pagesConfigFile
}

type ServerConfig struct {
	Port uint `yaml:"port"`
}

type PageConfig struct {
	Id       string `yaml:"id"`
	Answer   string `yaml:"answer"`
	Contents string `yaml:"contents"`
}

type pagesConfigFile struct {
	Pages []PageConfig `yaml:"pages"`
}

func NewDefaultConfig() *Config {
	return &Config{
		LogLevel: "info",
		Server: ServerConfig{
			Port: 8080,
		},
	}
}

func MakeConfYamlPath(confDirPath string) string {
	return filepath.Join(confDirPath, `conf.yaml`)
}

func MakeConfStaticDirPath(confDirPath string) string {
	return filepath.Join(confDirPath, `static`)
}

func EvalPagesContents(contents string) (string, error) {

	for {
		// Replace Special string /$(.*)/ to shell script.
		re := regexp.MustCompile(`\$\(([^)]*)\)`)
		matchs := re.FindStringSubmatch(contents)
		if len(matchs) != 2 {
			return contents, nil
		}

		var stderr bytes.Buffer
		cmd := exec.Command(`bash`, "-c", matchs[1])
		cmd.Stderr = &stderr
		stdout, err := cmd.Output()
		msg := strings.TrimSpace(string(stdout))
		if err != nil {
			errMsg := strings.TrimSpace(stderr.String())
			if len(errMsg) == 0 {
				errMsg = msg
			}
			return "", fmt.Errorf(`failed to exec "%s": %s, (%v)`, matchs[1], errMsg, err)
		}

		contents = strings.Replace(contents, matchs[0], msg, 1)
	}
}

func parsePagesWithFile(path string) ([]PageConfig, error) {

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	c := pagesConfigFile{}
	if err := yaml.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}

	// ID Duplication check
	ids := make(map[string]struct{})
	for i := range c.Pages {
		if _, exists := ids[c.Pages[i].Id]; exists {
			return nil, fmt.Errorf(`invalid configuration: "id: %v" is duplicated`, c.Pages[i].Id)
		}
		ids[c.Pages[i].Id] = struct{}{}
	}

	// Eval string format check
	for i := range c.Pages {
		if _, err := EvalPagesContents(c.Pages[i].Contents); err != nil {
			return nil, fmt.Errorf(`invalid configuration: %v`, err)
		}
	}

	return c.Pages, nil
}
