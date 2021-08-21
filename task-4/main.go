package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	gomail "gopkg.in/mail.v2"
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

	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	http.HandleFunc("/", index)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
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

	err := sendContact(details)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tpl.ExecuteTemplate(w, "contact.html", details)
}

//sendContact sends FormValue to mail.
func sendContact(details ContactDetails) error {

	//Get credentials for smtp.gmail.com
	mail := os.Getenv("GMAIL")
	password := os.Getenv("PASSWORD")

	m := gomail.NewMessage()

	m.SetAddressHeader("From", details.Email, details.Email)
	m.SetHeader("Reply-To", details.Email)
	m.SetHeader("To", mail)
	m.SetHeader("Subject", details.Subject)
	m.SetBody("text/plain", details.Message)

	d := gomail.NewDialer("smtp.gmail.com", 465, mail, password)

	//Now Send
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
