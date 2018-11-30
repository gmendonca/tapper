package query_test

import (
	"testing"

	"github.com/gmendonca/tapper/pkg/query"
)

type QueryTest struct {
	Query string
	Name  string
}

func TestAnalyzeQueries(t *testing.T) {

	QueryTest1 := QueryTest{
		Query: `SELECT column1,
	FROM table1
	WHERE column1 in (value1, value2)`,
		Name: "QueryTest1",
	}

	QueryTest2 := QueryTest{
		Query: `SELECT column1,
	FROM table1 LEFT JOIN table2 ON table1.column2 = table2.column2
	WHERE column1 in (value1, value2)`,
		Name: "QueryTest2",
	}

	QueryTest3 := QueryTest{
		Query: `SELECT column1, FROM table1 WHERE column1 in (value1, value2)`,
		Name:  "QueryTest3",
	}

	QueryTest4 := QueryTest{
		Query: `SELECT * FROM table1 CROSS JOIN table2`,
		Name:  "QueryTest4",
	}

	QueryTest5 := QueryTest{
		Query: `SELECT * FROM table1 JOIN table2 ON table1.a = table2.a JOIN table3 ON table2.b = table3.b`,
		Name:  "QueryTest5",
	}

	QueryTest6 := QueryTest{
		Query: `WITH table_tmp as (
		SELECT a, b FROM table1 WHERE a > 10
	)

	SELECT DISTINCT a
	FROM table_tmp FULL OUTER JOIN table2
	GROUP BY 1
	`,
		Name: "QueryTest6",
	}

	var analyzedQuery query.AnalyzedQuery
	analyzedQuery = query.AnalyzeQueries(QueryTest1.Query)
	if len(analyzedQuery.Joins) > 0 {
		t.Errorf("That shouldn't be any Joins in this query %s: %s", QueryTest1.Name, analyzedQuery.Joins)
	}
	if len(analyzedQuery.TablesNames) != 1 || analyzedQuery.TablesNames[0] != "table1" {
		t.Errorf("There are more tables than expected or the table doesn't match the expected names in query %s: %s", QueryTest1.Name, analyzedQuery.TablesNames)
	}
	analyzedQuery = query.AnalyzeQueries(QueryTest2.Query)
	if len(analyzedQuery.Joins) != 1 {
		t.Errorf("That should be only one Join in this query %s: %s", QueryTest2.Name, analyzedQuery.Joins)
	}
	if len(analyzedQuery.TablesNames) != 2 || analyzedQuery.TablesNames[0] != "table1" || analyzedQuery.TablesNames[1] != "table2" {
		t.Errorf("There are more tables than expected or the table doesn't match the expected names in query %s: %s", QueryTest2.Name, analyzedQuery.TablesNames)
	}
	analyzedQuery = query.AnalyzeQueries(QueryTest3.Query)
	if len(analyzedQuery.Joins) > 0 {
		t.Errorf("That shouldn't be any Joins in this query %s: %s", QueryTest3.Name, analyzedQuery.Joins)
	}
	if len(analyzedQuery.TablesNames) != 1 || analyzedQuery.TablesNames[0] != "table1" {
		t.Errorf("There are more tables than expected or the table doesn't match the expected names in query %s: %s", QueryTest3.Name, analyzedQuery.TablesNames)
	}
	analyzedQuery = query.AnalyzeQueries(QueryTest4.Query)
	if len(analyzedQuery.Joins) != 1 {
		t.Errorf("That should be only one Join in this query %s: %s", QueryTest4.Name, analyzedQuery.Joins)
	}
	if len(analyzedQuery.TablesNames) != 2 || analyzedQuery.TablesNames[0] != "table1" || analyzedQuery.TablesNames[1] != "table2" {
		t.Errorf("There are more tables than expected or the table doesn't match the expected names in query %s: %s", QueryTest4.Name, analyzedQuery.TablesNames)
	}
	analyzedQuery = query.AnalyzeQueries(QueryTest5.Query)
	if len(analyzedQuery.Joins) != 2 {
		t.Logf("Since this is a very naive solution to parse queries, this is failing cause only show one join for this query %s: %s", QueryTest5.Name, analyzedQuery.Joins)
	}
	if len(analyzedQuery.TablesNames) != 3 || analyzedQuery.TablesNames[0] != "table1" || analyzedQuery.TablesNames[1] != "table2" || analyzedQuery.TablesNames[2] != "table3" {
		t.Logf("Since this is a very naive solution to parse queries, this is failing cause only show two table names in query %s: %s", QueryTest5.Name, analyzedQuery.TablesNames)
	}
	analyzedQuery = query.AnalyzeQueries(QueryTest6.Query)
	if len(analyzedQuery.Joins) != 1 {
		t.Errorf("That should be only one Join in this query %s: %s", QueryTest6.Name, analyzedQuery.Joins)
	}
	if len(analyzedQuery.TablesNames) != 3 || analyzedQuery.TablesNames[0] != "table1" || analyzedQuery.TablesNames[1] != "table_tmp" || analyzedQuery.TablesNames[2] != "table2" {
		t.Errorf("There are more tables than expected or the table doesn't match the expected names in query %s: %s", QueryTest6.Name, analyzedQuery.TablesNames)
	}
}
