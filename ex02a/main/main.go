package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/saikrir/gophercises/urlshort"
)

type ViewModel struct {
	Keyword string
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	})

	return mux
}

func fallbackHandler() http.HandlerFunc {
	if notFoundHtml, err := os.OpenFile("../NotFound.html", os.O_RDONLY, 0666); err != nil {
		log.Fatalln("Fallback handler could not locate NotFound.html file ")
	} else {
		htmlFileContent, _ := io.ReadAll(notFoundHtml)
		if htmlTmpl, err := template.New("NotFound").Parse((string(htmlFileContent))); err != nil {
			log.Fatalln("Error Parsing HTML Template", err)
		} else {

			return func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("Path ", r.RequestURI[1:])
				out := new(bytes.Buffer)
				err := htmlTmpl.ExecuteTemplate(out, "NotFound", ViewModel{Keyword: r.RequestURI[1:]})
				if err != nil {
					log.Fatalln("Failed to execute Html Template")
				}
				fmt.Fprintln(w, out)
			}
		}
	}

	return nil
}

func main() {

	knownPaths := map[string]string{
		"/amazon": "http://www.amazon.com",
		"/google": "http://www.google.com",
	}

	mapHandler := urlshort.MapHandler(knownPaths, fallbackHandler())

	yaml := `
- path: /go
  url: https://go.dev/
- path: /reddit
  url: https://www.reddit.com/
`

	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)

	if err != nil {
		log.Fatalln("Failed to parse YAML Config ", err)
	}

	http.ListenAndServe(":4200", yamlHandler)
}
