package models

type ErrorWallet struct {
	Result interface{} `json:"result"`
	Error  Error       `json:"error"`
	ID     interface{} `json:"id"`
}
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
