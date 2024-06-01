package db

import (
	"fmt"
	"gorm.io/gorm"
)

type BaseList struct {
	PageNum  int    `json:"pageNum" form:"pageNum"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	Sort     string `json:"sort" form:"sort"`
}

func PageList(db *gorm.DB, baseList *BaseList) (*gorm.DB, error) {
	if baseList != nil {
		sort, err := NormalizeSorts(baseList.Sort)
		if err != nil {
			return nil, err
		}
		if sort == "" {
			sort = "created_at desc"
		}
		if len(sort) > 0 {
			db = db.Order(sort)
		}
		if baseList.PageNum > 0 && baseList.PageSize > 0 {
			db = db.Offset((baseList.PageNum - 1) * baseList.PageSize).Limit(baseList.PageSize)
		}
	}
	return db, nil
}

func NormalizeSorts(str string) (string, error) {
	var sorts string
	for _, v := range str {
		if v >= 'A' && v <= 'Z' {
			sorts += "_" + string(v+32)
		} else if v == '|' {
			sorts += " "
		} else if v >= 'a' && v <= 'z' || v == '_' || v == '-' || v == ',' {
			sorts += string(v)
		} else {
			return "", fmt.Errorf("invalid sort fields!")
		}
	}
	return sorts, nil
}
