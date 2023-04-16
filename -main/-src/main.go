package main

import (
	"fmt"
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

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/checkout", checkoutHandler)
	http.HandleFunc("/confirm", confirmHandler)

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

	// Route to display the form for adding a new item
	http.HandleFunc("/add-item", addItemFormHandler)

	// Route to handle form submission and add the new item
	http.HandleFunc("/add-item-submit", addItemSubmitHandler)

	// Route to display all items
	http.HandleFunc("/items", itemsHandler)

	// Start the web server and listen for incoming requests
	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func addItemFormHandler(w http.ResponseWriter, r *http.Request) {
	// Display the form for adding a new item
	tpl := template.Must(template.ParseFiles("add-item.html"))
	tpl.Execute(w, nil)
}

func addItemSubmitHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form data and add a new item to the list
	title := r.FormValue("title")
	description := r.FormValue("description")
	price := r.FormValue("price")

	item := Item{
		Title:       title,
		Description: description,
		Price:       price,
	}

	items = append(items, item)

	// Redirect back to the list of items
	http.Redirect(w, r, "/items", http.StatusSeeOther)
}

func checkoutHandler(w http.ResponseWriter, r *http.Request) {
	itemID := r.FormValue("item_id")

	if itemID == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	item, err := getItemByID(itemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/confirm", http.StatusSeeOther)
}

func confirmHandler(w http.ResponseWriter, r *http.Request) {
	var items []Item
	var total float64

	for _, cookie := range r.Cookies() {
		if cookie.Name[:5] == "item_" {
			itemID := cookie.Name[5:]
			item, err := getItemByID(itemID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			items = append(items, *item)
			total += item.Price
		}
	}
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	// Display a list of all items
	tpl := template.Must(template.ParseFiles("items.html"))
	tpl.Execute(w, items)
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

func addItemFormHandler(w http.ResponseWriter, r *http.Request) {
	// Display the form for adding a new item
	tpl := template.Must(template.ParseFiles("add-item.html"))
	tpl.Execute(w, nil)
}

func addItemSubmitHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form data and add a new item to the list
	title := r.FormValue("title")
	description := r.FormValue("description")
	price := r.FormValue("price")

	item := Item{
		Title:       title,
		Description: description,
		Price:       price,
	}

	items = append(items, item)

	// Redirect back to the list of items
	http.Redirect(w, r, "/items", http.StatusSeeOther)
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	// Display a list of all items
	tpl := template.Must(template.ParseFiles("items.html"))
	tpl.Execute(w, items)
}
