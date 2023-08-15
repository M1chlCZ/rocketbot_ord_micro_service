package models

type Inscriptions struct {
	Inscription string `json:"inscription"`
	Location    string `json:"location"`
	Explorer    string `json:"explorer"`
}

type Inscribe struct {
	Commit      string `json:"commit"`
	Inscription string `json:"inscription"`
	Reveal      string `json:"reveal"`
	Fees        int    `json:"fees"`
}

type EstimateRequest struct {
	NumberOfBlocks int    `json:"blocks"`
	ImageURL       string `json:"imageUrl"`
}

type EstimateQualityRequest struct {
	UrlPic  string `json:"url_pic"`
	Quality string `json:"quality"`
	FeeRate int    `json:"fee_rate"`
}

type EstimateQualityResponse struct {
	Size      float64 `json:"size"`
	B64       string  `json:"base64"`
	BTCAmount float64 `json:"btc_amount"`
}
