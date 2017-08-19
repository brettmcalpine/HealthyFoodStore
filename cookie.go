package main

import ("net/http"
				"database/sql"
				_ "github.com/mattn/go-sqlite3")

func setSession(u *User, w http.ResponseWriter) {
	value := map[string]string{
		"email": u.Email,
		"pass": u.Password,
		"firstname": u.Fname,
	}
	if encoded, err := cookieHandler.Encode("healthyfoodstore", value); err == nil {
		cookie := &http.Cookie{
			Name:  "healthyfoodstore",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func getUserName(r *http.Request) (firstname string) {
	if cookie, err := r.Cookie("healthyfoodstore"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("healthyfoodstore", cookie.Value, &cookieValue); err == nil {
			email := cookieValue["email"]
			var db, _ = sql.Open("sqlite3", "users.sqlite3")
		  defer db.Close()
			var em, fn string
			q, _ := db.Query("select email, firstname from users")
			for q.Next(){
				q.Scan(&em, &fn)
				if em == email{
					firstname = fn
				}
			}
		}
	}
	return firstname
}

func getUserDetails(r *http.Request) (u User) {
	if cookie, err := r.Cookie("healthyfoodstore"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("healthyfoodstore", cookie.Value, &cookieValue); err == nil {
			email := cookieValue["email"]
			var db, _ = sql.Open("sqlite3", "users.sqlite3")
		  defer db.Close()
			var em, fn string
			var cr float64
			q, _ := db.Query("select email, firstname, credit from users")
			for q.Next(){
				q.Scan(&em, &fn, &cr)
				if em == email{
					u = User{Fname: fn, Email: em, Credit: cr}
				}
			}
		}
	}
	return u
}

func clearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "healthyfoodstore",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
