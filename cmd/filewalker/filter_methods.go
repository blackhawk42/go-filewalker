package main

import (
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/blackhak42/go-filewalker/filter"
)

type FilterMethods map[string]func(pattern string) (filter.Filter, error)

func (fm FilterMethods) GetMethods() []string {
	methods := make([]string, 0, len(fm))
	for m := range fm {
		methods = append(methods, m)
	}

	sort.Strings(methods)

	return methods
}

var AvaiableFilterMethods = FilterMethods(map[string]func(pattern string) (filter.Filter, error){
	// Usual globbing pattern matching the basename
	"glob": func(pattern string) (filter.Filter, error) {
		// Detect errors early
		_, err := filepath.Match(pattern, "")
		if err != nil {
			return nil, err
		}

		f := filter.NewFunctionFilter(func(str string) bool {
			b, _ := filepath.Match(pattern, filepath.Base(str))
			return b
		})

		return f, nil
	},

	// Regex-based matching.
	//
	// Supposedly, as of Go 1.12 the same regex (and, therefore, the same filter)
	// should be able to be reused without lock contention.
	"regex": func(pattern string) (filter.Filter, error) {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}

		f := filter.NewFunctionFilter(func(str string) bool {
			return re.MatchString(str)
		})

		return f, nil
	},

	// Check if the whole path contains a substring
	"path-contains": func(pattern string) (filter.Filter, error) {
		f := filter.NewFunctionFilter(func(str string) bool {
			return strings.Contains(str, pattern)
		})

		return f, nil
	},

	// Check if only the basename contains a substring
	"contains": func(pattern string) (filter.Filter, error) {
		f := filter.NewFunctionFilter(func(str string) bool {
			return strings.Contains(filepath.Base(str), pattern)
		})

		return f, nil
	},

	// Check if the path conatins a suffix
	"suffix": func(pattern string) (filter.Filter, error) {
		f := filter.NewFunctionFilter(func(str string) bool {
			return strings.HasSuffix(str, pattern)
		})

		return f, nil
	},

	// Check if only the basename contains a prefix
	"prefix": func(pattern string) (filter.Filter, error) {
		f := filter.NewFunctionFilter(func(str string) bool {
			return strings.HasPrefix(filepath.Base(str), pattern)
		})

		return f, nil
	},
})
