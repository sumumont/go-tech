package book

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-tech/internal/dao"
	"go-tech/internal/dao/db"
	"go-tech/internal/dao/model"
	"go-tech/internal/dto"
	"go-tech/pkg/pke"
)

type ServerBook struct {
}

func (s *ServerBook) AddBook(c *gin.Context) (interface{}, error) {
	var param dto.Book
	err := c.ShouldBindJSON(&param)
	if err != nil {
		return nil, err
	}
	return nil, db.Transaction(c, func(runCtx context.Context) error {
		bookDao := dao.BookDao[model.Book]{}
		bookModel := model.Book{
			Name:   param.Name,
			Author: param.Author,
		}
		return bookDao.Save(runCtx, &bookModel)
	})

}
func (s *ServerBook) ListBook(c *gin.Context) (interface{}, error) {
	var param dto.Book
	err := c.ShouldBind(&param)
	if err != nil {
		return nil, err
	}
	var res pke.BaseListResp
	var books []model.Book
	var total int64
	err = db.Transaction(c, func(runCtx context.Context) error {
		bookDao := dao.BookDao[model.Book]{}
		bookModel := model.Book{
			Name:   param.Name,
			Author: param.Author,
		}
		books, total, err = bookDao.FindAll(runCtx, bookModel, &param.BaseList)
		res = pke.BaseListResp{
			Total: total,
			Items: books,
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
