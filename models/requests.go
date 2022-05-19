package models

type BadRequest struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}
