package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Houses struct {
	Name          string   `json:"name"`
	Url           string   `json:"url"`
	AltName       string   `json:"alt_name"`
	TitleColor    string   `json:"title_color"`
	TextColor     string   `json:"text_color"`
	Id            string   `json:"id"`
	ImageUrl      string   `json:"image_url"`
	Species       string   `json:"species"`
	Personality   string   `json:"personality"`
	Gender        string   `json:"gender"`
	BirthdayMonth string   `json:"birthday_month"`
	BirthdayDay   string   `json:"birthday_day"`
	Sign          string   `json:"sign"`
	Quote         string   `json:"quote"`
	Phrase        string   `json:"phrase"`
	Clothing      string   `json:"clothing"`
	Islander      bool     `json:"islander"`
	Debut         string   `json:"debut"`
	PrevPhrases   []string `json:"prev_phrases"`
	NhDetails     *struct {
		ImageUrl          string   `json:"image_url"`
		PhotoUrl          string   `json:"photo_url"`
		IconUrl           string   `json:"icon_url"`
		Quote             string   `json:"quote"`
		SubPersonality    string   `json:"sub-personality"`
		Catchphrase       string   `json:"catchphrase"`
		Clothing          string   `json:"clothing"`
		ClothingVariation string   `json:"clothing_variation"`
		FavStyles         []string `json:"fav_styles"`
		FavColors         []string `json:"fav_colors"`
		Hobby             string   `json:"hobby"`
		HouseInteriorUrl  string   `json:"house_interior_url"`
		HouseExteriorUrl  string   `json:"house_exterior_url"`
		HouseWallpaper    string   `json:"house_wallpaper"`
		HouseFlooring     string   `json:"house_flooring"`
		HouseMusic        string   `json:"house_music"`
		HouseMusicNote    string   `json:"house_music_note"`
	} `json:"nh_details"`
	Appearances []string `json:"appearances"`
}

func getHouses() []Houses {
	url := "https://api.nookipedia.com/villagers?nhdetails=true&game=nh&api_key=72f07541-84d2-4033-a088-93efa4297f71"

	httpClient := http.Client{}

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

	var response []Houses
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return response
}
