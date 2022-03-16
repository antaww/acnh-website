package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type test struct {
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

type HousewareRawData struct {
	URL               string `json:"url"`
	Name              string `json:"name"`
	NameSimplified    string
	Category          string        `json:"category"`
	ItemSeries        string        `json:"item_series"`
	ItemSet           string        `json:"item_set"`
	HhaCategory       string        `json:"hha_category"`
	Tag               string        `json:"tag"`
	HhaBase           int           `json:"hha_base"`
	Lucky             bool          `json:"lucky"`
	LuckySeason       string        `json:"lucky_season"`
	Sell              int           `json:"sell"`
	VariationTotal    int           `json:"variation_total"`
	PatternTotal      int           `json:"pattern_total"`
	Customizable      bool          `json:"customizable"`
	CustomKits        int           `json:"custom_kits"`
	CustomKitType     string        `json:"custom_kit_type"`
	CustomBodyPart    string        `json:"custom_body_part"`
	CustomPatternPart string        `json:"custom_pattern_part"`
	Height            string        `json:"height"`
	DoorDecor         bool          `json:"door_decor"`
	VersionAdded      string        `json:"version_added"`
	Unlocked          bool          `json:"unlocked"`
	Notes             string        `json:"notes"`
	GridWidth         float64       `json:"grid_width"`
	GridLength        float64       `json:"grid_length"`
	Themes            []string      `json:"themes"`
	Functions         []interface{} `json:"functions"`
	Availability      []struct {
		From string `json:"from"`
		Note string `json:"note"`
	} `json:"availability"`
	Buy []struct {
		Price    int    `json:"price"`
		Currency string `json:"currency"`
	} `json:"buy"`
	Variations []struct {
		Variation string   `json:"variation"`
		Pattern   string   `json:"pattern"`
		ImageURL  string   `json:"image_url"`
		Colors    []string `json:"colors"`
	} `json:"variations"`
}

type SimplifiedHouseware struct {
	Name           string
	NameSimplified string
	ImageURL       []string
	Buy            []struct {
		Price    int
		Currency string
	}
	From    string
	Note    string
	Themes  []string
	Version string
}

func toDataHouseware(housewaredata HousewareRawData) SimplifiedHouseware {
	var images []string
	for i := range housewaredata.Variations {
		images = append(images, housewaredata.Variations[i].ImageURL)
	}
	var buy []struct {
		Price    int
		Currency string
	}
	type buyData struct {
		Price    int
		Currency string
	}
	for i := range housewaredata.Buy {
		data := buyData{Price: housewaredata.Buy[i].Price, Currency: housewaredata.Buy[i].Currency}
		//buy[i].Price = housewaredata.Buy[i].Price
		//buy[i].Currency = housewaredata.Buy[i].Currency
		buy = append(buy, data)
	}
	var from string
	for i := range housewaredata.Availability {
		from = housewaredata.Availability[i].From
	}
	var note string
	for i := range housewaredata.Availability {
		note = housewaredata.Availability[i].Note
	}
	var themes []string
	for i := range housewaredata.Themes {
		themes = append(themes, housewaredata.Themes[i])
	}
	return SimplifiedHouseware{
		Name:           housewaredata.Name,
		NameSimplified: housewaredata.NameSimplified,
		ImageURL:       images,
		Buy:            buy,
		From:           from,
		Note:           note,
		Themes:         themes,
		Version:        housewaredata.VersionAdded,
	}
}

func getSimplifiedHouseware(name string, response []HousewareRawData) SimplifiedHouseware {
	var index int
	for i := range response {
		nameToCheck := strings.ToLower(response[i].Name)
		nameToCheck = strings.ReplaceAll(nameToCheck, " ", "")
		if nameToCheck == name {
			index = i
			break
		}
	}
	return toDataHouseware(response[index])
}

func getHouseware() []HousewareRawData {
	url := "https://api.nookipedia.com/nh/furniture?api_key=72f07541-84d2-4033-a088-93efa4297f71"

	httpClient := http.Client{ // define timeout
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

	var response []HousewareRawData
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return response
}
