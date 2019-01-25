package service

import (
	"context"
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"github.com/kyawmyintthein/golangRestfulAPISample/infrastructure"
)

type HealthServiceInterface interface {
	HealthCheck(ctx context.Context) error
	DBHealthCheck(ctx context.Context) (string, error)
}

type HealthService struct {
	Config *config.GeneralConfig
	MongoStore infrastructure.MongoStore
}


func (s *HealthService) HealthCheck(ctx context.Context) error{
	return nil
}

func (s *HealthService) DBHealthCheck(ctx context.Context) (string, error){
	return s.MongoStore.DatabaseName()
}

