package main

import (
	"encoding/json"
	"strconv"
	// "fmt"
	"io/ioutil"
	"log"
	"net/http"

	// "github.com/okanemo/Andi-Sultan-Asharil-Raphi/database"

	// "github.com/okanemo/Andi-Sultan-Asharil-Raphi/entity"
	"database/sql"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	// "github.com/okanemo/Andi-Sultan-Asharil-Raphi/controllers"
)

type Account struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Balance int `json:"balance"`
}

var db *sql.DB
var err error

func main() {
	initDB()

	router := mux.NewRouter()
	initaliseHandlers(router)
	log.Fatal(http.ListenAndServe(":8081", router))
}

func initaliseHandlers(router *mux.Router) {
	router.HandleFunc("/register", Register).Methods("POST")
	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/add-balance", AddBalance).Methods("POST")
	// router.HandleFunc("/get/{id}", controllers.GetPersonByID).Methods("GET")
	// router.HandleFunc("/update/{id}", controllers.UpdatePersonByID).Methods("PUT")
	// router.HandleFunc("/delete/{id}", controllers.DeletPersonByID).Methods("DELETE")
}

func initDB() {
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/jr_golang_test")
	if err != nil {
		panic(err.Error())
	}
	// defer db.Close()
}

func Register(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare("INSERT INTO accounts(name,balance) VALUES(?,0)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	name := keyVal["name"]
	_, err = stmt.Exec(name)
	if err != nil {
		panic(err.Error())
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	name := keyVal["name"]
	result, err := db.Query("SELECT * FROM accounts WHERE name = ?", name)
	if err != nil {
		panic(err.Error())
	}
	var account Account
	for result.Next() {
		err := result.Scan(&account.ID, &account.Name, &account.Balance)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(account)
}

func AddBalance(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	name := keyVal["name"]
	balance := keyVal["balance"]

	result, err := db.Query("SELECT * FROM accounts WHERE name = ?", name)
	if err != nil {
		panic(err.Error())
	}

	var account Account
	for result.Next() {
		err := result.Scan(&account.ID, &account.Name, &account.Balance)
		if err != nil {
			panic(err.Error())
		}
	}
	intBalance, err := strconv.Atoi(balance)

	var newBalance = intBalance + account.Balance
	account.Balance = newBalance

	log.Println(account)

	stmt, err := db.Prepare("UPDATE accounts SET balance = ? WHERE name = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(newBalance, name)
	if err != nil {
		panic(err.Error())
	}

	json.NewEncoder(w).Encode(account)
}
