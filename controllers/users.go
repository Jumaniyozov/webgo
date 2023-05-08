package controllers

import (
	"fmt"
	"github.com/jumaniyozov/gdo/views"
	"net/http"
)

type Users struct {
	Templates struct {
		New views.Template
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")

	u.Templates.New.Execute(w, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	_, err = fmt.Fprintln(w, "Email: ", r.FormValue("email"))
	if err != nil {
		fmt.Printf("Error printing into writer %s \n", err)
	}
	_, err = fmt.Fprintln(w, "Password: ", r.FormValue("password"))
	if err != nil {
		fmt.Printf("Error printing into writer %s \n", err)
	}
}
