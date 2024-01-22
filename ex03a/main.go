package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

type StoryArc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

func loadStories(name string) (storyArcs map[string]StoryArc, err error) {
	jsonFile, err := os.OpenFile(name, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Failed to load file ", err)
		return
	}

	err = json.NewDecoder(jsonFile).Decode(&storyArcs)
	if err != nil {
		fmt.Println("Failed to load JSON file ", err)
		return
	}

	defer jsonFile.Close()

	return
}

func loadWebage(storyArc StoryArc, w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("templates/Story.html")

	if err != nil {
		log.Fatalln(err)
	}

	templ.Execute(w, storyArc)
}

func main() {
	const storyFile = "gopher.json"
	storyArcs, err := loadStories(storyFile)
	if err != nil {
		log.Fatalln("Failed", err)
	}
	fmt.Println("Loaded all entries ", len(storyArcs))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		reqPath := r.URL.Path
		if reqPath == "/favicon.ico" {
			w.WriteHeader(200)
			return
		}

		pathRegEx := regexp.MustCompile("/([A-Za-z0-9- ]*)$")
		arcNames := pathRegEx.FindStringSubmatch(reqPath)

		if len(arcNames) < 2 {
			w.WriteHeader(400)
			return
		}

		arcName := arcNames[1]
		if len(arcName) == 0 {
			arcName = "intro"
		}

		storyArc, ok := storyArcs[arcName]
		if !ok {
			w.WriteHeader(404)
			return
		}
		loadWebage(storyArc, w, r)
	})
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
