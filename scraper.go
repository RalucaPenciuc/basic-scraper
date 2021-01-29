package main

import (
	"fmt"
	"log"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, wtf is %s ?!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
