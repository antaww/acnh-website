package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

type RawData struct {
	Id       int    `json:"id"`
	FileName string `json:"file-name"`
	Name     struct {
		NameUSen string `json:"name-USen"`
		NameEUen string `json:"name-EUen"`
		NameEUde string `json:"name-EUde"`
		NameEUes string `json:"name-EUes"`
		NameUSes string `json:"name-USes"`
		NameEUfr string `json:"name-EUfr"`
		NameUSfr string `json:"name-USfr"`
		NameEUit string `json:"name-EUit"`
		NameEUnl string `json:"name-EUnl"`
		NameCNzh string `json:"name-CNzh"`
		NameTWzh string `json:"name-TWzh"`
		NameJPja string `json:"name-JPja"`
		NameKRko string `json:"name-KRko"`
		NameEUru string `json:"name-EUru"`
	} `json:"name"`
	Personality       string `json:"personality"`
	BirthdayString    string `json:"birthday-string"`
	Birthday          string `json:"birthday"`
	Species           string `json:"species"`
	Gender            string `json:"gender"`
	Subtype           string `json:"subtype"`
	Hobby             string `json:"hobby"`
	CatchPhrase       string `json:"catch-phrase"`
	IconUri           string `json:"icon_uri"`
	ImageUri          string `json:"image_uri"`
	BubbleColor       string `json:"bubble-color"`
	TextColor         string `json:"text-color"`
	Saying            string `json:"saying"`
	CatchTranslations struct {
		CatchUSen string `json:"catch-USen"`
		CatchEUen string `json:"catch-EUen"`
		CatchEUde string `json:"catch-EUde"`
		CatchEUes string `json:"catch-EUes"`
		CatchUSes string `json:"catch-USes"`
		CatchEUfr string `json:"catch-EUfr"`
		CatchUSfr string `json:"catch-USfr"`
		CatchEUit string `json:"catch-EUit"`
		CatchEUnl string `json:"catch-EUnl"`
		CatchCNzh string `json:"catch-CNzh"`
		CatchTWzh string `json:"catch-TWzh"`
		CatchJPja string `json:"catch-JPja"`
		CatchKRko string `json:"catch-KRko"`
		CatchEUru string `json:"catch-EUru"`
	} `json:"catch-translations"`
}

func (rawdata *RawData) toData() Data {
	return Data{
		Name:        rawdata.Name.NameEUfr,
		Icon:        rawdata.IconUri,
		Image:       rawdata.ImageUri,
		Catch:       rawdata.CatchTranslations.CatchEUfr,
		BubbleColor: rawdata.BubbleColor,
		TextColor:   rawdata.TextColor,
	}
}

type Data struct {
	Name        string
	Icon        string
	Image       string
	Catch       string
	BubbleColor string
	TextColor   string
}

func getCharacters() []RawData {
	url := "https://acnhapi.com/v1a/villagers/"

	httpClient := http.Client{
		Timeout: time.Second * 2, // define timeout
	}

	//create request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "antaww")

	//make api call
	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(res.Body)
	}

	//parse response
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var response []RawData
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return response
}

func acnh(name string, response []RawData) Data {
	var index int
	for i, data := range response {
		if strings.ToLower(data.Name.NameEUfr) == strings.ToLower(name) {
			index = i
		}
	}
	return response[index].toData()
}

func characterExistence(name string, response []RawData) bool {
	for _, data := range response {
		if strings.ToLower(data.Name.NameEUfr) == strings.ToLower(name) {
			fmt.Println("dans func : perso existant")
			return true
		}
	}
	fmt.Println("dans func : perso inexistant")
	return false
}

type ErrStruct struct {
	Path string
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	templ := template.Must(template.ParseFiles("error.gohtml"))
	w.WriteHeader(status)
	if status == http.StatusNotFound {

		errData := &ErrStruct{Path: r.URL.Path}
		err := templ.Execute(w, errData)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	templ := template.Must(template.ParseFiles("character.gohtml")) //define html file

	characters := getCharacters()
	array2 := []string{} //debug
	for _, character := range characters {
		array2 = append(array2, character.Name.NameEUfr) //debug

		http.HandleFunc(fmt.Sprintf("/%s", strings.ToLower(character.Name.NameEUfr)), func(writer http.ResponseWriter, request *http.Request) {
			name := strings.TrimPrefix(request.URL.Path, "/")
			ch := acnh(name, characters)
			fmt.Println(request.URL, "=", ch)
			err := templ.Execute(writer, ch)
			if err != nil {
				log.Fatal(err)
			}
		})
	}

	//debug print dans l'ordre alphabétique
	sort.Strings(array2)
	for _, character := range array2 {
		fmt.Println(character)
	}
	//fin debug
	//test
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		path := TrimURLPrefix(r.URL.Path)
		if path == "favicon.ico" {
			return
		}
		if !characterExistence(path, characters) {
			errorHandler(w, r, http.StatusNotFound)
		}
	})

	handleDirectory(".", "/static/")
	handleDirectory("./fonts", "/fonts/")

	err := http.ListenAndServe(":8010", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleDirectory(directory, route string) {
	fonts := http.FileServer(http.Dir(directory))
	http.Handle(route, http.StripPrefix(route, fonts))
}
