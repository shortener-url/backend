package service

import "github.com/shortener-url/backend/repository"

type Service struct {
	Repository *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{Repository: repo}
}
