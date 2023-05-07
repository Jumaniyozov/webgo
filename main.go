package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jumaniyozov/gdo/controllers"
	"github.com/jumaniyozov/gdo/templates"
	"github.com/jumaniyozov/gdo/views"
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

func main() {
	PORT := "8000"
	r := chi.NewRouter()
	fs := http.FileServer(http.Dir("static"))
	r.Use(middleware.Logger)

	r.Handle("/static/*", http.StripPrefix("/static/", fs))
	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "layout-page.gohtml", "home.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "layout-page.gohtml", "contact.gohtml"))))

	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "layout-page.gohtml", "faq.gohtml"))))

	fmt.Printf("Server is running on PORT %s...\n", PORT)

	err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), r)
	if err != nil {
		panic(err)
	}
}
