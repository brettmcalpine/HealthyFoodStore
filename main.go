package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

var router = mux.NewRouter()

func indexPage(w http.ResponseWriter, r *http.Request) {
	msg, _ := getMsg(w, r, "message")
	if msg != nil {
		tmpl, _ := template.ParseFiles("base.html", "index.html", "main.html", "flash.html")
		err := tmpl.ExecuteTemplate(w, "base", msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {

		u := &User{}
		tmpl, _ := template.ParseFiles("base.html", "index.html", "main.html")
		err := tmpl.ExecuteTemplate(w, "base", u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	pass := r.FormValue("password")
	u := &User{Email: email, Password: pass}
	redirect := "/"
	if email != "" && pass != "" {
		if userReal(u) == true {
			setSession(u, w)
			redirect = "/buysell"
		} else {
			setMsg(w, "message", []byte("Please signup or enter a valid email and password!"))
		}
	} else {
		setMsg(w, "message", []byte("Email or Password field are empty!"))
	}
	http.Redirect(w, r, redirect, 302)
}

func logout(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/", 302)
}

func buysell(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("base.html", "index.html")
	firstname := getUserName(r)
	if firstname != "" {
		err := tmpl.ExecuteTemplate(w, "base", &User{Fname: firstname})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		setMsg(w, "message", []byte("Please login first!"))
		http.Redirect(w, r, "/", 302)
	}
}

func buy(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("buy.html", "index.html", "internal.html")
	firstname := getUserName(r)
	if firstname != "" {
		err := tmpl.ExecuteTemplate(w, "buy", &User{Fname: firstname})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		setMsg(w, "message", []byte("Please login first!"))
		http.Redirect(w, r, "/", 302)
	}
}

func sell(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("base.html", "index.html", "internal.html")
	firstname := getUserName(r)
	if firstname != "" {
		err := tmpl.ExecuteTemplate(w, "base", &User{Fname: firstname})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		setMsg(w, "message", []byte("Please login first!"))
		http.Redirect(w, r, "/", 302)
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl, _ := template.ParseFiles("signup.html", "index.html", "base.html")
		u := &User{}
		tmpl.ExecuteTemplate(w, "base", u)
	case "POST":
		f := r.FormValue("fName")
		l := r.FormValue("lName")
		em := r.FormValue("email")
		pass := r.FormValue("password")

		u := &User{Fname: f, Lname: l, Email: em, Password: pass}
		createUser(u)
		http.Redirect(w, r, "/", 302)
	}
}

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	router.HandleFunc("/", indexPage)
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/logout", logout).Methods("POST")
	router.HandleFunc("/buysell", buysell)
	router.HandleFunc("/buy", buy)
	router.HandleFunc("/sell", sell)
	router.HandleFunc("/signup", signup).Methods("POST", "GET")
	http.Handle("/", router)
	http.ListenAndServe(":5050", nil)
}
