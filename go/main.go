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

	println("\n")
	houses := getHouses()

	sort.Slice(characters, func(indexFirst, indexSecond int) bool {
		return characters[indexFirst].Name.NameEUen < characters[indexSecond].Name.NameEUen
	})
	//for _, chara := range characters {
	//	println("Villagers : ", chara.Name.NameEUen)
	//}
	//for _, hou := range houses {
	//println("Villagers : ", hou.Name)
	//println("Villagers T POSE : ", hou.NhDetails.ImageUrl)
	//println("House Interior : ", hou.NhDetails.HouseInteriorUrl)
	//println("House Exterior : ", hou.NhDetails.HouseExteriorUrl)
	//}
	for _, character := range characters {
		http.HandleFunc(fmt.Sprintf("/%s", strings.ToLower(character.Name.NameEUen)), func(writer http.ResponseWriter, request *http.Request) {
			name := strings.TrimPrefix(request.URL.Path, "/")
			simplifiedVillager := getSimplified(name, characters, houses)
			println("Villager : ", simplifiedVillager.Name)
			fmt.Println(request.URL, "url request")
			err := templ.ExecuteTemplate(writer, "character.gohtml", simplifiedVillager)
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
			fmt.Println(path, "=> introuvable")
		}
	})

	http.HandleFunc("/charalist", func(w http.ResponseWriter, r *http.Request) {
		for _, chara := range characters {
			print(chara.SelectedSpeccy)
			break
		}
		for i := range characters {
			characters[i].IsEmpty = true
		}
		path := TrimURLPrefix(r.URL.Path)
		switch r.Method {
		case "GET":

			templ.ExecuteTemplate(w, "charalist.gohtml", characters)
		case "POST":
			villager := r.FormValue("name")
			fmt.Println(villager)
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
				//fmt.Println("espèce choisie : " + species)
				//http.Redirect(w, r, species, http.StatusSeeOther)
				//fmt.Println("characters len = ", len(characters)-1)
				for i, chara := range characters {
					if chara.Species == species {
						//println("Villager : ", chara.Name.NameEUen, " speccy : ", chara.Species)
						characters[i].SelectedSpeccy = species
						//fmt.Println("debug speccy1")
						//fmt.Println(chara.SelectedSpeccy)
						//fmt.Println("debug speccy2")
					}
					fmt.Println("i = ", i)
					if i == len(characters)-1 {
						templ.ExecuteTemplate(w, "charalist.gohtml", characters)
						for i := range characters {
							characters[i].SelectedSpeccy = ""
						}
						//fmt.Println("charalist F5")
					}
					i++
				}
			}
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
