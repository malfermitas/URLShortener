package service

import (
	"context"
	"urlshortener/internal/core/model"
	"urlshortener/internal/core/port/in"
)

type urlService struct {
}

func NewUrlService() in.URLService {
	return &urlService{}
}

func (u urlService) Create(ctx context.Context, originalURL string, customKey string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (u urlService) GetOriginal(ctx context.Context, shortKey string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (u urlService) GetAnalytics(ctx context.Context, shortKey string) (*model.Analytics, error) {
	//TODO implement me
	panic("implement me")
}
