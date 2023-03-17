package main

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

//go:embed ui.tmpl
var tmpl embed.FS

// Tmpl is the embedded ui template
var Tmpl = template.Must(template.ParseFS(tmpl, "ui.tmpl"))

// input names
const (
	InputNameField          = "field"
	InputNameSearch         = "search"
	InputNamePreviousSearch = "previousSearch"
	InputNamePage           = "page"
	InputNamePageSize       = "pageSize"
)

// Form is the data from a submitted form
type Form struct {
	Field          string
	Search         string
	PreviousSearch string
	Page           int
	PageSize       int
	Rows           []*Row
	Error          error
}

// HandleUI is the HTTP handler for the UI
func (db *DB) HandleUI(r *http.Request) (*Form, error) {
	form := &Form{Field: "user.username", Page: 1, PageSize: 10}
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			return nil, fmt.Errorf("could not parse form: %w", err)
		}
		if f := r.PostFormValue(InputNameField); f != "" {
			form.Field = f
		}
		if s := r.PostFormValue(InputNameSearch); s != "" {
			form.Search = s
		}
		if s := r.PostFormValue(InputNamePreviousSearch); s != "" {
			form.PreviousSearch = s
		}
		if p := r.PostFormValue(InputNamePage); p != "" {
			i, err := strconv.Atoi(p)
			if err == nil {
				form.Page = i
			}
		}
		if p := r.PostFormValue(InputNamePageSize); p != "" {
			i, err := strconv.Atoi(p)
			if err == nil {
				form.PageSize = i
			}
		}

		if form.Search != form.PreviousSearch {
			form.Page = 1
		}
		form.PreviousSearch = form.Search

		rows, err := db.Query(form.Field, form.Search, form.Page-1, form.PageSize)
		if err != nil {
			return nil, fmt.Errorf("could not query database: %w", err)
		}
		if len(rows) == 0 {
			form.Error = errors.New("no rows found")
		}
		form.Rows = rows
	}

	return form, nil
}
