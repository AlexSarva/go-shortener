package models

// BadRequest message when respond bad request
type BadRequest struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}
