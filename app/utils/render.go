package utils

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type Props map[string]interface{}

func (props Props) ValidatePresence(field string) bool {
	switch value := props[field].(type) {
	case string:
		if value == "" {
			props.SetError(field, fmt.Sprintf("%s can't be blank", field))
			return false
		}
	case float64:
		if value == 0.0 {
			props.SetError(field, fmt.Sprintf("%s can't be blank", field))
			return false
		}
	}
	return true
}

func (props Props) ValidateNoSpace(field string) bool {
	if StripSpaces(props[field].(string)) != props[field] {
		props.SetError(field, fmt.Sprintf("%s may not contain spaces", field))
		return false
	}
	return true
}

func (props Props) ValidateConfirmation(field string, confirmationField string) bool {
	if props[field] != props[confirmationField] {
		props.SetError(confirmationField, fmt.Sprintf("%s and %s must match", field, confirmationField))
		return false
	}
	return true
}

func (props Props) ValidateEmail(field string) bool {
	if !(strings.Contains(props[field].(string), "@")) {
		props.SetError(field, fmt.Sprintf("%s must be an email", field))
		return false
	}
	return true
}

func (props Props) ValidateLength(field string, min int, max int) bool {
	length := len(props[field].(string))
	if length < min || length > max {
		props.SetError(field, fmt.Sprintf("%s must be between %d and %d characters in length", field, min, max))
		return false
	}
	return true
}

func (props Props) FieldIsValid(field string) bool {
	return props["errors"].(map[string]string)[field] == ""
}

func (props Props) IsValid() bool {
	return len(props["errors"].(map[string]string)) == 0
}

func (props Props) SetError(field string, value string) {
	props["errors"].(map[string]string)[field] = value
}

var basePath = os.Getenv("base_path")

func formatDate(date time.Time) string {
	return date.Format("02/01/2006")
}

var funcMap = template.FuncMap{
	"formatDate": formatDate,
}

func Render(w http.ResponseWriter, r *http.Request, filename string, props interface{}) {
	tmpl := template.Must(template.New("base").Funcs(funcMap).ParseFiles(path.Join(basePath, "./app/views/layout.html"), path.Join(basePath, "app/views", filename)))

	endProps := make(map[string]interface{})
	for k, v := range *props.(*Props) {
		endProps[k] = v
	}

	data, ok := r.Context().Value("data").(*Props)

	if ok {
		for k, v := range *data {
			endProps[k] = v
		}
	}

	if err := tmpl.ExecuteTemplate(w, "layout", endProps); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
