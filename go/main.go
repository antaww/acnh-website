package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

type ErrStruct struct {
	Path string
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	templ := template.Must(template.ParseGlob("templates/error.gohtml")) //define gohtml file
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
	templateFunctions := template.FuncMap{
		"isBirthday": func(birthday string) bool {
			date := time.Now().Format("2/1")
			return birthday == date
		},
	}
	templ, _ := template.New("").Funcs(templateFunctions).ParseGlob("templates/*.gohtml") //define gohtml file

	// Gestion de tous les fichiers css
	fs := http.FileServer(http.Dir("style"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Gestion de tous les fichiers pics
	pics := http.FileServer(http.Dir("pics"))
	http.Handle("/image/", http.StripPrefix("/image/", pics))

	handleDirectory("./fonts", "/fonts/")
	handleDirectory("./scripts", "/scripts/")

	characters := getCharacters()
	houses := getHouses()
	housewares := getHouseware()

	sort.Slice(characters, func(indexFirst, indexSecond int) bool {
		return characters[indexFirst].Name.NameEUen < characters[indexSecond].Name.NameEUen
	})

	for _, character := range characters {
		http.HandleFunc(fmt.Sprintf("/%s", strings.ToLower(character.Name.NameEUen)), func(writer http.ResponseWriter, request *http.Request) {
			name := strings.TrimPrefix(request.URL.Path, "/")
			simplifiedVillager := getSimplified(name, characters, houses)
			fmt.Println(request.URL, "url request")
			err := templ.ExecuteTemplate(writer, "character.gohtml", simplifiedVillager)
			if err != nil {
				log.Fatal(err)
			}
		})
	}

	for i, houseware := range housewares {
		housewName := strings.ReplaceAll(houseware.Name, " ", "")
		housewName = strings.ToLower(housewName)
		housewares[i].NameSimplified = housewName
		simplifiedname := getSimplifiedHouseware(housewName, housewares)
		http.HandleFunc(fmt.Sprintf("/%s"+"_houseware", housewName), func(writer http.ResponseWriter, request *http.Request) {
			fmt.Println(request.URL, "url request")
			err := templ.ExecuteTemplate(writer, "houseware.gohtml", simplifiedname)
			if err != nil {
				log.Fatal(err)
			}
		})
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := TrimURLPrefix(r.URL.Path)
		if path == "favicon.ico" {
			return
		}
		if path == "" {
			templ.ExecuteTemplate(w, "index.gohtml", characters)
		} else if path == "test" {
			templ.ExecuteTemplate(w, "test.gohtml", houses)
		} else if !characterExistence(path, characters) {
			errorHandler(w, r, http.StatusNotFound)
		}
	})

	http.HandleFunc("/charalist", func(w http.ResponseWriter, r *http.Request) {
		for i := range characters {
			characters[i].IsEmpty = true
		}
		path := TrimURLPrefix(r.URL.Path)
		switch r.Method {
		case "GET":
			templ.ExecuteTemplate(w, "charalist.gohtml", characters)
		case "POST":
			villager := r.FormValue("name")
			if villager != "" {
				if !characterExistence(villager, characters) {
					path = villager
					http.Redirect(w, r, path, http.StatusSeeOther)
					errorHandler(w, r, http.StatusNotFound)
					fmt.Println(villager, "=> introuvable")
				} else {
					http.Redirect(w, r, fmt.Sprintf("/%s", strings.ToLower(villager)), http.StatusSeeOther)
				}
			}
			species := r.FormValue("species")
			if species == "All" {
				for i := range characters {
					characters[i].IsEmpty = true
				}
				templ.ExecuteTemplate(w, "charalist.gohtml", characters)

			} else if species != "" {
				for i := range characters {
					characters[i].IsEmpty = false
				}
				for i, chara := range characters {
					if chara.Species == species {
						characters[i].SelectedSpeccy = species
					}
					if i == len(characters)-1 {
						templ.ExecuteTemplate(w, "charalist.gohtml", characters)
						for i := range characters {
							characters[i].SelectedSpeccy = ""
						}
					}
					i++
				}
			}
		}
	})

	http.HandleFunc("/character", func(w http.ResponseWriter, r *http.Request) {
		queries := r.URL.Query()
		if queries.Has("name") {
			name := r.URL.Query().Get("name")
			name = strings.ToLower(name)
			http.Redirect(w, r, fmt.Sprintf("/%s", name), http.StatusSeeOther)
		}
	})

	http.HandleFunc("/houseware", func(w http.ResponseWriter, r *http.Request) {
		queries := r.URL.Query()
		if queries.Has("name") {
			name := r.URL.Query().Get("name")
			name = strings.ToLower(name)
			strings.ReplaceAll(name, " ", "")
			http.Redirect(w, r, fmt.Sprintf("/%s"+"_houseware", name), http.StatusSeeOther)
		}
	})

	http.HandleFunc("/housewaresList", func(w http.ResponseWriter, r *http.Request) {
		templ.ExecuteTemplate(w, "housewaresList.gohtml", housewares)
	})

	fmt.Println("Server started on localhost:8010")
	err := http.ListenAndServe(":8010", nil)
	if err != nil {
		log.Fatal(err)
	}
}
