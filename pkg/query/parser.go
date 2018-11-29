package query

import "strings"

//AnalyzeQueries gets the most used tables and most used joins to send to Datadog
func AnalyzeQueries(query string) []string {
	query = strings.ToUpper(query)
	parts := strings.Fields(query)

	var blacklistedTables []string
	var tableNames []string
	var joins []string

	fromStatement := false

	for idx, part := range parts {
		if fromStatement {
			// Here we are sure that's a FROM statement starting
			// Now we should check if this is not the beginning of a new query
			if part != "(" && part != "(SELECT" {
				//Now we should check if w
				//Then this should be a table name
				if !tableNames.Contains(part) {
					tableNames = append(tableNames, part)
				}
			}
		}

		if part == "FROM" {
			// Starting FROM Statment
			fromStatement = true
		}
	}

	return parts
}
