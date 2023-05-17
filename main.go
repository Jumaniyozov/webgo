package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
	"github.com/jumaniyozov/gdo/controllers"
	"github.com/jumaniyozov/gdo/migrations"
	"github.com/jumaniyozov/gdo/models"
	"github.com/jumaniyozov/gdo/templates"
	"github.com/jumaniyozov/gdo/views"
	"net/http"
	"os"
	"strconv"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

const (
	PORT = "8001"
)

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	cfg.PSQL = models.DefaultPostgresConfig()

	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	cfg.CSRF.Key = "ghjkghjkghjkghjkghjkghjkghjkghjk"
	cfg.CSRF.Secure = false

	cfg.Server.Address = ":8001"

	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	userService := &models.UserService{DB: db}
	sessionService := &models.SessionService{DB: db}
	pwResetService := &models.PasswordResetService{DB: db}
	emailService := models.NewEmailService(cfg.SMTP)

	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	csrfMw := csrf.Protect([]byte(cfg.CSRF.Key), csrf.Secure(cfg.CSRF.Secure))

	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
	}

	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "layout-page.gohtml", "signup.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "layout-page.gohtml", "signin.gohtml"))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "layout-page.gohtml", "forgot-pw.gohtml"))
	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(templates.FS, "layout-page.gohtml", "check-your-email.gohtml"))
	usersC.Templates.ResetPassword = views.Must(views.ParseFS(templates.FS, "layout-page.gohtml", "reset-pw.gohtml"))

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
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)
	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProccessResetPassword)
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
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}
}
