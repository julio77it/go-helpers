package helpers

import (
	"database/sql"
	"errors"
)

// NewSQLRows : build a SQLRows struct from a database/sql.Rows
func NewSQLRows(rows *sql.Rows) (*SQLRows, error) {
	ct, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	rh := &SQLRows{
		rows:        rows,
		columnTypes: ct,
		columnBytes: make([]interface{}, len(ct)),
	}
	return rh, nil
}

// SQLRows : holds the info about sql.Row fields
type SQLRows struct {
	rows        *sql.Rows
	columnTypes []*sql.ColumnType
	columnBytes []interface{}
}

// Next : prepares the next result row for reading with the Fetch method
func (rh SQLRows) Next() bool {
	return rh.rows.Next()
}

// Length : number of fields of the result set
func (rh SQLRows) Length() int {
	return len(rh.columnTypes)
}

// Fetch : read the bytes of the current row
func (rh *SQLRows) Fetch() error {
	// reserve memory space
	for i := 0; i < len(rh.columnBytes); i++ {
		rh.columnBytes[i] = new(sql.RawBytes)
	}
	// read as variadic parameters
	err := rh.rows.Scan(rh.columnBytes...)
	return err
}

// GetFieldByIndex : find a field By index. Return name, value and error
func (rh SQLRows) GetFieldByIndex(index int) (string, interface{}, error) {
	// Check the input parameters
	if index < 0 || index >= len(rh.columnTypes) {
		// return zerov, zerov, error
		return "", nil, errors.New("index out of bound")
	}
	rb := rh.columnBytes[index].(*sql.RawBytes)
	// convert bytes in string
	// TODO convert in the right value
	// return name, value, zerov
	return rh.columnTypes[index].Name(), string(*rb), nil
}

// GetFieldByName : find a field By name. Returns index, value and error
func (rh SQLRows) GetFieldByName(name string) (int, interface{}, error) {
	err := errors.New("field not found")

	for i, v := range rh.columnTypes {
		// Find the index of the field name
		if name == v.Name() {
			// Get the field by index
			_, value, err := rh.GetFieldByIndex(i)
			if err != nil {
				break
			}
			// return index, value, zerov
			return i, value, nil
		}
	}
	// return zerov, zerov, nil
	return 0, nil, err
}
