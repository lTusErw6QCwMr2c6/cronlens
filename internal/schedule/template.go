package schedule

// CommonTemplates holds a registry of well-known cron expressions with labels.
var CommonTemplates = []Template{
	{Name: "every-minute", Expr: "* * * * *", Description: "Every minute"},
	{Name: "every-5-minutes", Expr: "*/5 * * * *", Description: "Every 5 minutes"},
	{Name: "every-15-minutes", Expr: "*/15 * * * *", Description: "Every 15 minutes"},
	{Name: "every-30-minutes", Expr: "*/30 * * * *", Description: "Every 30 minutes"},
	{Name: "hourly", Expr: "0 * * * *", Description: "Every hour at minute 0"},
	{Name: "daily-midnight", Expr: "0 0 * * *", Description: "Daily at midnight"},
	{Name: "daily-noon", Expr: "0 12 * * *", Description: "Daily at noon"},
	{Name: "weekly-sunday", Expr: "0 0 * * 0", Description: "Weekly on Sunday at midnight"},
	{Name: "weekly-monday", Expr: "0 0 * * 1", Description: "Weekly on Monday at midnight"},
	{Name: "monthly-first", Expr: "0 0 1 * *", Description: "Monthly on the 1st at midnight"},
	{Name: "monthly-last", Expr: "0 0 28 * *", Description: "Monthly on the 28th at midnight"},
	{Name: "yearly", Expr: "0 0 1 1 *", Description: "Yearly on January 1st at midnight"},
	{Name: "weekdays", Expr: "0 9 * * 1-5", Description: "Weekdays at 9am"},
	{Name: "weekends", Expr: "0 10 * * 6,0", Description: "Weekends at 10am"},
}

// Template represents a named, reusable cron schedule.
type Template struct {
	Name        string
	Expr        string
	Description string
}

// FindTemplate returns a Template by name, and whether it was found.
func FindTemplate(name string) (Template, bool) {
	for _, t := range CommonTemplates {
		if t.Name == name {
			return t, true
		}
	}
	return Template{}, false
}

// ListTemplates returns all available template names.
func ListTemplates() []string {
	names := make([]string, len(CommonTemplates))
	for i, t := range CommonTemplates {
		names[i] = t.Name
	}
	return names
}

// TemplateByExpr searches for a template matching the given cron expression.
func TemplateByExpr(expr string) (Template, bool) {
	for _, t := range CommonTemplates {
		if t.Expr == expr {
			return t, true
		}
	}
	return Template{}, false
}
