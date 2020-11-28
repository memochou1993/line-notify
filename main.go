package main

import (
	_ "github.com/joho/godotenv/autoload"

	"log"
	"net/http"
	"net/url"
	"os"

	"html/template"
)

var (
	clientID     string
	clientSecret string
	callbackURL  string
	token        string
)

func main() {
	clientID = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	callbackURL = os.Getenv("CALLBACK_URL")

	http.HandleFunc("/callback", callbackHandler)
	http.HandleFunc("/notify", notifyHandler)
	http.HandleFunc("/auth", authHandler)

	log.Fatal(http.ListenAndServe(":85", nil))
}

func notifyHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err.Error())
	}

	msg := r.Form.Get("msg")

	data := url.Values{}
	data.Add("message", msg)

	payload, err := call("POST", host+"/api/notify", data, token)

	if err != nil {
		log.Println(err.Error())
	}

	res := parse(payload)

	token = res.AccessToken

	if _, err := w.Write(payload); err != nil {
		log.Println(err.Error())
	}
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err.Error())
	}

	code := r.Form.Get("code")

	data := url.Values{}
	data.Add("grant_type", "authorization_code")
	data.Add("code", code)
	data.Add("redirect_uri", callbackURL)
	data.Add("client_id", clientID)
	data.Add("client_secret", clientSecret)

	payload, err := call("POST", host+"/oauth/token", data, "")

	if err != nil {
		log.Println(err.Error())
	}

	res := parse(payload)

	token = res.AccessToken

	if _, err := w.Write(payload); err != nil {
		log.Println(err.Error())
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("index.html"))

	err := tmpl.Execute(w, struct {
		ClientID    string
		CallbackURL string
	}{
		ClientID:    clientID,
		CallbackURL: callbackURL,
	})

	if err != nil {
		log.Println(err.Error())
	}
}
