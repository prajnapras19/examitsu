package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/prajnapras19/project-form-exam-sman2/backend/config"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
	"google.golang.org/api/option"
)

type Service interface {
	GetUploadURL(req *GetUploadURLRequest) (*GetUploadURLResponse, error)
	UploadWithSignedURL(signedURL string, content []byte, contentType string) error
}

type service struct {
	cfg config.StorageConfig
}

func NewService(
	cfg config.StorageConfig,
) Service {
	return &service{
		cfg: cfg,
	}
}

func (s *service) GetUploadURL(req *GetUploadURLRequest) (*GetUploadURLResponse, error) {
	client, err := storage.NewClient(context.Background(), option.WithCredentialsFile(s.cfg.ServiceAccountKeyPath))
	if err != nil {
		log.Println("[storage][service][GetUploadURL] failed to initialize gcs client:", err.Error())
		return nil, lib.ErrFailedToGetUploadURL
	}
	defer client.Close()

	signedUrl, err := client.Bucket(s.cfg.BucketName).SignedURL(req.FileName, &storage.SignedURLOptions{
		Method:      http.MethodPut,
		Expires:     time.Now().Add(s.cfg.UploadURLExpiryDuration),
		ContentType: req.FileType,
	})
	if err != nil {
		log.Println("[storage][service][GetUploadURL] failed to get signed url:", err.Error())
		return nil, lib.ErrFailedToGetUploadURL
	}
	return &GetUploadURLResponse{
		UploadURL: signedUrl,
		PublicURL: fmt.Sprintf("https://storage.googleapis.com/%s/%s", s.cfg.BucketName, req.FileName),
	}, nil
}

func (s *service) UploadWithSignedURL(signedURL string, content []byte, contentType string) error {
	req, err := http.NewRequest(http.MethodPut, signedURL, bytes.NewReader(content))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", contentType)

	// Perform the HTTP request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed: status %d, body: %s", resp.StatusCode, body)
	}
	return nil
}
