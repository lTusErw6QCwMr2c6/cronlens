package schedule

import (
	"sort"
	"strings"
)

// SuggestTemplates returns templates whose name or description contains the query (case-insensitive).
func SuggestTemplates(query string) []Template {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return CommonTemplates
	}
	var results []Template
	for _, t := range CommonTemplates {
		if strings.Contains(strings.ToLower(t.Name), q) ||
			strings.Contains(strings.ToLower(t.Description), q) {
			results = append(results, t)
		}
	}
	return results
}

// MatchTemplate returns the closest template name match for a given partial name.
// Returns empty string if no reasonable match is found.
func MatchTemplate(partial string) string {
	partial = strings.ToLower(strings.TrimSpace(partial))
	if partial == "" {
		return ""
	}
	var candidates []string
	for _, t := range CommonTemplates {
		if strings.HasPrefix(strings.ToLower(t.Name), partial) {
			candidates = append(candidates, t.Name)
		}
	}
	if len(candidates) == 0 {
		return ""
	}
	sort.Strings(candidates)
	return candidates[0]
}
