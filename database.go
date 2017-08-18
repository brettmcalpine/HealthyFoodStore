package main

// These functions add, update and list all the database entries

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
  "fmt"
	"time"
)

type User struct {
  Email    string
  Password string
	Fname    string
	Lname    string
  Credit   float64
}

type Item struct{
  Itemname string
  Value float64
  Quantity int
}

type Transaction struct{
  Date time.Time
  User
  Item
	Tax float64
}

func userReal(u *User) bool {
	var db, _ = sql.Open("sqlite3", "users.sqlite3")
	defer db.Close()
	db.Exec("create table if not exists users (email text, password text, firstname text, lastname text, credit real)")
	var em, pw string
	q, err := db.Query("select email, password from users where email = '" + u.Email +"'")
	if err != nil {
		return false
	}
	for q.Next(){
		q.Scan(&em, &pw)
	}
	if em == u.Email && pw == u.Password{
		return true
	}
  return false
}

func userExists(u *User) bool {
	var db, _ = sql.Open("sqlite3", "users.sqlite3")
	defer db.Close()
	db.Exec("create table if not exists users (email text, password text, firstname text, lastname text, credit real)")
	var em string
	q, _ := db.Query("select email from users where email = '" + u.Email +"'")
	for q.Next(){
		q.Scan(&em)
	}
	if em == u.Email{
		return true
	}
  return false
}

func createUser(u *User) error {
	var db, _ = sql.Open("sqlite3", "users.sqlite3")
	defer db.Close()
	db.Exec("create table if not exists users (email text, password text, firstname text, lastname text, credit real)")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into users (email, password, firstname, lastname, credit) values (?, ?, ?, ?, ?)")
	_, err := stmt.Exec(u.Email, u.Password, u.Fname, u.Lname, u.Credit)
	tx.Commit()
	return err
}

func listUsers() error{
  var db, _ = sql.Open("sqlite3", "users.sqlite3")
  defer db.Close()
  var em, fn string
  var cash float64
  q, err := db.Query("select email, firstname, credit from users")
  for q.Next(){
    q.Scan(&em, &fn, &cash)
    fmt.Print("Email: ")
    fmt.Print(em)
    fmt.Print("\tName: ")
    fmt.Print(fn)
    fmt.Print("\tCredit: $")
    fmt.Printf("%.2f\n",cash)
  }
  return err
}

func totalUserCredit() (float64, error){
  var db, _ = sql.Open("sqlite3", "users.sqlite3")
  defer db.Close()
  var cash, totcash float64
  q, err := db.Query("select credit from users")
  for q.Next(){
    q.Scan(&cash)
    totcash = totcash + cash
  }
  return totcash, err
}

func createItem(i *Item) error {
	var db, _ = sql.Open("sqlite3", "items.sqlite3")
	defer db.Close()
	db.Exec("create table if not exists items (itemname text, value text, quantity integer)")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into items (itemname, value, quantity) values (?, ?, ?)")
	_, err := stmt.Exec(i.Itemname, i.Value, i.Quantity)
	tx.Commit()
	return err
}

func itemEmpty(i *Item) bool {
	var db, _ = sql.Open("sqlite3", "items.sqlite3")
	defer db.Close()
	var it string
	q, err := db.Query("select itemname from items where itemname = '" + i.Itemname +"'")
	if err != nil {
		return true
	}
	for q.Next(){
		q.Scan(&it)
	}
	if it == i.Itemname {
		return false
	}
  return true
}

func listItems() []Item{
  var db, _ = sql.Open("sqlite3", "items.sqlite3")
  defer db.Close()
  var it string
  var val float64
  var qty int
	var list_of_items []Item
  q, _ := db.Query("select itemname, value, quantity from items")
  for q.Next(){
    q.Scan(&it, &val, &qty)
		var thing Item
		thing.Itemname = it
	  thing.Value = val
	  thing.Quantity = qty
		list_of_items = append(list_of_items, thing)
  }
  return list_of_items
}

func assetValue() (float64, error){
  var db, _ = sql.Open("sqlite3", "items.sqlite3")
  defer db.Close()
  var val float64
  var qty int
  var totval float64
  q, err := db.Query("select value, quantity from items")
  for q.Next(){
    q.Scan(&val, &qty)
    totval = totval + (val*float64(qty))
  }
  return totval, err
}

func buyItem(u *User, i *Item, qty int){
  var users, _ = sql.Open("sqlite3", "users.sqlite3")
	defer users.Close()
  var items, _ = sql.Open("sqlite3", "items.sqlite3")
	defer items.Close()
}

/*func main(){  //for debugging I suppose
  fmt.Println("Welcome to the Healthy Food Store!")
  fmt.Println("----------------------------------\n")

	j := Item{"Coke", 1.00, 30}
  if itemEmpty(&j){
    createItem(&j)
  }

  var assval, _ = assetValue()
  fmt.Print("Total asset value: $")
  fmt.Printf("%.2f\n",assval)
  fmt.Println("")

  var totcash, _ = totalUserCredit()
  fmt.Print("Total cash value: $")
  fmt.Printf("%.2f\n",totcash)

}*/
