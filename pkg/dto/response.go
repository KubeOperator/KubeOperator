package dto

type Response struct {
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}
