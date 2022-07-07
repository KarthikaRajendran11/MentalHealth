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
	"github.com/pkg/errors"
)

type Writer interface {
	Write(url ...string) error
}

type Uploader interface {
	Upload(ctx context.Context, fileName, bucket string, content []byte) error
}

type Service struct {
	writer   Writer
	uploader Uploader
}

func (s *Service) RegisterRoutes(r gin.IRoutes) {
	r.GET("_status", s.handleHealthCheck)
	r.POST("history", s.handleWrite)
	r.POST("website", s.handleWebsite)
	r.POST("upload", s.handleUpload)
	r.OPTIONS("website", s.handleOptions)
}

func NewService(writer Writer, uploader Uploader) *Service {
	return &Service{
		writer:   writer,
		uploader: uploader,
	}
}

func (s *Service) handleOptions(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, User-Agent, Origin")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Status(http.StatusOK)
}

func (s *Service) handleHealthCheck(c *gin.Context) {
	c.Status(http.StatusOK)
}

// TODO: Secure recog-bucket-mental-helathco bucket to protect user privacy
// Set an expiry data on bucket so the user images are automatically deleted after certain amount of time
// Set ACL for bucket. As of July 6th 2022, the bucket has public access
// Must do above before production
func (s *Service) handleUpload(c *gin.Context) {
	var request UploadOptions
	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Fprintln(os.Stdout, errors.Wrap(err, "failed to bind json payload to request struct").Error())
		c.Status(http.StatusBadRequest)
		return
	}

	fmt.Fprintln(os.Stdout, request.FileName)
	fmt.Fprintln(os.Stdout, request.BucketName)
	data, err := base64.StdEncoding.DecodeString(request.Content)
	if err != nil {
		fmt.Fprintln(os.Stdout, errors.Wrap(err, "failed to base64 decode request content").Error())
		c.Status(http.StatusBadRequest)
		return
	}
	err = s.uploader.Upload(c, request.FileName, request.BucketName, data)
	if err != nil {
		fmt.Fprintln(os.Stdout, errors.Wrap(err, "failed to upload file to s3").Error())
		c.Status(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}

// TODO: Encrypt email (encryption at rest maybe) while storing in postgres DB to protect user identity
// Must to before production
func (s *Service) handleWebsite(c *gin.Context) {

	var request WebsiteOptions
	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Fprintln(os.Stdout, errors.Wrap(err, "failed to bind json payload to request struct").Error())
		c.Status(http.StatusBadRequest)
		return
	}

	url := request.URL
	email := request.Email

	err := s.writer.Write(url, email)
	if err != nil {
		fmt.Fprintln(os.Stdout, errors.Wrap(err, "failed to write to postgres DB").Error())
		c.Status(http.StatusBadRequest)
		return
	}
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	// access-control-allow-headers
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, User-Agent, Origin")
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
	urlEmail := strings.Split(content[1], ", email: ")
	url := urlEmail[0]
	email := urlEmail[1]

	err = s.writer.Write(url, email)
	if err != nil {
		fmt.Fprintln(os.Stdout, errors.Wrap(err, "failed to write to postgres DB").Error())
		c.Status(http.StatusBadRequest)
		return
	}
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, User-Agent, Origin")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Status(http.StatusOK)
}
