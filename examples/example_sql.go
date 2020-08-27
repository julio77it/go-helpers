package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/julio77it/go-helpers/helpers"

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

func main() {
	// datasource : for sqlite3 is a temporary file
	dataSource, err := ioutil.TempFile("/var/tmp", "helpers.go")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(dataSource.Name())

	// create new data database
	db, err := sql.Open(driver, dataSource.Name())
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
		fmt.Printf("INSERT (id=%d,author=%s,quote=%s)\n", i, q.QuoteAuthor, q.QuoteText)
		if _, err := db.Exec(insertStmst, i, q.QuoteAuthor, q.QuoteText); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}

	// query the table
	rows, err := db.Query(queryStmt)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer rows.Close()

	// the core of the example starts here
	// helpers/sql.SQLRows : get the fields info from the resultset
	sqlRows, err := helpers.NewSQLRows(rows)

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	for ri := 0; sqlRows.Next(); ri++ {
		// Every row need the call of SQLRows.Fetch method
		// it reads the fields bytes
		if err := sqlRows.Fetch(); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Printf("ROW[%2d]\n", ri)

		for fi := 0; fi < sqlRows.Length(); fi++ {
			name, value, err := sqlRows.GetFieldByIndex(fi)
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
			fmt.Printf("\tFIELD[%2d] %s = %v\n", fi, name, value)
		}

	}
	if err := sqlRows.Err(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
