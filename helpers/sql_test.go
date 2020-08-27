package helpers

import (
	"database/sql"
	"fmt"
	"os"
	"runtime"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const (
	// SQL : query the table
	queryStmt string = "SELECT * FROM quotes"
)

var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	// open database
	db, err = sql.Open("sqlite3", "sql_test.db")
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
	// start test session
	os.Exit(m.Run())
}

func TestNewSQLRows(t *testing.T) {

	rows, err := db.Query(queryStmt)
	if err != nil {
		t.Errorf("db.Query failed : %v", err)
	}
	fmt.Println(runtime.Caller(1))

	// OK
	_, err = NewSQLRows(rows)
	if err != nil {
		t.Errorf("NewSQLRows failed : %v", err)
	}
	rows.Close()

	// KO
	_, err = NewSQLRows(rows)
	if err == nil {
		t.Errorf("NewSQLRows error expected, got %v", err)
	}
}

func TestErr(t *testing.T) {
	rows, err := db.Query(queryStmt)
	if err != nil {
		t.Errorf("db.Query failed : %v", err)
	}
	defer rows.Close()

	rh, err := NewSQLRows(rows)
	if err != nil {
		t.Errorf("NewSQLRows failed : %v", err)
	}
	// OK
	if err = rh.Err(); err != nil {
		t.Errorf("SQLRows.Err expected nil, got %v", err)
	}
}

func TestNext(t *testing.T) {
	rows, err := db.Query(queryStmt)
	if err != nil {
		t.Errorf("db.Query failed : %v", err)
	}
	defer rows.Close()
	// OK
	rh, err := NewSQLRows(rows)
	if err != nil {
		t.Errorf("NewSQLRows failed : %v", err)
	}
	if rh.Next() == false {
		t.Errorf("SQLRows.Length expected true, got false")
	}
}

func TestLength(t *testing.T) {
	rows, err := db.Query(queryStmt)
	if err != nil {
		t.Errorf("db.Query failed : %v", err)
	}
	defer rows.Close()
	// OK
	rh, err := NewSQLRows(rows)
	if err != nil {
		t.Errorf("NewSQLRows failed : %v", err)
	}
	length := rh.Length()
	if length != 3 {
		t.Errorf("SQLRows.Length expected 3, got %d", length)
	}
}

func TestFetch(t *testing.T) {
	rows, err := db.Query(queryStmt)
	if err != nil {
		t.Errorf("db.Query failed : got %v", err)
	}
	rh, err := NewSQLRows(rows)
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

func TestGet(t *testing.T) {
	rows, err := db.Query(queryStmt)
	if err != nil {
		t.Errorf("db.Query failed : got %v", err)
	}
	defer rows.Close()
	rh, err := NewSQLRows(rows)
	if err != nil {
		t.Errorf("NewSQLRows failed : got %v", err)
	}
	rows.Next()
	if err := rh.Fetch(); err != nil {
		t.Errorf("SQLRows.Fetch failed : got %v", err)
	}
	row := rh.Get()

	if row["author"] == "" {
		t.Errorf("SQLRows.Get failed expected field: got empty %v", row)
	}
}

func TestGetFieldByIndex(t *testing.T) {
	rows, err := db.Query(queryStmt)
	if err != nil {
		t.Errorf("db.Query failed : got %v", err)
	}
	defer rows.Close()
	rh, err := NewSQLRows(rows)
	if err != nil {
		t.Errorf("NewSQLRows failed : got %v", err)
	}
	rows.Next()
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
	rows, err := db.Query(queryStmt)
	if err != nil {
		t.Errorf("db.Query failed : got %v", err)
	}
	defer rows.Close()
	rh, err := NewSQLRows(rows)
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
