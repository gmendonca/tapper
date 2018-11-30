package query_test

import (
	"testing"

	"github.com/gmendonca/tapper/pkg/query"
)

type QueryTest struct {
	Query string
	Name  string
}

const (
	QueryTest1 = `SELECT column1,
	FROM table1
	WHERE column1 in (value1, value2)`

	QueryTest2 = `SELECT column1,
	FROM table1 LEFT JOIN table2 ON table1.column2 = table2.column2
	WHERE column1 in (value1, value2)`

	QueryTest3 = `SELECT column1, FROM table1 WHERE column1 in (value1, value2)`

	QueryTest4 = `SELECT * FROM table1 CROSS JOIN table2`

	QueryTest5 = `SELECT * FROM table1 JOIN table2 ON table1.a = table2.a JOIN table3 ON table2.b = table3.b`

	QueryTest6 = `WITH table_tmp as (
		SELECT a, b FROM table1 WHERE a > 10
	)

	SELECT DISTINCT a
	FROM table_tmp FULL OUTER JOIN table2
	GROUP BY 1
	`
)

func TestAnalyzeQueries(t *testing.T) {
	var analyzedQuery query.AnalyzedQuery
	analyzedQuery = query.AnalyzeQueries(QueryTest1)
	if len(analyzedQuery.Joins) > 0 {
		t.Errorf("That shouldn't be any Joins in this query: %s", analyzedQuery.Joins)
	}
	if len(analyzedQuery.TablesNames) != 1 || analyzedQuery.TablesNames[0] != "table1" {
		t.Errorf("There are more tables than expected or the table doesn't match table1: %s", analyzedQuery.TablesNames[0])
	}
	analyzedQuery = query.AnalyzeQueries(QueryTest2)
	if len(analyzedQuery.Joins) != 1 {
		t.Errorf("That should be only one Join in this query: %s", analyzedQuery.Joins)
	}
	if len(analyzedQuery.TablesNames) != 2 || analyzedQuery.TablesNames[0] != "table1" || analyzedQuery.TablesNames[1] != "table2" {
		t.Errorf("There are more tables than expected or the table doesn't match the expected names: %s", analyzedQuery.TablesNames)
	}
	analyzedQuery = query.AnalyzeQueries(QueryTest3)
	if len(analyzedQuery.Joins) > 0 {
		t.Errorf("That shouldn't be any Joins in this query: %s", analyzedQuery.Joins)
	}
	if len(analyzedQuery.TablesNames) != 1 || analyzedQuery.TablesNames[0] != "table1" {
		t.Errorf("There are more tables than expected or the table doesn't match table1: %s", analyzedQuery.TablesNames[0])
	}
	analyzedQuery = query.AnalyzeQueries(QueryTest4)
	if len(analyzedQuery.Joins) != 1 {
		t.Errorf("That should be only one Join in this query: %s", analyzedQuery.Joins)
	}
	if len(analyzedQuery.TablesNames) != 2 || analyzedQuery.TablesNames[0] != "table1" || analyzedQuery.TablesNames[1] != "table2" {
		t.Errorf("There are more tables than expected or the table doesn't match the expected names: %s", analyzedQuery.TablesNames)
	}
	analyzedQuery = query.AnalyzeQueries(QueryTest5)
	if len(analyzedQuery.Joins) != 2 {
		t.Logf("Since this is a very naive solution to parse queries, this is failing cause only show one join for this query: %s", analyzedQuery.Joins)
	}
	if len(analyzedQuery.TablesNames) != 3 || analyzedQuery.TablesNames[0] != "table1" || analyzedQuery.TablesNames[1] != "table2" || analyzedQuery.TablesNames[2] != "table3" {
		t.Logf("Since this is a very naive solution to parse queries, this is failing cause only show two table names: %s", analyzedQuery.TablesNames)
	}
	analyzedQuery = query.AnalyzeQueries(QueryTest6)
	if len(analyzedQuery.Joins) != 1 {
		t.Errorf("That should be only one Join in this query: %s", analyzedQuery.Joins)
	}
	if len(analyzedQuery.TablesNames) != 3 || analyzedQuery.TablesNames[0] != "table1" || analyzedQuery.TablesNames[1] != "table_tmp" || analyzedQuery.TablesNames[2] != "table2" {
		t.Errorf("There are more tables than expected or the table doesn't match the expected names: %s", analyzedQuery.TablesNames)
	}
}
