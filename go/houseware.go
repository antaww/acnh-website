package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

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
	Availability []struct {
		From string
		Note string
	}
	Versions []struct {
		ImageURL   string
		Colors     []string
		ColorsText []string
	}
	Themes       []string
	Category     string
	ItemSeries   string
	Tag          string
	Sell         int
	VersionAdded string
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
		buy = append(buy, data)
	}

	var availability []struct {
		From string
		Note string
	}
	type availabilityData struct {
		From string
		Note string
	}
	for i := range housewaredata.Availability {
		data := availabilityData{From: housewaredata.Availability[i].From, Note: housewaredata.Availability[i].Note}
		availability = append(availability, data)
	}

	var versions []struct {
		ImageURL   string
		Colors     []string
		ColorsText []string
	}
	type versionsData struct {
		ImageURL   string
		Colors     []string
		ColorsText []string
	}
	for i := range housewaredata.Variations {
		data := versionsData{ImageURL: housewaredata.Variations[i].ImageURL, Colors: housewaredata.Variations[i].Colors}
		versions = append(versions, data)
	}

	var themes []string
	for i := range housewaredata.Themes {
		themes = append(themes, housewaredata.Themes[i])
	}

	return SimplifiedHouseware{
		Name:           housewaredata.Name,
		NameSimplified: housewaredata.NameSimplified,
		Category:       housewaredata.Category,
		ItemSeries:     housewaredata.ItemSeries,
		Tag:            housewaredata.Tag,
		ImageURL:       images,
		Buy:            buy,
		Availability:   availability,
		Themes:         themes,
		Versions:       versions,
		Sell:           housewaredata.Sell,
		VersionAdded:   housewaredata.VersionAdded,
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
