package dto

type Response struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
}

