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
