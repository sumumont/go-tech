package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-tech/internal/logging"
	"go-tech/internal/service/book"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	r.Use(cors.Default())
	r.RedirectTrailingSlash = false

	r.Use(logging.Trace(), loggerHandler())
	r.Use(gin.Recovery())
	addBookServer(r)
	return r
}

func addBookServer(r *gin.Engine) {

	group := r.Group("/book")
	s := &book.ServerBook{}
	group.GET("", wrapper(s.ListBook)) // //上传模型
	group.POST("", wrapper(s.AddBook)) // //上传模型
}
