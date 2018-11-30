package query

import (
	"strings"
)

type AnalyzedQuery struct {
	TablesNames []string
	Joins       []string
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//AnalyzeQueries gets the most used tables and most used joins to send to Datadog
func AnalyzeQueries(query string) AnalyzedQuery {
	joinWords := []string{
		"LEFT", "left",
		"RIGHT", "right",
		"FULL", "full",
		"INNER", "inner",
		"OUTER", "outer",
		"CROSS", "cross",
		"JOIN", "join",
	}

	tokens := strings.Fields(query)

	var joinToken string
	var tableNames []string
	var joins []string

	joinStatement := false
	fromStatement := false

	for idx, token := range tokens {
		if joinStatement {
			if contains(joinWords, token) {
				joinToken = joinToken + " " + token
			} else {
				//The second table that we're doing the join with
				joinToken = joinToken + " " + token
				joinStatement = false
				fromStatement = false

				if !contains(tableNames, token) {
					tableNames = append(tableNames, token)
				}

				if !contains(joins, joinToken) {
					joins = append(joins, joinToken)
				}
			}
		} else if fromStatement {
			//Here we are sure that this is a from Statement
			if token != "(" && token != "(SELECT" && token != "(select" {
				//Here we are sure that is a table name
				if !contains(tableNames, token) {
					tableNames = append(tableNames, token)
				}
				if idx+1 < len(tokens) && contains(joinWords, tokens[idx+1]) {
					joinStatement = true
					joinToken = token
				} else {
					fromStatement = false
				}
			} else {
				//Since it's a new query we shouldn't care about it anymore
				fromStatement = false
			}
		} else if token == "FROM" || token == "from" {
			fromStatement = true
		}
	}

	a := AnalyzedQuery{
		TablesNames: tableNames,
		Joins:       joins,
	}

	return a
}
