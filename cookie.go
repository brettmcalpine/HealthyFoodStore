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
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func getUserName(r *http.Request) (firstname string) {
	if cookie, err := r.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
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

func clearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
