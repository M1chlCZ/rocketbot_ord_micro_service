package models

type ErrorHTTP struct {
	ErrorMessage string `json:"errorMessage"`
	Status       string `json:"status"`
	HasError     bool   `json:"hasError"`
}

type HttpSuccess struct {
	HasError bool   `json:"hasError"`
	Message  string `json:"message"`
}

type TxRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type MintRequest struct {
	FeeRate int    `json:"feeRate"`
	Format  string `json:"format"`
	Base64  string `json:"base64"`
}

type MintLaunchRequest struct {
	ID      int64  `json:"id"`
	Addr    string `json:"addr"`
	FeeRate int    `json:"feeRate"`
	Format  string `json:"format"`
	Base64  string `json:"base64"`
}

type SendRequest struct {
	ID            int64  `json:"id"`
	FeeRate       int    `json:"feeRate"`
	Address       string `json:"Address"`
	InscriptionID string `json:"InscriptionID"`
}

type NewAddressRequest struct {
	Address string `json:"address"`
}

type ListInscriptionsResponse struct {
	HasError     bool      `json:"hasError"`
	Status       string    `json:"status"`
	Inscriptions []TxTable `json:"inscriptions"`
}

type NSFWInscriptionsResponse struct {
	HasError     bool        `json:"hasError"`
	Status       string      `json:"status"`
	Inscriptions []NSFWTable `json:"inscriptions"`
}

type FeeResponse struct {
	HasError bool   `json:"hasError"`
	FeeRate  int    `json:"feeRate"`
	Status   string `json:"status"`
}

type RawTxResponse struct {
	HasError bool           `json:"hasError"`
	RawTx    RawTransaction `json:"rawTx"`
	Status   string         `json:"status"`
}

type InscriptionPicResponse struct {
	HasError bool   `json:"hasError"`
	Base64   string `json:"base64"`
	Status   string `json:"status"`
}

type TestPicReq struct {
	Base64   string `json:"base64"`
	Filename string `json:"filename"`
}

type TestPicResponse struct {
	HasError bool   `json:"hasError"`
	Nsfw     bool   `json:"nsfwPicture"`
	NsfwText bool   `json:"nsfwText"`
	Status   string `json:"status"`
}
