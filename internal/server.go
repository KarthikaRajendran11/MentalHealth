// build an empty server struct
// import packages required to connect to postgres DB
// use them to perform an insert
package server

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Writer interface {
	Write(url string) error
}

type Service struct {
	// should accept an interface to write stuff to postgres DB
	writer Writer
}

func (s *Service) RegisterRoutes(r gin.IRoutes) {
	r.POST("history", s.handleWrite)
}

func NewService(writer Writer) *Service {
	return &Service{
		writer: writer,
	}
}

func (s *Service) handleWrite(c *gin.Context) {

	urlBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Fprintf(os.Stdout, err.Error())
		c.Status(http.StatusBadRequest)
		return
	}
	fmt.Fprintf(os.Stdout, "url : %s\n", string(urlBytes))

	err = s.writer.Write(string(urlBytes))
	if err != nil {
		fmt.Fprintf(os.Stdout, err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
	return
}
