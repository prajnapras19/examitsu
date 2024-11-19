package storage

import (
	"context"
	"fmt"
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
