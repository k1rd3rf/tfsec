package scanner

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/k1rd3rf/tfsec/internal/app/tfsec/parser"
)

// Scanner scans HCL blocks by running all registered checks against them
type Scanner struct {
}

// New creates a new Scanner
func New() *Scanner {
	return &Scanner{}
}

// Scan takes all available hcl blocks and an optional context, and returns a slice of results. Each result indicates a potential security problem.
func (scanner *Scanner) Scan(blocks []*parser.Block) []Result {
	var results []Result
	context := &Context{blocks: blocks}
	for _, block := range blocks {
		for _, check := range GetRegisteredChecks() {
			if check.IsRequiredForBlock(block) {
				for _, result := range check.Run(block, context) {
					if !scanner.checkRangeIgnored(result.RuleID, result.Range) {
						result.Link = fmt.Sprintf("https://github.com/liamg/tfsec/wiki/%s", result.RuleID)
						results = append(results, result)
					}
				}
			}
		}
	}
	return results
}

func (scanner *Scanner) checkRangeIgnored(code RuleID, r parser.Range) bool {
	raw, err := ioutil.ReadFile(r.Filename)
	if err != nil {
		return false
	}
	ignoreAll := "tfsec:ignore:*"
	ignoreCode := fmt.Sprintf("tfsec:ignore:%s", code)
	lines := append([]string{""}, strings.Split(string(raw), "\n")...)
	for number := r.StartLine; number <= r.EndLine; number++ {
		if number <= 0 || number >= len(lines) {
			continue
		}
		if strings.Contains(lines[number], ignoreAll) || strings.Contains(lines[number], ignoreCode) {
			return true
		}
	}

	if r.StartLine-1 > 0 {
		line := lines[r.StartLine-1]
		line = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(line, "//", ""), "#", ""))
		segments := strings.Split(line, " ")
		for _, segment := range segments {
			if segment == ignoreAll || segment == ignoreCode {
				return true
			}
		}

	}

	return false
}
