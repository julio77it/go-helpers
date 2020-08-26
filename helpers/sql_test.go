package helpers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const (
	// for a simple example, I'v chosen sqlite3
	driver string = "sqlite3"
	// SQL : table creation
	createStmst string = "CREATE TABLE quotes (id INTEGER, author TEXT, quoteText TEXT)"
	// SQL : insert rows
	insertStmst string = "INSERT INTO quotes (id, author, quoteText) VALUES (?, ?, ?)"
	// SQL : query the table
	queryStmt string = "SELECT * FROM quotes"
	// quotes resource
	quotesURL string = "https://raw.githubusercontent.com/JamesFT/Database-Quotes-JSON/master/quotes.json"
)

type quote struct {
	QuoteText   string `json:"quoteText"`
	QuoteAuthor string `json:"quoteAuthor"`
}

func getQuotes() ([]quote, error) {
	// Read quotes from github
	resp, err := http.Get(quotesURL)
	if err != nil {
		return nil, err
	}
	var quotes []quote
	// convert JSON to a slice of quotes
	err = json.NewDecoder(resp.Body).Decode(&quotes)
	if err != nil {
		return nil, err
	}
	return quotes, nil
}

var db *sql.DB

func TestMain(m *testing.M) {
	// datasource : for sqlite3 is a temporary file
	dataSource, err := ioutil.TempFile("/var/tmp", "helpers.go")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(dataSource.Name())

	// create new data database
	db, err = sql.Open(driver, dataSource.Name())
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer db.Close()

	// check the connection
	if err := db.Ping(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// create table
	if _, err := db.Exec(createStmst); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// get quotes from http
	quotes, err := getQuotes()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// insert rows
	for i, q := range quotes {
		if i == 100 {
			break
		}
		if _, err := db.Exec(insertStmst, i, q.QuoteAuthor, q.QuoteText); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}
	os.Exit(m.Run())
}

func TestNewSQLRowHeaders(t *testing.T) {
	rows, err := db.Query(queryStmt)
	if err != nil {
		t.Errorf("db.Query failed : %v", err)
	}

	// OK
	_, err = NewSQLRowHeaders(rows)
	if err != nil {
		t.Errorf("NewSQLRowHeaders failed : %v", err)
	}
	rows.Close()

	// KO
	_, err = NewSQLRowHeaders(rows)
	if err == nil {
		t.Errorf("NewSQLRowHeaders error expected, got %v", err)
	}
}

func TestLength(t *testing.T) {
	rows, err := db.Query(queryStmt)
	if err != nil {
		t.Errorf("db.Query failed : %v", err)
	}
	defer rows.Close()
	// OK
	rh, err := NewSQLRowHeaders(rows)
	if err != nil {
		t.Errorf("NewSQLRowHeaders failed : %v", err)
	}
	length := rh.Length()
	if length != 3 {
		t.Errorf("SQLRowHeaders.Length expected 3, got %d", length)
	}
}

func TestFetch(t *testing.T) {
	rows, err := db.Query(queryStmt)
	if err != nil {
		t.Errorf("db.Query failed : got %v", err)
	}
	rh, err := NewSQLRowHeaders(rows)
	if err != nil {
		t.Errorf("NewSQLRowHeaders failed : got %v", err)
	}
	rows.Next()
	if err := rh.Fetch(); err != nil {
		t.Errorf("SQLRowHeaders.Fetch failed : got %v", err)
	}
	rows.Close()

	if err := rh.Fetch(); err == nil {
		t.Errorf("SQLRowHeaders.Fetch not failed : error expected")
	}
}
func TestGetFieldByIndex(t *testing.T) {
	rows, err := db.Query(queryStmt)
	if err != nil {
		t.Errorf("db.Query failed : got %v", err)
	}
	defer rows.Close()
	rh, err := NewSQLRowHeaders(rows)
	if err != nil {
		t.Errorf("NewSQLRowHeaders failed : got %v", err)
	}
	rows.Next()
	if err := rh.Fetch(); err != nil {
		t.Errorf("SQLRowHeaders.Fetch failed : got %v", err)
	}
	if _, _, err := rh.GetFieldByIndex(-1); err == nil {
		t.Errorf("SQLRowHeaders.Fetch not failed : error expected")
	}
	if _, _, err := rh.GetFieldByIndex(1); err != nil {
		t.Errorf("SQLRowHeaders.Fetch failed : got %v", err)
	}
}
func TestGetFieldByName(t *testing.T) {
	rows, err := db.Query(queryStmt)
	if err != nil {
		t.Errorf("db.Query failed : got %v", err)
	}
	defer rows.Close()
	rh, err := NewSQLRowHeaders(rows)
	if err != nil {
		t.Errorf("NewSQLRowHeaders failed : got %v", err)
	}
	rows.Next()
	if err := rh.Fetch(); err != nil {
		t.Errorf("SQLRowHeaders.Fetch failed : got %v", err)
	}
	if _, _, err := rh.GetFieldByName("BOH"); err == nil {
		t.Errorf("SQLRowHeaders.Fetch not failed : error expected")
	}
	if _, _, err := rh.GetFieldByName("author"); err != nil {
		t.Errorf("SQLRowHeaders.Fetch failed : got %v", err)
	}
}
