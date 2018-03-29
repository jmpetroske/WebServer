package main

import (
	"net/http"
	"strings"
//    "golang.org/x/crypto/bcrypt"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message

	w.Write([] byte(message))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("usr")
	password := r.FormValue("pwd")

	if username != "" && password != "" {

		
		databaseConnString := "postgres://petroske:@localhost/accounts?sslmode=disable"
		
		db, err := sql.Open("postgres", databaseConnString)
		checkErr(err)
		
		_, err = db.Query(
			`INSERT INTO user_account (username, password_hash) VALUES ($1, $2);`,
			username,
			"password")
		checkErr(err)

		//err := bcrypt.CompareHashAndPassword
	}
}

func createAccountHandler(w http.ResponseWriter, r *http.Request) {	
	username := r.FormValue("usr")
	password := r.FormValue("pwd")

//	redirectTarget := "/"
	if username != "" && password != "" {
		databaseConnString := "postgres://petroske:@localhost/accounts?sslmode=disable"
		
		db, err := sql.Open("postgres", databaseConnString)
		checkErr(err)
		
		_, err = db.Query(
			`INSERT INTO user_account (username, password_hash) VALUES ($1, $2);`,
			username,
			"password")
		checkErr(err)

	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/hello", sayHello)
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/createAccount", createAccountHandler).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	http.Handle("/", router)
	
	panic(http.ListenAndServe(":8080", router))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
