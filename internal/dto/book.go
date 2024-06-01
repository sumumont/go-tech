package dto

import "go-tech/internal/dao/db"

type Book struct {
	Name   string `json:"name" form:"name"`
	Author string `json:"author" form:"author"`
	db.BaseList
}
