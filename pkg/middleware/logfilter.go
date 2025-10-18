package middleware

import "strings"

type LogFilter struct {
	SkipPrefixes []string
	SkipExact    []string
	SkipStatus   []int
}

func DefaultLogFilter() *LogFilter {
	return &LogFilter{
		SkipPrefixes: []string{
			"/.well-known",
		},
		SkipExact: []string{
			"/favicon.ico",
			"/robots.txt",
			"/sitemap.xml",
			"/ads.txt",
			"/security.txt",
			"/health",
		},
		SkipStatus: []int{304},
	}
}

func (f *LogFilter) ShouldSkip(path string, status int) bool {
	for _, exact := range f.SkipExact {
		if path == exact {
			return true
		}
	}

	for _, prefix := range f.SkipPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	if (status == 200 || status == 304) && strings.HasPrefix(path, "/static/") {
		return true
	}
	return false
}
