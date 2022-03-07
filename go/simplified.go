package main

import "strings"

type simplifiedData struct {
	Name           string
	Icon           string
	Image          string
	Catch          string
	BubbleColor    string
	TextColor      string
	Saying         string
	Personality    string
	Hobby          string
	BirthString    string
	Birth          string
	Species        string
	Gender         string
	Subtype        string
	FavStyles      []string
	FavColors      []string
	VillagerBody   string
	HouseInterior  string
	HouseExterior  string
	HouseWallpaper string
	HouseFlooring  string
	HouseMusic     string
	HouseMusicNote string
}

func toData(villagerRawData VillagerRawData, houseRawData Houses) simplifiedData {
	return simplifiedData{
		Name:           villagerRawData.Name.NameEUen,
		Icon:           villagerRawData.IconUri,
		Image:          villagerRawData.ImageUri,
		Catch:          villagerRawData.CatchTranslations.CatchEUen,
		BubbleColor:    villagerRawData.BubbleColor,
		TextColor:      villagerRawData.TextColor,
		Saying:         villagerRawData.Saying,
		Personality:    villagerRawData.Personality,
		Hobby:          villagerRawData.Hobby,
		BirthString:    villagerRawData.BirthdayString,
		Birth:          villagerRawData.Birthday,
		Species:        villagerRawData.Species,
		Gender:         villagerRawData.Gender,
		Subtype:        villagerRawData.Subtype,
		FavStyles:      houseRawData.NhDetails.FavStyles,
		FavColors:      houseRawData.NhDetails.FavColors,
		VillagerBody:   houseRawData.NhDetails.ImageUrl,
		HouseInterior:  houseRawData.NhDetails.HouseInteriorUrl,
		HouseExterior:  houseRawData.NhDetails.HouseExteriorUrl,
		HouseWallpaper: houseRawData.NhDetails.HouseWallpaper,
		HouseFlooring:  houseRawData.NhDetails.HouseFlooring,
		HouseMusic:     houseRawData.NhDetails.HouseMusic,
		HouseMusicNote: houseRawData.NhDetails.HouseMusicNote,
	}
}

func getSimplified(name string, response []VillagerRawData, houseRawData []Houses) simplifiedData {
	var indexFirstApi int
	for i, data := range response {
		if strings.ToLower(data.Name.NameEUen) == strings.ToLower(name) {
			indexFirstApi = i
			break
		}
	}
	var indexSecondApi int
	for i, houseData := range houseRawData {
		if strings.ToLower(houseData.Name) == strings.ToLower(name) {
			indexSecondApi = i
			break
		}
	}
	return toData(response[indexFirstApi], houseRawData[indexSecondApi])
}
