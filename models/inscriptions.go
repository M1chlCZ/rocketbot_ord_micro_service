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
