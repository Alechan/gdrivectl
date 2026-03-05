package service

import (
	"context"

	"github.com/Alechan/gdrivectl/internal/auth"
	"github.com/Alechan/gdrivectl/internal/googleapi"
)

type FileMetaRequest struct {
	ID     string `json:"id"`
	Fields string `json:"fields,omitempty"`
}

type FileMetaService struct {
	tokens auth.TokenProvider
	drive  googleapi.DriveClient
}

func NewFileMetaService(tokens auth.TokenProvider, drive googleapi.DriveClient) *FileMetaService {
	return &FileMetaService{tokens: tokens, drive: drive}
}

func (s *FileMetaService) Run(ctx context.Context, req FileMetaRequest) (map[string]any, error) {
	token, err := s.tokens.AccessToken(ctx)
	if err != nil {
		return nil, err
	}
	return s.drive.FileMeta(ctx, token, googleapi.FileMetaRequest(req))
}
