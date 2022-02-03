package main

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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
		Name:  rawdata.Name.NameEUfr,
		Icon:  rawdata.IconUri,
		Image: rawdata.ImageUri,
		Catch: rawdata.CatchTranslations.CatchEUfr,
	}
}

type Data struct {
	Name  string
	Icon  string
	Image string
	Catch string
}

func main() {

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

	//for _, character := range response {
	//	fmt.Println(character.Name, "gender -->", character.Gender)
	//}

	character := response
	println(character[0].IconUri)

	templ := template.Must(template.ParseFiles("index.gohtml")) //define html file
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var values []Data
		for _, data := range character {
			values = append(values, data.toData())
		}

		err := templ.Execute(w, values)
		if err != nil {
			log.Fatal(err)
		} //execute template
		//fmt.Println(values.Name)
	})

	css := http.FileServer(http.Dir("."))
	http.Handle("/static/", http.StripPrefix("/static/", css))

	err = http.ListenAndServe(":8010", nil)
	if err != nil {
		log.Fatal(err)
	}
}
