package main

import (
	"html/template"
	"net/http"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

var router = mux.NewRouter()

func login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl, _ := template.ParseFiles("login.html", "footer.html")
		u := &User{}
		tmpl.ExecuteTemplate(w, "login", u)
	case "POST":
		email := r.FormValue("email")
		pass := r.FormValue("password")
		u := &User{Email: email, Password: pass}
		redirect := "/login"
		if email != "" && pass != "" {
			if userReal(u) == true {
				setSession(u, w)
				redirect = "/buy"
			} else {
				setMsg(w, "message", []byte("Please signup or enter a valid email and password!"))
			}
		} else {
			setMsg(w, "message", []byte("Email or Password field are empty!"))
		}
		http.Redirect(w, r, redirect, 302)

		}
}

func signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl, _ := template.ParseFiles("signup.html", "footer.html")
		u := &User{}
		tmpl.ExecuteTemplate(w, "signup", u)
	case "POST":
		f := r.FormValue("fName")
		l := r.FormValue("lName")
		em := r.FormValue("email")
		pass := r.FormValue("password")
		u := &User{Fname: f, Lname: l, Email: em, Password: pass}
		if !userExists(u){
			fmt.Println("Creating new user")
			createUser(u)
			http.Redirect(w, r, "/buy", 302)
		} else {
			fmt.Println("User already exists")
			setMsg(w, "message", []byte("Email already in use!"))
			http.Redirect(w, r, "/", 302)
		}
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/", 302)
}

func buy(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
	tmpl, _ := template.ParseFiles("buy.html", "shopheader.html", "footer.html")
	userdata := getUserDetails(r)
	firstname := getUserName(r)
	items := listItems()
	data := struct{
		U User
		I []Item
	}{userdata, items}
	if firstname != "" {
		err := tmpl.ExecuteTemplate(w, "buy", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		setMsg(w, "message", []byte("Please login first!"))
		http.Redirect(w, r, "/login", 302)
	}
	case "POST":
		i := r.FormValue("item")
		firstname := getUserName(r)
		buyItem(firstname, i)
		http.Redirect(w, r, "/buy", 302)
	}
}

func sell(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
	tmpl, _ := template.ParseFiles("sell.html", "shopheader.html", "footer.html")
	userdata := getUserDetails(r)
	firstname := getUserName(r)
	items := listItems()
	data := struct{
		U User
		I []Item
	}{userdata, items}
	if firstname != "" {
		err := tmpl.ExecuteTemplate(w, "sell", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		setMsg(w, "message", []byte("Please login first!"))
		http.Redirect(w, r, "/", 302)
	}
	case "POST":
		i := r.FormValue("sell-item")
		q := r.FormValue("sell-quantity")
		firstname := getUserName(r)
		sellItems(firstname, i, q)
		http.Redirect(w, r, "/sell", 302)
	}
}

func shopkeeping(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("shopkeeping.html", "shopheader.html", "footer.html")
	userdata := getUserDetails(r)
	firstname := getUserName(r)
	items := listItems()
	data := struct{
		U User
		I []Item
	}{userdata, items}
	if firstname != "" {
		err := tmpl.ExecuteTemplate(w, "shopkeeping", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		setMsg(w, "message", []byte("Please login first!"))
		http.Redirect(w, r, "/", 302)
	}
}

func stocktakePage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
	tmpl, _ := template.ParseFiles("stocktake.html", "shopheader.html", "footer.html")
	userdata := getUserDetails(r)
	firstname := getUserName(r)
	items := listItems()
	data := struct{
		U User
		I []Item
	}{userdata, items}
	if firstname != "" {
		err := tmpl.ExecuteTemplate(w, "stocktake", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		setMsg(w, "message", []byte("Please login first!"))
		http.Redirect(w, r, "/", 302)
	}
	case "POST":
		i := r.FormValue("stocktake-item")
		q := r.FormValue("stocktake-quantity")
		firstname := getUserName(r)
		stocktake(firstname, i, q)
		http.Redirect(w, r, "/stocktake", 302)
	}
}

func createPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
	tmpl, _ := template.ParseFiles("newitem.html", "shopheader.html", "footer.html")
	userdata := getUserDetails(r)
	firstname := getUserName(r)
	items := listItems()
	data := struct{
		U User
		I []Item
	}{userdata, items}
	if firstname != "" {
		err := tmpl.ExecuteTemplate(w, "newitem", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		setMsg(w, "message", []byte("Please login first!"))
		http.Redirect(w, r, "/", 302)
	}
	case "POST":
		i := r.FormValue("create-name")
		v := r.FormValue("create-value")
		//firstname := getUserName(r)
		createItem(i, v)
		http.Redirect(w, r, "/buy", 302)
	}
}

func deleteItemPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
	tmpl, _ := template.ParseFiles("deleteitem.html", "shopheader.html", "footer.html")
	userdata := getUserDetails(r)
	firstname := getUserName(r)
	items := listItems()
	data := struct{
		U User
		I []Item
	}{userdata, items}
	if firstname != "" {
		err := tmpl.ExecuteTemplate(w, "deleteitem", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		setMsg(w, "message", []byte("Please login first!"))
		http.Redirect(w, r, "/", 302)
	}
	case "POST":
		i := r.FormValue("delete-name")
		//firstname := getUserName(r)
		deleteItem(i)
		http.Redirect(w, r, "/buy", 302)
	}
}

func deleteUserPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
	tmpl, _ := template.ParseFiles("deleteuser.html", "shopheader.html", "footer.html")
	userdata := getUserDetails(r)
	firstname := getUserName(r)
	items := listItems()
	data := struct{
		U User
		I []Item
	}{userdata, items}
	if firstname != "" {
		err := tmpl.ExecuteTemplate(w, "deleteuser", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		setMsg(w, "message", []byte("Please login first!"))
		http.Redirect(w, r, "/", 302)
	}
	case "POST":
		check := r.FormValue("delete-user")
		firstname := getUserName(r)
		deleteUser(firstname, check)
		http.Redirect(w, r, "/signup", 302)
	}
}

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	router.HandleFunc("/", login)
	router.HandleFunc("/login", login).Methods("POST", "GET")
	router.HandleFunc("/logout", logout).Methods("POST")
	router.HandleFunc("/buy", buy).Methods("POST", "GET")
	router.HandleFunc("/shopkeeping", shopkeeping).Methods("POST", "GET")
	router.HandleFunc("/stocktake", stocktakePage).Methods("POST", "GET")
	router.HandleFunc("/newitem", createPage).Methods("POST", "GET")
	router.HandleFunc("/deleteitem", deleteItemPage).Methods("POST", "GET")
	router.HandleFunc("/deleteuser", deleteUserPage).Methods("POST", "GET")
	router.HandleFunc("/sell", sell).Methods("POST", "GET")
	router.HandleFunc("/signup", signup).Methods("POST", "GET")
	http.Handle("/", router)
	http.ListenAndServe(":5050", nil)
}
