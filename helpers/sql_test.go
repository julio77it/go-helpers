package helpers_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/julio77it/go-helpers/helpers"
	_ "github.com/mattn/go-sqlite3"
)

const (
	// SQL : query the table
	queryStmt string = "SELECT * FROM quotes LIMIT 50"
)

func TestNewSQLRows(t *testing.T) {
	// open database
	db, err := sql.Open("sqlite3", "sql_test.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer db.Close()

	// check the connection
	if err = db.Ping(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	rows, err := db.Query(queryStmt)
	if err != nil {
		t.Errorf("db.Query failed : %v", err)
	}

	// OK
	_, err = helpers.NewSQLRows(rows)
	if err != nil {
		t.Errorf("NewSQLRows failed : %v", err)
	}
	rows.Close()

	// KO
	_, err = helpers.NewSQLRows(rows)
	if err == nil {
		t.Errorf("NewSQLRows error expected, got %v", err)
	}
}

func TestLength(t *testing.T) {
	// open database
	db, err := sql.Open("sqlite3", "sql_test.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer db.Close()

	// check the connection
	if err = db.Ping(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	rows, err := db.Query(queryStmt)
	if err != nil {
		t.Errorf("db.Query failed : %v", err)
	}
	defer rows.Close()
	// OK
	rh, err := helpers.NewSQLRows(rows)
	if err != nil {
		t.Errorf("NewSQLRows failed : %v", err)
	}
	length := rh.Length()
	if length != 3 {
		t.Errorf("SQLRows.Length expected 3, got %d", length)
	}
}

func TestFetch(t *testing.T) {
	// open database
	db, err := sql.Open("sqlite3", "sql_test.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer db.Close()

	// check the connection
	if err = db.Ping(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	rows, err := db.Query(queryStmt)
	if err != nil {
		t.Errorf("db.Query failed : got %v", err)
	}
	rh, err := helpers.NewSQLRows(rows)
	if err != nil {
		t.Errorf("NewSQLRows failed : got %v", err)
	}
	rows.Next()
	if err := rh.Fetch(); err != nil {
		t.Errorf("SQLRows.Fetch failed : got %v", err)
	}
	rows.Close()

	if err := rh.Fetch(); err == nil {
		t.Errorf("SQLRows.Fetch not failed : error expected")
	}
}

func TestGetFields(t *testing.T) {
	// open database
	db, err := sql.Open("sqlite3", "sql_test.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer db.Close()

	// check the connection
	if err = db.Ping(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	rows, err := db.Query(queryStmt)
	if err != nil {
		t.Errorf("db.Query failed : got %v", err)
	}
	defer rows.Close()
	rh, err := helpers.NewSQLRows(rows)
	if err != nil {
		t.Errorf("NewSQLRows failed : got %v", err)
	}
	rh.Next()
	if err := rh.Fetch(); err != nil {
		t.Errorf("SQLRows.Fetch failed : got %v", err)
	}
	row := rh.GetFields()

	if row["author"] == "" {
		t.Errorf("SQLRows.Get failed expected field: got empty %v", row)
	}
}

func TestGetFieldByIndex(t *testing.T) {
	// open database
	db, err := sql.Open("sqlite3", "sql_test.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer db.Close()

	// check the connection
	if err = db.Ping(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	rows, err := db.Query(queryStmt)
	if err != nil {
		t.Errorf("db.Query failed : got %v", err)
	}
	defer rows.Close()
	rh, err := helpers.NewSQLRows(rows)
	if err != nil {
		t.Errorf("NewSQLRows failed : got %v", err)
	}
	rh.Next()
	if err := rh.Fetch(); err != nil {
		t.Errorf("SQLRows.Fetch failed : got %v", err)
	}
	if _, _, err := rh.GetFieldByIndex(-1); err == nil {
		t.Errorf("SQLRows.Fetch not failed : error expected")
	}
	if _, _, err := rh.GetFieldByIndex(1); err != nil {
		t.Errorf("SQLRows.Fetch failed : got %v", err)
	}
}

func TestGetFieldByName(t *testing.T) {
	// open database
	db, err := sql.Open("sqlite3", "sql_test.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer db.Close()

	// check the connection
	if err = db.Ping(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	rows, err := db.Query(queryStmt)
	if err != nil {
		t.Errorf("db.Query failed : got %v", err)
	}
	defer rows.Close()
	rh, err := helpers.NewSQLRows(rows)
	if err != nil {
		t.Errorf("NewSQLRows failed : got %v", err)
	}
	rows.Next()
	if err := rh.Fetch(); err != nil {
		t.Errorf("SQLRows.Fetch failed : got %v", err)
	}
	if _, _, err := rh.GetFieldByName("BOH"); err == nil {
		t.Errorf("SQLRows.Fetch not failed : error expected")
	}
	if _, _, err := rh.GetFieldByName("author"); err != nil {
		t.Errorf("SQLRows.Fetch failed : got %v", err)
	}
}

// BenchmarkRows : test database/sql.Rows
func BenchmarkRows(b *testing.B) {
	// open database
	db, err := sql.Open("sqlite3", "sql_test.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer db.Close()

	// check the connection
	if err = db.Ping(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	b.ResetTimer()

	var author, quote string
	for n := 0; n < b.N; n++ {
		rows, err := db.Query(queryStmt)
		if err != nil {
			b.Errorf("db.Query failed : got %v", err)
		}
		for rows.Next() {
			rows.Scan(author, quote)
		}
		rows.Close()
	}
}

// BenchmarkSQLRowsGet : test helpers.SQLRows.Get
func BenchmarkSQLRowsGetFields(b *testing.B) {
	// open database
	db, err := sql.Open("sqlite3", "sql_test.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer db.Close()

	// check the connection
	if err = db.Ping(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		rows, err := db.Query(queryStmt)
		if err != nil {
			b.Errorf("db.Query failed : got %v", err)
		}
		rh, err := helpers.NewSQLRows(rows)
		if err != nil {
			b.Errorf("NewSQLRows failed : got %v", err)
		}
		for rh.Next() {
			rh.Fetch()
			rh.GetFields()
		}
		rows.Close()
	}
}

// BenchmarkSQLRowsGetByIndex : test helpers.SQLRows.GetFieldByIndex
func BenchmarkSQLRowsGetByIndex(b *testing.B) {
	// open database
	db, err := sql.Open("sqlite3", "sql_test.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer db.Close()

	// check the connection
	if err = db.Ping(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		rows, err := db.Query(queryStmt)
		if err != nil {
			b.Errorf("db.Query failed : got %v", err)
		}
		rh, err := helpers.NewSQLRows(rows)
		if err != nil {
			b.Errorf("NewSQLRows failed : got %v", err)
		}
		for rh.Next() {
			rh.Fetch()
			rh.GetFieldByIndex(0)
			rh.GetFieldByIndex(1)
		}
		rows.Close()
	}
}

// BenchmarkSQLRowsGetByName : test helpers.SQLRows.GetFieldByName
func BenchmarkSQLRowsGetByName(b *testing.B) {
	// open database
	db, err := sql.Open("sqlite3", "sql_test.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer db.Close()

	// check the connection
	if err = db.Ping(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		rows, err := db.Query(queryStmt)
		if err != nil {
			b.Errorf("db.Query failed : got %v", err)
		}
		rh, err := helpers.NewSQLRows(rows)
		if err != nil {
			b.Errorf("NewSQLRows failed : got %v", err)
		}
		for rh.Next() {
			rh.Fetch()
			rh.GetFieldByName("author")
			rh.GetFieldByName("quoteText")
		}
		rows.Close()
	}
}

// BenchmarkSQLRowsGetAllByIndex : test helpers.SQLRows.GetFieldByIndex
func BenchmarkSQLRowsGetAllByIndex(b *testing.B) {
	// open database
	db, err := sql.Open("sqlite3", "sql_test.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer db.Close()

	// check the connection
	if err = db.Ping(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		rows, err := db.Query(queryStmt)
		if err != nil {
			b.Errorf("db.Query failed : got %v", err)
		}
		rh, err := helpers.NewSQLRows(rows)
		if err != nil {
			b.Errorf("NewSQLRows failed : got %v", err)
		}
		for rh.Next() {
			rh.Fetch()
			for i := 0; i < rh.Length(); i++ {
				rh.GetFieldByIndex(i)
			}
		}
		rows.Close()
	}
}

func ExampleSQLRows() {
	// open database
	db, err := sql.Open("sqlite3", "sql_test.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	// query
	rows, err := db.Query("SELECT * FROM quotes LIMIT 1")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	// promote sql.Rows in helpers.SQLRows
	rh, err := helpers.NewSQLRows(rows)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 1st row
	if rh.Next() {
		// read fields from bytes
		rh.Fetch()
		// get fields
		row := rh.GetFields()

		fmt.Println(len(row))
	}
	// Output:
	// 3
}
