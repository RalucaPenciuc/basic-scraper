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

	p, ok := os.LookupEnv("PORT")

	if !ok {
		p = "8080"
	}

	port := fmt.Sprintf(":%s", p)

	http.HandleFunc("/", rootHandler)

	log.Fatal(http.ListenAndServe(port, nil))
}
