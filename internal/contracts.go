package server

type UploadOptions struct {
	FileName   string `json:"file_name"`
	BucketName string `json:"bucket_name"`
	Content    string `json:"content"`
}
