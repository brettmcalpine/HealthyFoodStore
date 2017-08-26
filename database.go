package main

// These functions add, update and list all the database entries

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
	"strconv"
	"fmt"
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
  Timestamp time.Time
	TType string //Transaction Type (Buy, Sell, Transfer, Stocktake, ItemAdd, ItemDelete, UserCreate, UserDelete, UserSetCredit)
  UID	string
	Dollars float64
  Unit string
	Count int
	Tax float64 //For this transaction
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

	time := time.Now()
	log := Transaction{Timestamp:time, TType:"usercreate", UID:u.Fname, Dollars:0, Unit:"", Count:0, Tax:0.0}
	LogTransaction(log)

	return err
}

func listUsers() ([]User, error){
	var db, _ = sql.Open("sqlite3", "users.sqlite3")
  defer db.Close()
  var em, fn, ln string
	var cr float64
	var list_of_users []User
  q, err := db.Query("select email, firstname, lastname, credit from users")
	if err != nil{
		return list_of_users, err
	}
  for q.Next(){
    q.Scan(&em, &fn, &ln, &cr)
		var person User
		person.Email = em
		person.Fname = fn
	  person.Lname = ln
	  person.Credit = cr
		list_of_users = append(list_of_users, person)
  }
  return list_of_users, err
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

func tax()float64{
	cash, _ := totalUserCredit()
	assets, _ := assetValue()
	tax := cash/assets - 1
	if tax <= 0 || assets == 0{
		return 0
	}
	return tax
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

	time := time.Now()
	log := Transaction{Timestamp:time, TType:"itemcreate", UID:"", Dollars:value, Unit:i, Count:0, Tax:0.0}
	LogTransaction(log)

	return err
}

func deleteItem(i string) error {
	var db, _ = sql.Open("sqlite3", "items.sqlite3")
	defer db.Close()
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("delete from items where (itemname, quantity) = (?, ?)")
	_, err := stmt.Exec(i, 0)
	tx.Commit()

	time := time.Now()
	log := Transaction{Timestamp:time, TType:"itemdelete", UID:"", Dollars:0, Unit:i, Count:0, Tax:0.0}
	LogTransaction(log)

	return err
}

func deleteUser(firstname string, check string){
	if check == "on"{
		if userCredit(firstname)>0{
			var db, _ = sql.Open("sqlite3", "users.sqlite3")
			defer db.Close()
			tx, _ := db.Begin()
			stmt, _ := tx.Prepare("delete from users where firstname = ?")
			stmt.Exec(firstname)

			time := time.Now()
			log := Transaction{Timestamp:time, TType:"userdelete", UID:firstname, Dollars:0, Unit:"", Count:0, Tax:0.0}
			LogTransaction(log)

			tx.Commit()
		}
	}
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

func listItems() ([]Item, error){
  var db, _ = sql.Open("sqlite3", "items.sqlite3")
  defer db.Close()
  var it string
  var val float64
  var qty int
	var list_of_items []Item
  q, err := db.Query("select itemname, value, quantity from items")
	if err != nil{
		return list_of_items, err
	}
  for q.Next(){
    q.Scan(&it, &val, &qty)
		var thing Item
		thing.Itemname = it
	  thing.Value = val
	  thing.Quantity = qty
		list_of_items = append(list_of_items, thing)
  }
  return list_of_items, err
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

	tax := tax()
	cost := changeItemQuantity(i, -1)
	charge := cost * (1 + tax)
	adjustUserCredit(n, -charge)

	time := time.Now()
	logtax := cost*tax
	log := Transaction{Timestamp:time, TType:"buy", UID:n, Dollars:cost, Unit:i, Count:-1, Tax:logtax}
	LogTransaction(log)
}

func sellItems(n string, i string, q string){

	quantity, _ := strconv.Atoi(q)
	unitprice := changeItemQuantity(i, quantity)
	price := unitprice*float64(quantity)
	adjustUserCredit(n, price)

	time := time.Now()
	log := Transaction{Timestamp:time, TType:"sell", UID:n, Dollars:price, Unit:i, Count:quantity, Tax:0.0}
	LogTransaction(log)
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

	time := time.Now()
	log := Transaction{Timestamp:time, TType:"stocktake", UID:name, Dollars:0, Unit:i, Count:quantityadjustment, Tax:0.0}
	LogTransaction(log)

}

func adjustUserCredit(name string, charge float64){
	var db, _ = sql.Open("sqlite3", "users.sqlite3")
	defer db.Close()

	credit := userCredit(name)
	newcredit := credit + charge

	r, _ := db.Prepare("update users set credit = ? where firstname = '" + name + "'")
	r.Exec(newcredit)
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

func transferFunds(from string, to string, dollars string){
	transfer, _ := strconv.ParseFloat(dollars, 64)
	adjustUserCredit(from, -transfer)
	adjustUserCredit(to, transfer)

}

func LogTransaction(t Transaction) error{
	var db, _ = sql.Open("sqlite3", "transactions.sqlite3")
	defer db.Close()
	db.Exec("create table if not exists transactions (datetime timestamp, ttype text, uid text, dollars real, item text, count int, tax real)")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into transactions (datetime, ttype, uid, dollars, item, count, tax) values (?, ?, ?, ?, ?, ?, ?)")
	_, err := stmt.Exec(t.Timestamp, t.TType, t.UID, t.Dollars, t.Unit, t.Count, t.Tax)
	tx.Commit()
	return err
}

func undoTransaction(date string){

	var db, _ = sql.Open("sqlite3", "transactions.sqlite3")

	var fn, it, tt string
	var cnt int
	var dol, tax float64
	var tm time.Time

	var t Transaction

	dumbtime, _ := time.Parse("2006-01-02 15:04:05.999999 -0700 -0700", date)

	y, _ := db.Query("select datetime, ttype, uid, dollars, item, count, tax from transactions") // where datetime = " + t.Timestamp
	for y.Next(){
		y.Scan(&tm, &tt, &fn, &dol, &it, &cnt, &tax)
			t.Timestamp = tm
			t.TType = tt
			t.UID = fn
			t.Dollars = dol
			t.Unit = it
			t.Count = cnt
			t.Tax = tax
			if t.Timestamp.Equal(dumbtime){
				break
			}
		}
	db.Close()

	if t.TType == "sell"{
		changeItemQuantity(t.Unit, -t.Count)
		adjustUserCredit(t.UID, -t.Dollars)
		time := time.Now()
		log := Transaction{Timestamp:time, TType:"undosell", UID:t.UID, Dollars:t.Dollars, Unit:t.Unit, Count:t.Count, Tax:t.Tax}
		fmt.Println("Fuck?")
		LogTransaction(log)
	}
	if t.TType == "buy"{
		changeItemQuantity(t.Unit, -t.Count)
		cost := t.Dollars + t.Tax
		adjustUserCredit(t.UID, cost)
		time := time.Now()
		log := Transaction{Timestamp:time, TType:"undobuy", UID:t.UID, Dollars:cost, Unit:t.Unit, Count:t.Count, Tax:t.Tax}
		LogTransaction(log)
	}
}

func listTransactions() ([]Transaction, error){
  var db, _ = sql.Open("sqlite3", "transactions.sqlite3")
  defer db.Close()
  var tm time.Time
  var tp, id, it string
  var dol, tax float64
	var cnt int
	var list_of_transactions []Transaction
  q, err := db.Query("select datetime, ttype, uid, dollars, item, count, tax from transactions")
	if err != nil{
		return list_of_transactions, err
	}
  for q.Next(){
    q.Scan(&tm, &tp, &id, &dol, &it, &cnt, &tax)
		var check Transaction
		check.Timestamp = tm
	  check.TType = tp
	  check.UID = id
		check.Dollars = dol
		check.Unit = it
		check.Count = cnt
		check.Tax = tax
		list_of_transactions = append(list_of_transactions, check)
  }
  return list_of_transactions, err
}

/*func main(){  //for debugging I suppose
  fmt.Println("Welcome to the Healthy Food Store!")
  fmt.Println("----------------------------------\n")

	listT, _ := listTransactions()
	fmt.Printf("%+v\n", listT)
}*/
