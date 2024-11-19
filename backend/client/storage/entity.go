package storage

type GetUploadURLRequest struct {
	FileName string
	FileType string
}

type GetUploadURLResponse struct {
	UploadURL string
	PublicURL string
}
