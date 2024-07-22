package analyze

import (
	"regexp"
	"strings"

	"github.com/blwsh/llmt/config"
)

func compileRegexesFromConfig(c config.Config) map[*config.ProjectAnalyzerConfig]*regexp.Regexp {
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

func condition(a *config.ProjectAnalyzerConfig, compiledRegexes map[*config.ProjectAnalyzerConfig]*regexp.Regexp) func(filePath string) bool {
	return func(filePath string) bool {
		if a.NotIn != nil {
			for _, notIn := range *a.NotIn {
				if strings.Contains(filePath, notIn) {
					return false
				}
			}
		}

		if a.In != nil {
			var inMatchFound = false
			for _, in := range *a.In {
				if strings.Contains(filePath, in) {
					inMatchFound = true
					break
				}
			}
			if !inMatchFound {
				return false
			}
		}

		if compile, ok := compiledRegexes[a]; ok {
			return compile.MatchString(filePath)
		}

		return true
	}
}
