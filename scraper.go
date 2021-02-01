package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var baseSite string = "https://knoxon.co"

func isNewPage(link string) bool {

	if strings.Contains(link, "#") {
		return false
	}

	if strings.Contains(link, "https") {
		return false
	}

	return true
}

func parseDoc(node *html.Node, w http.ResponseWriter) {

	if node.Type == html.ElementNode && node.Data == "title" {
		fmt.Fprintf(w, "SITE TITLE: %s\n", node.FirstChild.Data)
	}

	if node.Type == html.ElementNode && node.Data == "a" {
		for _, a := range node.Attr {
			if a.Key == "href" && isNewPage(a.Val) {
				fmt.Fprintf(w, "%s\n", a.Val)
				// newSite := baseSite + a.Val
				break
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		parseDoc(c, w)
	}
}

func scrape(site string, w http.ResponseWriter, r *http.Request) {

	response, errorGet := http.Get(site)
	if errorGet != nil {
		log.Fatal(errorGet)
	}

	bytes, errorRead := ioutil.ReadAll(response.Body)
	if errorRead != nil {
		log.Fatal(errorRead)
	}

	stringBody := string(bytes)
	response.Body.Close()

	doc, errorParse := html.Parse(strings.NewReader(stringBody))
	if errorParse != nil {
		log.Fatal(errorParse)
	}

	parseDoc(doc, w)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	scrape(baseSite, w, r)
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
