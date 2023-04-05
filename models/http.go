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
	Format string `json:"format"`
	Base64 string `json:"base64"`
}

type SendRequest struct {
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
