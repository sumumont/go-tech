package common

import "regexp"

var (
	Regexps = map[string]*regexp.Regexp{
		"image.name":       regexp.MustCompile("emailRegexString"),
		"image.tag":        regexp.MustCompile("^[a-zA-Z0-9]+([._-][a-zA-Z0-9]+)*([._-][a-zA-Z0-9]+)*$"),
		"image.tag.failed": regexp.MustCompile("([.-][.-])+"),
	}
)
