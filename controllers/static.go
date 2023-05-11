package controllers

import (
	"net/http"
)

func StaticHandler(tpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, nil)
	}
}

func FAQ(tpl Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   string
	}{
		{
			Question: "Is there a free version?",
			Answer:   "Yes ! we offer a fre trial for 30 days on any paid plans.",
		},
		{
			Question: "Is there a free version?",
			Answer:   "Yes ! we offer a fre trial for 30 days on any paid plans.",
		},
		{
			Question: "Is there a free version?",
			Answer:   "Yes ! we offer a fre trial for 30 days on any paid plans.",
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, questions)
	}
}
