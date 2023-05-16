package main

import (
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/jumaniyozov/gdo/migrations"
	"github.com/jumaniyozov/gdo/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jumaniyozov/gdo/controllers"
	"github.com/jumaniyozov/gdo/templates"
	"github.com/jumaniyozov/gdo/views"
)

const (
	PORT = "8001"
)

func main() {
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	userService := models.UserService{DB: db}
	sessionService := models.SessionService{DB: db}

	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	csrfKey := "ghjkghjkghjkghjkghjkghjkghjkghjk"
	csrfMw := csrf.Protect([]byte(csrfKey), csrf.Secure(false))

	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}

	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "layout-page.gohtml", "signup.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "layout-page.gohtml", "signin.gohtml"))

	r := chi.NewRouter()
	fs := http.FileServer(http.Dir("static"))
	r.Use(csrfMw)
	r.Use(umw.SetUser)

	r.Handle("/static/*", http.StripPrefix("/static/", fs))
	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "layout-page.gohtml", "home.gohtml"))))
	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "layout-page.gohtml", "contact.gohtml"))))
	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "layout-page.gohtml", "faq.gohtml"))))
	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
		r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello")
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Printf("Server is running on PORT %s...\n", PORT)
	err = http.ListenAndServe(fmt.Sprintf(":%s", PORT), r)
	if err != nil {
		panic(err)
	}
}
