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

type VillagerRawData struct {
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

type HousewareRawData struct {
	Variant             interface{} `json:"variant"`
	BodyTitle           interface{} `json:"body-title"`
	Pattern             interface{} `json:"pattern"`
	PatternTitle        interface{} `json:"pattern-title"`
	IsDIY               bool        `json:"isDIY"`
	CanCustomizeBody    bool        `json:"canCustomizeBody"`
	CanCustomizePattern bool        `json:"canCustomizePattern"`
	KitCost             interface{} `json:"kit-cost"`
	Color1              string      `json:"color-1"`
	Color2              string      `json:"color-2"`
	Size                string      `json:"size"`
	Source              string      `json:"source"`
	SourceDetail        string      `json:"source-detail"`
	Version             string      `json:"version"`
	HhaConcept1         string      `json:"hha-concept-1"`
	HhaConcept2         interface{} `json:"hha-concept-2"`
	HhaSeries           interface{} `json:"hha-series"`
	HhaSet              interface{} `json:"hha-set"`
	IsInteractive       bool        `json:"isInteractive"`
	Tag                 string      `json:"tag"`
	IsOutdoor           bool        `json:"isOutdoor"`
	SpeakerType         interface{} `json:"speaker-type"`
	LightingType        interface{} `json:"lighting-type"`
	IsCatalog           bool        `json:"isCatalog"`
	FileName            string      `json:"file-name"`
	VariantId           interface{} `json:"variant-id"`
	InternalId          int         `json:"internal-id"`
	Name                struct {
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
	BuyPrice  int    `json:"buy-price"`
	SellPrice int    `json:"sell-price"`
	ImageUri  string `json:"image_uri"`
}

func (rawdata *HousewareRawData) toData() HousewareData {
	return HousewareData{
		SellPrice: rawdata.SellPrice,
		ImageUri:  rawdata.ImageUri,
	}
}

func (rawdata *VillagerRawData) toData() Data {
	return Data{
		Name:         rawdata.Name.NameEUen,
		Icon:         rawdata.IconUri,
		Image:        rawdata.ImageUri,
		Catch:        rawdata.CatchTranslations.CatchEUen,
		BubbleColor:  rawdata.BubbleColor,
		TextColor:    rawdata.TextColor,
		Saying:       rawdata.Saying,
		Personnality: rawdata.Personality,
		Hobby:        rawdata.Hobby,
		Birth:        rawdata.BirthdayString,
	}
}

type HousewareData struct {
	SellPrice int
	ImageUri  string
}

type Data struct {
	Name         string
	Icon         string
	Image        string
	Catch        string
	BubbleColor  string
	TextColor    string
	Saying       string
	Personnality string
	Hobby        string
	Birth        string
}

func getCharacters() []VillagerRawData {
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

	var response []VillagerRawData
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return response
}

func acnh(name string, response []VillagerRawData) Data {
	var index int
	for i, data := range response {
		if strings.ToLower(data.Name.NameEUen) == strings.ToLower(name) {
			index = i
		}
	}
	return response[index].toData()
}

func characterExistence(name string, response []VillagerRawData) bool {
	for _, data := range response {
		if strings.ToLower(data.Name.NameEUen) == strings.ToLower(name) {
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

func main() {
	templ := template.Must(template.ParseGlob("templates/*.gohtml")) //define gohtml file

	// Gestion de tous les fichiers css
	fs := http.FileServer(http.Dir("style"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

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

	////debug print dans l'ordre alphabétique
	//sort.Strings(array2)
	//for _, character := range array2 {
	//	fmt.Println(character)
	//}
	//fin debug
	//test
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := TrimURLPrefix(r.URL.Path)
		if path == "favicon.ico" {
			return
		}
		fmt.Println("index")
		if path == "" {
			templ.ExecuteTemplate(w, "index.gohtml", "")
		}
		if !characterExistence(path, characters) {
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

	err := http.ListenAndServe(":8010", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleDirectory(directory, route string) {
	fonts := http.FileServer(http.Dir(directory))
	http.Handle(route, http.StripPrefix(route, fonts))
}
