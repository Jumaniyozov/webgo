package views

import (
	"bytes"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/jumaniyozov/gdo/context"
	"github.com/jumaniyozov/gdo/models"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
)

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl := template.New(patterns[0])
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField not implemented")
			},
			"currentUser": func() (template.HTML, error) {
				return "", fmt.Errorf("current user not implemented")
			},
		},
	)

	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	return Template{HTMLTpl: tpl}, nil
}

//func Parse(filepath string) (Template, error) {
//	tpl, err := template.ParseFiles(filepath)
//	if err != nil {
//		return Template{}, fmt.Errorf("parsing embbeded file: %w", err)
//	}
//
//	return Template{HTMLTpl: tpl}, nil
//}

type Template struct {
	HTMLTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data any) {
	tpl, err := t.HTMLTpl.Clone()
	if err != nil {
		log.Printf("error: %w", err)
		http.Error(w, "There was an error rendering the page", http.StatusInternalServerError)
		return
	}

	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
		},
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error parsing the template", http.StatusInternalServerError)
		return
	}

	io.Copy(w, &buf)
}
