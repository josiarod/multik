package util

import "path"

// MathcAny returns tru if name matches any glob pattern in patterns
// If patterns is empty, it returns true.
func MatchAny(name string, patterns []string) bool {
	if len(patterns) == 0 {
		return true
	}
	for _, p := range patterns {
		if ok, _ := path.Match(p, name); ok {
			return true
		}
	}
	return false
}
