package schedule

import "fmt"

// Alias maps a human-friendly shorthand (e.g. "@hourly") to a standard
// 5-field cron expression.
type Alias struct {
	Name        string
	Expression  string
	Description string
}

// predefinedAliases contains the standard cron aliases supported by most
// cron implementations.
var predefinedAliases = []Alias{
	{"@yearly", "0 0 1 1 *", "Run once a year at midnight on January 1st"},
	{"@annually", "0 0 1 1 *", "Same as @yearly"},
	{"@monthly", "0 0 1 * *", "Run once a month at midnight on the first day"},
	{"@weekly", "0 0 * * 0", "Run once a week at midnight on Sunday"},
	{"@daily", "0 0 * * *", "Run once a day at midnight"},
	{"@midnight", "0 0 * * *", "Same as @daily"},
	{"@hourly", "0 * * * *", "Run once an hour at the beginning of the hour"},
	{"@every_minute", "* * * * *", "Run every minute"},
}

// aliasIndex is a lookup map built from predefinedAliases.
var aliasIndex map[string]Alias

func init() {
	aliasIndex = make(map[string]Alias, len(predefinedAliases))
	for _, a := range predefinedAliases {
		aliasIndex[a.Name] = a
	}
}

// ResolveAlias returns the 5-field cron expression for the given alias name.
// If the name is not a recognised alias the original value is returned
// unchanged together with ok=false.
func ResolveAlias(name string) (expr string, ok bool) {
	if a, found := aliasIndex[name]; found {
		return a.Expression, true
	}
	return name, false
}

// LookupAlias returns the full Alias struct for the given name, or an error
// if the name is not recognised.
func LookupAlias(name string) (Alias, error) {
	if a, found := aliasIndex[name]; found {
		return a, nil
	}
	return Alias{}, fmt.Errorf("unknown alias %q", name)
}

// ListAliases returns all predefined aliases in definition order.
func ListAliases() []Alias {
	result := make([]Alias, len(predefinedAliases))
	copy(result, predefinedAliases)
	return result
}

// AliasForExpr returns the first alias whose expression matches expr, or
// ok=false when no alias matches.
func AliasForExpr(expr string) (Alias, bool) {
	for _, a := range predefinedAliases {
		if a.Expression == expr {
			return a, true
		}
	}
	return Alias{}, false
}
