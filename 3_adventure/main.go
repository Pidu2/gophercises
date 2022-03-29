// use html/template
// use http.Handler
// use encoding/json
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type Arc struct {
	Title   string
	Story   []string
	Options []struct {
		Text string
		Arc  string
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getStoryArcsFromJSON(story_file string) map[string]Arc {
	dat, err := os.Open(story_file)
	check(err)
	defer dat.Close()
	byteValue, _ := ioutil.ReadAll(dat)

	var objmap map[string]json.RawMessage
	json.Unmarshal([]byte(byteValue), &objmap)

	arcs := make(map[string]Arc)

	for key, _ := range objmap {
		var arc Arc
		json.Unmarshal(objmap[key], &arc)
		arcs[key] = arc
	}
	return arcs
}

func handler(arcs map[string]Arc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("layout.html"))
		arc := arcs[strings.TrimLeft(r.URL.Path, "/")]
		tmpl.Execute(w, arc)
	}
}

func main() {
	story_file := os.Args[1]
	arcs := getStoryArcsFromJSON(story_file)
	handler := handler(arcs)
	http.ListenAndServe(":8080", handler)
}
