package main

import (
	"fmt"
	"github.com/jumaniyozov/gdo/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jumaniyozov/gdo/controllers"
	"github.com/jumaniyozov/gdo/templates"
	"github.com/jumaniyozov/gdo/views"
)

func main() {
	PORT := "8001"
	r := chi.NewRouter()
	fs := http.FileServer(http.Dir("static"))
	r.Use(middleware.Logger)

	r.Handle("/static/*", http.StripPrefix("/static/", fs))
	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "layout-page.gohtml", "home.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "layout-page.gohtml", "contact.gohtml"))))

	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "layout-page.gohtml", "faq.gohtml"))))

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	userService := models.UserService{
		DB: db,
	}

	usersC := controllers.Users{
		UserService: &userService,
	}

	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "layout-page.gohtml", "signup.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "layout-page.gohtml", "signin.gohtml"))
	r.Get("/signup", usersC.New)
	r.Get("/signin", usersC.SignIn)
	r.Post("/users", usersC.Create)

	fmt.Printf("Server is running on PORT %s...\n", PORT)

	err = http.ListenAndServe(fmt.Sprintf(":%s", PORT), r)
	if err != nil {
		panic(err)
	}
}
