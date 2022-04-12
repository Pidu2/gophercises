package main

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/net/html"
)

type Link struct {
	Text string
	Href string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	file_to_parse, err := os.Open(os.Args[1])
	check(err)
	parseHtmlFile(file_to_parse)
}

func parseHtmlFile(content io.Reader) {
	doc, err := html.Parse(content)
	check(err)
	var linklist []Link
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		parseNode(*c, &linklist)
	}
	fmt.Println(linklist)
}

func parseNode(node html.Node, linklist *[]Link) {
	for _, a := range node.Attr {
		if a.Key == "href" {
			text := collectText(&node)
			*linklist = append(*linklist, Link{text, a.Val})
			return
		}
	}
	if node.FirstChild == nil {
		return
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		parseNode(*c, linklist)
	}
}

func collectText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	} else {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			return collectText(c)
		}
	}
	return ""
}
