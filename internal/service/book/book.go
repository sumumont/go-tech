package book

import "github.com/gin-gonic/gin"

type ServerBook struct {
}

func (s *ServerBook) AddBook(c *gin.Context) (interface{}, error) {
	return "hello world", nil
}
func (s *ServerBook) ListBook(c *gin.Context) (interface{}, error) {
	return []string{"a", "b", "hello world"}, nil
}
