package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/html"
)

type Link struct {
	Href, Text string
}

const htmlFile = "ex1.html"

func getLinkAttrs(aNode *html.Node) Link {
	for _, attr := range aNode.Attr {
		if attr.Key == "href" {
			return Link{Href: attr.Val, Text: aNode.FirstChild.Data}
		}
	}
	return Link{}
}

func scanChildren(someNode *html.Node) []Link {

	links := make([]Link, 0)

	for node := someNode.FirstChild; node != nil; node = node.NextSibling {
		if node.Type == html.ElementNode && node.Data == "a" {
			links = append(links, getLinkAttrs(node))
		}

		if node.FirstChild != nil {
			links = append(links, scanChildren(node)...)
		}
	}

	return links
}

func main() {

	file, err := os.OpenFile(htmlFile, os.O_RDONLY, 0666)

	if err != nil {
		log.Fatalf("Failed to open file %s", htmlFile)
	}

	defer file.Close()

	rootNode, err := html.Parse(file)

	if err != nil {
		log.Fatalf("Failed to load html, reason %s ", err)
	}

	fmt.Println("Links ", scanChildren(rootNode))
}
