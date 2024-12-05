package analyze

import (
	"regexp"
	"strings"

	"github.com/blwsh/llmt/cmd/llmt/config"
)

type checker struct {
	a               *config.ProjectAnalyzerConfig
	compiledRegexes map[*config.ProjectAnalyzerConfig]*regexp.Regexp
}

func (c checker) check(filePath string) bool {
	if c.a.NotIn != nil {
		for _, notIn := range *c.a.NotIn {
			if strings.Contains(filePath, notIn) {
				return false
			}
		}
	}

	if c.a.In != nil {
		var inMatchFound = false
		for _, in := range *c.a.In {
			if strings.Contains(filePath, in) {
				inMatchFound = true
				break
			}
		}
		if !inMatchFound {
			return false
		}
	}

	if compile, ok := c.compiledRegexes[c.a]; ok {
		return compile.MatchString(filePath)
	}

	return true
}

func compileRegexesFromConfig(c config.AnalyzeConfig) map[*config.ProjectAnalyzerConfig]*regexp.Regexp {
	var compiledRegexes = make(map[*config.ProjectAnalyzerConfig]*regexp.Regexp)

	for _, analyzer := range c.Analyzers {
		if analyzer.Regex != nil {
			compiledRegex, err := regexp.Compile(*analyzer.Regex)
			if err != nil {
				log.Warnf("failed to compile regex %s: %v", analyzer.Regex, err)
			}

			compiledRegexes[analyzer] = compiledRegex
		}
	}

	return compiledRegexes
}
