package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const query = `
SELECT
	log.time,
	address.internetip, 
	device.serial, 
	device.hostname, 
	user.username, 
	user.fullname 
FROM log 
	JOIN identity 
		ON log.identity_id = identity.id 
	JOIN address 
		ON identity.address_id = address.id 
	JOIN device 
		ON identity.device_id = device.id 
	JOIN user 
		ON identity.user_id = user.id 
WHERE %s = ?
ORDER BY log.time DESC 
LIMIT %d,%d;
`

var allowedColumns = map[string]struct{}{
	"address.internetip": {},
	"device.serial":      {},
	"device.hostname":    {},
	"user.username":      {},
}

const minPageSize = 10
const maxPageSize = 100

type ValidationError struct {
	Value  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf(`could not validate "%s": %s`, e.Value, e.Reason)
}

type DB struct {
	*sql.DB
}

func NewDB(dsn string) (*DB, error) {
	if !strings.Contains(dsn, "parseTime=true") {
		return nil, errors.New(`DSN must contain "parseTime=true"`)
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not open database: %w", err)
	}
	return &DB{DB: db}, nil
}

type Row struct {
	Time     time.Time
	IP       string
	Serial   string
	Hostname string
	Username string
	Name     string
}

func (db *DB) Query(col, search string, page, pageSize int) ([]*Row, error) {
	if _, ok := allowedColumns[col]; !ok {
		return nil, &ValidationError{Value: col, Reason: "not a valid column"}
	}

	if search == "" {
		return nil, &ValidationError{Reason: "search must not be empty"}
	}

	if page < 0 {
		return nil, &ValidationError{Value: strconv.Itoa(page), Reason: "not a valid page"}
	}

	if pageSize < minPageSize || pageSize > maxPageSize {
		return nil, &ValidationError{Value: strconv.Itoa(pageSize), Reason: "not a valid page size"}
	}

	offset := page * pageSize
	q := fmt.Sprintf(query, col, offset, pageSize)

	srows, err := db.DB.Query(q, search)
	if err != nil {
		return nil, fmt.Errorf("could not query database: %w", err)
	}
	defer srows.Close()

	var rows []*Row
	for srows.Next() {
		r := new(Row)
		if err := srows.Scan(&r.Time, &r.IP, &r.Serial, &r.Hostname, &r.Username, &r.Name); err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}
		rows = append(rows, r)
	}

	if err = srows.Err(); err != nil {
		return nil, fmt.Errorf("could not read all rows: %w", err)
	}

	return rows, nil
}
