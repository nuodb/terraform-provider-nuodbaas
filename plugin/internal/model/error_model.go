/*
(C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/
package model

type ErrorModel struct {
	Code   string `json:"code"`
	Status string `json:"status"`
	Detail string `json:"detail"`
}