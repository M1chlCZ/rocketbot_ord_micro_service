package models

type HttpErr struct {
	HasError bool   `json:"hasError"`
	Message  string `json:"errorMessage"`
}

type HttpSuccess struct {
	HasError bool   `json:"hasError"`
	Message  string `json:"message"`
}
