package main

// These functions add, update and list all the database entries

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
  "fmt"
	"time"
	"strconv"
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
	db.Exec("create table if not exists users (email text, password text, firstname text, lastname text, credit float)")
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

func createItem(i string, v string) error {

	value, _ := strconv.ParseFloat(v, 64)

	var db, _ = sql.Open("sqlite3", "items.sqlite3")
	defer db.Close()
	db.Exec("create table if not exists items (itemname text, value text, quantity integer)")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into items (itemname, value, quantity) values (?, ?, ?)")
	_, err := stmt.Exec(i, value, 0)
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

func buyItem(n string, i string){

	cost := changeItemQuantity(i, -1)

	var tax float64 = 0.0

	charge := -cost - tax

	adjustUserCredit(n, charge)
}

func sellItems(n string, i string, q string){

	quantity, _ := strconv.Atoi(q)

	unitprice := changeItemQuantity(i, quantity)

	price := unitprice*float64(quantity)

	adjustUserCredit(n, price)
}

func changeItemQuantity(i string, q int) float64{
	var db, _ = sql.Open("sqlite3", "items.sqlite3")
	defer db.Close()

	x, _ := db.Query("select itemname, value, quantity from items")

	var it string
	var val float64
	var qty int

	var charge float64
	var newqty int

	for x.Next(){
		x.Scan(&it, &val, &qty)
		if it == i{
			charge = val
			newqty = qty+q
		}
	}
	r, _ := db.Prepare("update items set quantity = '" + strconv.Itoa(newqty) + "' where itemname = '" + i + "'")
	r.Exec()
	fmt.Sprintf("Item %s changed by %d to %d", i, q, newqty)
	return charge
}

func itemDetails(i string) (float64, int){
	var db, _ = sql.Open("sqlite3", "items.sqlite3")
	defer db.Close()

	x, _ := db.Query("select itemname, value, quantity from items")

	var it string
	var val float64
	var qty int

	var value float64
	var currentQty int

	for x.Next(){
		x.Scan(&it, &val, &qty)
		if it == i{
			value = val
			currentQty = qty
		}
	}
	return value, currentQty
}

func stocktake(name string, i string, q string){
	newquantity, _ := strconv.Atoi(q)
	_, oldquantity := itemDetails(i)
	quantityadjustment := newquantity - oldquantity
	changeItemQuantity(i, quantityadjustment)
	fmt.Printf("Name: %s\tItem: %s\tNew Qty: %d\n", name, i, newquantity)
}

func adjustUserCredit(name string, charge float64){

	var db, _ = sql.Open("sqlite3", "users.sqlite3")
	defer db.Close()

	credit := userCredit(name)

	newcredit := credit + charge

	r, _ := db.Prepare("update users set credit = ? where firstname = '" + name + "'")
	r.Exec(newcredit)

	fmt.Printf("User %s adjusted by $%6.2f and has a final credit of $%6.2f\n", name, charge, newcredit)
}

func userCredit(name string)float64{
	var db, _ = sql.Open("sqlite3", "users.sqlite3")
	defer db.Close()

	var fn string
	var cr float64
	var credit float64

	y, _ := db.Query("select firstname, credit from users")
	for y.Next(){
		y.Scan(&fn, &cr)
		if fn == name{
			credit = cr
		}
	}
	return credit
}

/*func main(){  //for debugging I suppose
  fmt.Println("Welcome to the Healthy Food Store!")
  fmt.Println("----------------------------------\n")

  var assval, _ = assetValue()
  fmt.Print("Total asset value: $")
  fmt.Printf("%.2f\n",assval)
  fmt.Println("")

  var totcash, _ = totalUserCredit()
  fmt.Print("Total cash value: $")
  fmt.Printf("%.2f\n",totcash)

	listUsers()

}*/
