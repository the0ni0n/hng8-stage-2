package main

import (
	"html/template"
	"net/http"
	"os"
)

var tpl *template.Template

type ContactDetails struct {
	Email   string
	Subject string
	Message string
}

func init() {
	tpl = template.Must(template.ParseGlob("*.html"))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
}

func main() {
	http.HandleFunc("/", index)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	http.ListenAndServe(":"+port, nil)
}

func index(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		tpl.ExecuteTemplate(w, "index.html", nil)
		return
	}

	details := ContactDetails{
		Email:   r.FormValue("email"),
		Subject: r.FormValue("subject"),
		Message: r.FormValue("message"),
	}

	tpl.ExecuteTemplate(w, "contact.html", details)
}
