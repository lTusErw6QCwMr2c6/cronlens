package schedule

import (
	"fmt"
	"strings"
)

// FormatTemplate returns a human-readable string for a single Template.
func FormatTemplate(t Template) string {
	return fmt.Sprintf("%-20s %-20s %s", t.Name, t.Expr, t.Description)
}

// FormatTemplateList formats a slice of templates as a table.
func FormatTemplateList(templates []Template) string {
	if len(templates) == 0 {
		return "No templates available."
	}
	var sb strings.Builder
	header := fmt.Sprintf("%-20s %-20s %s", "NAME", "EXPRESSION", "DESCRIPTION")
	sb.WriteString(header)
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("-", len(header)))
	sb.WriteString("\n")
	for _, t := range templates {
		sb.WriteString(FormatTemplate(t))
		sb.WriteString("\n")
	}
	return sb.String()
}

// TemplateSummaryLine returns a one-line summary for a template lookup result.
func TemplateSummaryLine(name string, found bool) string {
	if !found {
		return fmt.Sprintf("Template %q not found.", name)
	}
	return fmt.Sprintf("Template %q found.", name)
}
