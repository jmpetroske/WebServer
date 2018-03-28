package main

import (
	"net/http"
	"strings"
	"fmt"
    // "golang.org/x/crypto/bcrypt"
	"database/sql"
	_ "github.com/lib/pq"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message

	w.Write([] byte(message))
}

func main() {
	connString := "postgres://petroske:@localhost/accounts?sslmode=disable"
	
	db, err := sql.Open("postgres", connString)
	checkErr(err)
	
	_, err = db.Query(
		`INSERT INTO user_account (username, password_hash) VALUES ($1, $2);`,
		"jmpetroske",
		"password")
	checkErr(err)

	http.HandleFunc("/", sayHello)
	panic(http.ListenAndServe(":8080", nil))
	fmt.Print("test")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
