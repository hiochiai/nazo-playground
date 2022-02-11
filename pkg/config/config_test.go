package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestEvalPagesContents(t *testing.T) {
	noErr := fmt.Sprintf("%v", nil)

	tests := []struct {
		contents  string
		expect    string
		expectErr string
	}{
		{
			contents:  "(あれは$(Y=`date +%Y`; expr $Y - 2016)年前)",
			expect:    fmt.Sprintf(`(あれは%d年前)`, time.Now().Year()-2016),
			expectErr: noErr,
		},
		{
			contents:  "$(echo 1) + $(echo 1) = $(expr 1 + 1) です",
			expect:    `1 + 1 = 2 です`,
			expectErr: noErr,
		},
		{
			contents:  "$(invalid command)",
			expect:    "",
			expectErr: "command not found",
		},
	}

	for i, test := range tests {
		actual, err := EvalPagesContents(test.contents)

		if !strings.Contains(fmt.Sprintf(`%v`, err), test.expectErr) {
			t.Fatalf("test[%d] expects error %v, but %v", i, test.expectErr, err)
		}

		if !reflect.DeepEqual(test.expect, actual) {
			t.Fatalf("test[%d] expects %v, but %v", i, test.expect, actual)
		}
	}
}

func TestParsePagesWithFile(t *testing.T) {
	noErr := fmt.Sprintf("%v", nil)

	tests := []struct {
		yaml        string
		expectPages []PageConfig
		expectErr   string
	}{
		{
			yaml:        ` `,
			expectPages: nil,
			expectErr:   noErr,
		},
		{
			yaml: `
pages:
  - id: 0
    answer: a
    contents: con ten ts
  - id: 00000000-0000-0000-0000-000000000000
    answer: test
    contents: |
      {"json":"contents"}`,
			expectPages: []PageConfig{
				{Id: "0", Answer: "a", Contents: "con ten ts"},
				{
					Id:       `00000000-0000-0000-0000-000000000000`,
					Answer:   `test`,
					Contents: `{"json":"contents"}`,
				},
			},
			expectErr: noErr,
		},
		{
			yaml: `
pages:
  - id: 0
    answer: a
    contents: con ten ts
  - id: 0
    answer: b
    contents: duplicated id`,
			expectPages: nil,
			expectErr:   `invalid configuration`,
		},
		{
			yaml: `
pages:
  - id: 0
    answer: a
    contents: $(invalid)`,
			expectPages: nil,
			expectErr:   `command not found`,
		},
		{
			yaml:        ``,
			expectPages: nil,
			expectErr:   `no such file`,
		},
	}

	for i, test := range tests {
		tempDir, _ := ioutil.TempDir("", "")
		defer os.RemoveAll(tempDir)
		yamlPath := filepath.Join(tempDir, "conf.yaml")
		if len(test.yaml) > 0 {
			_ = ioutil.WriteFile(yamlPath, []byte(test.yaml), 0644)
		}

		pages, err := parsePagesWithFile(yamlPath)

		if !strings.Contains(fmt.Sprintf(`%v`, err), test.expectErr) {
			t.Fatalf("test[%d] expects error %v, but %v", i, test.expectErr, err)
		}

		if !reflect.DeepEqual(test.expectPages, pages) {
			t.Fatalf("test[%d] expects %v, but %v", i, test.expectPages, pages)
		}
	}
}
