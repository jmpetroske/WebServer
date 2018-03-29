package main

import (
	"fmt"
	"net/http"
	"strings"
	"golang.org/x/crypto/bcrypt"
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

	redirectTarget := "/?badlogin=1"
	if username != "" && password != "" {
		databaseConnString := "postgres://petroske:@localhost/accounts?sslmode=disable"
		
		db, err := sql.Open("postgres", databaseConnString)
		checkErr(err)

		fmt.Printf("Signing in user: %s\n", username)
		rows, err := db.Query(`SELECT password_hash FROM user_account WHERE username=$1`,
			username)
		checkErr(err)
		defer rows.Close()

		var passwordHash string
		if rows.Next() {
			err = rows.Scan(&passwordHash)
			checkErr(err)

			err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))

			if err == nil {
				redirectTarget = "/hello"
			}
		}
	}
	
	http.Redirect(w, r, redirectTarget, 302)
}

func createAccountHandler(w http.ResponseWriter, r *http.Request) {	
	username := r.FormValue("usr")
	password := r.FormValue("pwd")

	redirectTarget := "/"
	if username != "" && password != "" {
		// TODO check if username already exists
		
		databaseConnString := "postgres://petroske:@localhost/accounts?sslmode=disable"
		
		db, err := sql.Open("postgres", databaseConnString)
		checkErr(err)

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		checkErr(err)
		
		_, err = db.Query(
			`INSERT INTO user_account (username, password_hash) VALUES ($1, $2);`,
			username,
			passwordHash)
		checkErr(err)

		redirectTarget = "/hello";
	}
	http.Redirect(w, r, redirectTarget, 302)
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
