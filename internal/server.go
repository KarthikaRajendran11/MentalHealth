// build an empty server struct
// import packages required to connect to postgres DB
// use them to perform an insert
package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type Writer interface {
	Write(url ...string) error
}

type Uploader interface {
	Upload(ctx context.Context, fileName, bucket string, content []byte) error
}

type Service struct {
	// should accept an interface to write stuff to postgres DB
	writer   Writer
	uploader Uploader
}

func (s *Service) RegisterRoutes(r gin.IRoutes) {
	r.GET("_status", s.handleHealthCheck)
	r.POST("history", s.handleWrite)
	r.POST("upload", s.handleUpload)
}

func NewService(writer Writer, uploader Uploader) *Service {
	return &Service{
		writer:   writer,
		uploader: uploader,
	}
}

func (s *Service) handleHealthCheck(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (s *Service) handleUpload(c *gin.Context) {
	var request UploadOptions
	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Fprintln(os.Stdout, err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	fmt.Fprintln(os.Stdout, request.FileName)
	fmt.Fprintln(os.Stdout, request.BucketName)
	data, err := base64.StdEncoding.DecodeString(request.Content)
	if err != nil {
		fmt.Fprintln(os.Stdout, err.Error())
		c.Status(http.StatusBadRequest)
		return
	}
	err = s.uploader.Upload(c, request.FileName, request.BucketName, data)
	if err != nil {
		fmt.Fprintln(os.Stdout, err.Error())
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}

func (s *Service) handleWrite(c *gin.Context) {

	urlBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Fprintln(os.Stdout, err.Error())
		c.Status(http.StatusBadRequest)
		return
	}
	_s := string(urlBytes)
	fmt.Fprintln(os.Stdout, _s)
	content := strings.Split(_s, "url: ")
	fmt.Fprintf(os.Stdout, "len content : %d\n", len(content))
	fmt.Fprintf(os.Stdout, "content : %s\n", content[0])
	fmt.Fprintf(os.Stdout, "content : %s\n", content[1])
	urlEmail := strings.Split(content[1], ", email: ")
	fmt.Fprintf(os.Stdout, "urlEmail : %s\n", urlEmail[0])
	fmt.Fprintf(os.Stdout, "urlEmail : %s\n", urlEmail[1])
	url := urlEmail[0]
	email := urlEmail[1]

	err = s.writer.Write(url, email)
	if err != nil {
		fmt.Fprintln(os.Stdout, err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}
