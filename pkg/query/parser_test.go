package query_test

import (
	"testing"

	"github.com/gmendonca/tapper/pkg/query"
)

const (
	QueryTest1 = `SELECT column1,
	FROM table1
	WHERE column1 in (value1, value2)`

	QueryTest2 = `SELECT column1,
	FROM table1 LEFT JOIN table2 ON table1.column2 = table2.column2
	WHERE column1 in (value1, value2)`

	QueryTest3 = `SELECT column1,\nFROM table1\nWHERE column1 in (value1, value2)`
)

func TestAnalyzeQueries(t *testing.T) {
	var parts []string
	parts = query.AnalyzeQueries(QueryTest1)
	t.Log(parts)
	parts = query.AnalyzeQueries(QueryTest2)
	t.Log(parts)
	parts = query.AnalyzeQueries(QueryTest3)
	t.Log(parts)
}
