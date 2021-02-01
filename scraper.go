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

// MakeSet initialize the set
func makeSet() *Set {
	return &Set{
		container: make(map[string]struct{}),
	}
}

var links = makeSet()

// Set implementation
type Set struct {
	container map[string]struct{}
}

// Size of the set
func (s *Set) Size() int {
	return len(s.container)
}

// Exists check
func (s *Set) Exists(key string) bool {
	_, exists := s.container[key]
	return exists
}

// Add element to set
func (s *Set) Add(key string) {
	s.container[key] = struct{}{}
}

// Remove element from set
func (s *Set) Remove(key string) error {
	_, exists := s.container[key]
	if !exists {
		return fmt.Errorf("Remove Error: Item does not exist")
	}
	delete(s.container, key)
	return nil
}

func isNewPage(link string) bool {

	if strings.Contains(link, "#") {
		return false
	}

	if strings.Contains(link, "https") {
		return false
	}

	if strings.Contains(link, "@") {
		return false
	}

	if strings.Contains(link, "tel") {
		return false
	}

	return true
}

func parseDoc(node *html.Node, w http.ResponseWriter, r *http.Request) {

	if node.Type == html.ElementNode && node.Data == "title" {
		fmt.Fprintf(w, "SITE TITLE: %s\n", node.FirstChild.Data)
	}

	if node.Type == html.ElementNode && node.Data == "a" {
		for _, a := range node.Attr {
			if a.Key == "href" && isNewPage(a.Val) {
				fmt.Fprintf(w, "%s\n", a.Val)
				newSite := baseSite + a.Val
				// fmt.Println(newSite)
				if !links.Exists(newSite) {
					links.Add(newSite)
					scrape(newSite, w, r)
				}
				break
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		parseDoc(c, w, r)
	}
}

func scrape(site string, w http.ResponseWriter, r *http.Request) {

	links.Add(site)

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

	parseDoc(doc, w, r)
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
