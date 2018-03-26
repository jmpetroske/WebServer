package main

import (
	"net/http"
	"strings"
	"fmt"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message

	w.Write([] byte(message))
}

func main() {
	http.HandleFunc("/", sayHello)
	panic(http.ListenAndServe(":8080", nil))
	fmt.Print("test")
}
