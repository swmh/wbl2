package main

import (
	"regexp"
)

type Matcher interface {
	Match([]byte) []Match
}

type RegexpMatcher struct {
	re *regexp.Regexp
}

func NewRegexpMatcher(pattern string, fixed, ignoreCase bool) (*RegexpMatcher, error) {
	if fixed {
		pattern = regexp.QuoteMeta(pattern)
	}

	if ignoreCase {
		pattern = "(?i)" + pattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &RegexpMatcher{re: re}, nil
}

func (r *RegexpMatcher) Match(b []byte) []Match {
	e := r.re.FindAllIndex(b, -1)
	if len(e) == 0 {
		return nil
	}

	result := make([]Match, len(e))

	for i, v := range e {
		result[i] = Match{
			start: v[0],
			len:   v[1] - v[0],
		}
	}

	return result
}
