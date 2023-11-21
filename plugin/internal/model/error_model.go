package model

type ErrorModel struct {
	Code   string `json:"code"`
	Status string `json:"status"`
	Detail string `json:"detail"`
}