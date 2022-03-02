package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
)

type ErrStruct struct {
	Path string
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	templ := template.Must(template.ParseGlob("templates/*.gohtml")) //define gohtml file
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		errData := &ErrStruct{Path: r.URL.Path}
		err := templ.ExecuteTemplate(w, "error.gohtml", errData)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func handleDirectory(directory, route string) {
	fonts := http.FileServer(http.Dir(directory))
	http.Handle(route, http.StripPrefix(route, fonts))
}

func main() {
	templ := template.Must(template.ParseGlob("templates/*.gohtml")) //define gohtml file

	housesData := getHouses()
	for _, data := range housesData {
		println("Villagers : ", data.Name)
		println("House Interior : ", data.NhDetails.HouseInteriorUrl, "\n")
	}
	// Gestion de tous les fichiers css
	fs := http.FileServer(http.Dir("style"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Gestion de tous les fichiers pics
	pics := http.FileServer(http.Dir("pics"))
	http.Handle("/image/", http.StripPrefix("/image/", pics))

	handleDirectory("./fonts", "/fonts/")
	handleDirectory("./scripts", "/scripts/")

	characters := getCharacters()
	sort.Slice(characters, func(indexFirst, indexSecond int) bool {
		return characters[indexFirst].Name.NameEUen < characters[indexSecond].Name.NameEUen
	})
	for _, character := range characters {

		http.HandleFunc(fmt.Sprintf("/%s", strings.ToLower(character.Name.NameEUen)), func(writer http.ResponseWriter, request *http.Request) {
			name := strings.TrimPrefix(request.URL.Path, "/")
			ch := acnh(name, characters)
			fmt.Println(request.URL, "url request")
			err := templ.Execute(writer, ch)
			if err != nil {
				log.Fatal(err)
			}
		})
	}

	http.HandleFunc("/charalist", func(w http.ResponseWriter, r *http.Request) {
		templ.ExecuteTemplate(w, "charalist.gohtml", characters)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := TrimURLPrefix(r.URL.Path)
		if path == "favicon.ico" {
			return
		}
		fmt.Println("index")
		if path == "" {
			templ.ExecuteTemplate(w, "index.gohtml", "")
		} else if !characterExistence(path, characters) {
			errorHandler(w, r, http.StatusNotFound)
			fmt.Println(path, "=> introuvable")
		}
	})

	http.HandleFunc("/character", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("debug1")
		queries := r.URL.Query()
		if queries.Has("name") {
			fmt.Println("debug2")
			name := r.URL.Query().Get("name")
			name = strings.ToLower(name)
			fmt.Println("name =>", name)
			http.Redirect(w, r, fmt.Sprintf("/%s", name), http.StatusSeeOther)
		}
	})

	fmt.Println("Server started on localhost:8010")
	err := http.ListenAndServe(":8010", nil)
	if err != nil {
		log.Fatal(err)
	}
}
