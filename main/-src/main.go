package main

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var tpl *template.Template

var users_a = map[string]string{"user1": "password", "user2": "password"}
var store_data = sessions.NewCookieStore([]byte("secret_key"))

func init() {
	tpl = template.Must(template.ParseGlob("./static/*.html"))
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe(":7777", nil)

	r := mux.NewRouter()
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/logout", logoutHandler).Methods("GET")
	r.HandleFunc("/healthcheck", healthcheck).Methods("GET")
	httpServer := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
	}
	log.Fatal(httpServer.ListenAndServe())
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Pass the data as URL form encoded", http.StatusBadRequest)
		return
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	storedPassword, exists := users_a[username]
	if exists {
		session, _ := store_data.Get(r, "session.id")
		if storedPassword == password {
			session.Values["authenticated"] = true
			session.Save(r, w)
		} else {
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		}
		w.Write([]byte("Login!"))
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := store_data.Get(r, "session.id")
	session.Values["authenticated"] = false
	session.Save(r, w)
	w.Write([]byte("Logout"))
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	session, _ := store_data.Get(r, "session.id")
	authenticated := session.Values["authenticated"]
	if authenticated != nil && authenticated != false {
		w.Write([]byte("Welcome!"))
		return
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
}
