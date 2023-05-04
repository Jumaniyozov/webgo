package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Address struct {
	Street       string
	Country      string
	StreetNumber int
}

type User struct {
	Name    string
	Address Address
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	user := User{
		Name: "Islom",
		Address: Address{
			Street:       "Kamarniso",
			Country:      "Uzbekistan",
			StreetNumber: 6,
		},
	}
	tplPath := filepath.Join("templates", "hello.gohtml")
	executeTemplate(w, tplPath, user)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "contact.gohtml")
	executeTemplate(w, tplPath, nil)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "faq.gohtml")
	executeTemplate(w, tplPath, nil)
}

func main() {
	PORT := "8001"

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)

	fmt.Printf("Server is running on PORT %s...\n", PORT)

	err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), r)
	if err != nil {
		panic(err)
	}
}

func executeTemplate(w http.ResponseWriter, filePath string, fields any) {
	t, err := template.ParseFiles(filePath)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "There was an error parsing the template", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, fields)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}
}
