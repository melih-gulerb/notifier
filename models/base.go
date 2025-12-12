package models

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var Success = "Başarılı"
var ServerError = "Sunucu hatası"
