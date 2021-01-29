package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, wtf is %s ?!", r.URL.Path[1:])
}

func main() {

	p := os.Getenv("PORT")

	if p == "" {
		p = "8080"
	}

	port := fmt.Sprintf(":%s", p)

	http.HandleFunc("/", rootHandler)

	log.Println("Listening on port: ", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
