package main

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

type HousewareData struct {
	SellPrice int
	ImageUri  string
}

func (rawdata *HousewareRawData) toData() HousewareData {
	return HousewareData{
		SellPrice: rawdata.SellPrice,
		ImageUri:  rawdata.ImageUri,
	}
}
