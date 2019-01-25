package service

import (
	"context"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant/ecodes"
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"github.com/kyawmyintthein/golangRestfulAPISample/infrastructure"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/errors"
	"testing"
)

func TestHealthService_HealthCheck(t *testing.T) {
	t.Run("Heath file not exist", func(t *testing.T) {
		generalConfig := &config.GeneralConfig{}

		mockMongoStore := &infrastructure.MockMongoStore{}
		healthService := &HealthService{
			Config: generalConfig,
			MongoStore: mockMongoStore,
		}

		err := healthService.HealthCheck(context.Background())
		if err == nil{
			t.Errorf("Expected error shoudn't be nil. Got '%v'", err)
		}

		healthService.MongoStore.(*infrastructure.MockMongoStore).AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {

		generalConfig := &config.GeneralConfig{}

		mockMongoStore := &infrastructure.MockMongoStore{}
		healthService := &HealthService{
			Config: generalConfig,
			MongoStore: mockMongoStore,
		}

		err := healthService.HealthCheck(context.Background())
		if err != nil{
			t.Errorf("Expected error to be nil. Got '%v'", err)
		}


		healthService.MongoStore.(*infrastructure.MockMongoStore).AssertExpectations(t)
	})
}

func TestHealthService_DBHealthCheck(t *testing.T) {
	t.Run("Failed to get database name", func(t *testing.T) {
		expected := map[string]interface{}{
			"database_name":  "",
			"error":  errors.New(ecodes.DatabaseConnnectionFailed, constant.DatabaseConnnectionFailedErr),
		}

		generalConfig := &config.GeneralConfig{}

		mockMongoStore := &infrastructure.MockMongoStore{}
		mockMongoStore.On("DatabaseName").Return(expected["database_name"], expected["error"])
		healthService := &HealthService{
			Config: generalConfig,
			MongoStore: mockMongoStore,
		}

		databaseName, err := healthService.DBHealthCheck(context.Background())
		if err == nil{
			t.Errorf("Expected error shoudn't be 'nil'. Got '%v'", err)
		}

		clerror, ok := err.(errors.CustomError)
		if !ok{
			t.Errorf("Expected error type should be errors.CustomError type")
		}

		if clerror.GetCode() != ecodes.DatabaseConnnectionFailed{
			t.Errorf("Expected custom error code should be '%v'. Got '%v'", clerror.GetCode() , ecodes.DatabaseConnnectionFailed)
		}

		if databaseName != expected["database_name"]{
			t.Errorf("Expected database name should be empty. Got '%v'", databaseName)
		}

		healthService.MongoStore.(*infrastructure.MockMongoStore).AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		expected := map[string]interface{}{
			"database_name":  "api_backend_dev",
			"error":  nil,
		}
		generalConfig := &config.GeneralConfig{}

		mockMongoStore := &infrastructure.MockMongoStore{}
		mockMongoStore.On("DatabaseName").Return(expected["database_name"], expected["error"])
		healthService := &HealthService{
			Config: generalConfig,
			MongoStore: mockMongoStore,
		}

		databaseName, err := healthService.DBHealthCheck(context.Background())
		if err != nil{
			t.Errorf("Expected error to be 'nil'. Got '%v'", err)
		}

		if databaseName != expected["database_name"]{
			t.Errorf("Expected database name to be '%v'. Got '%v'", expected["database_name"], databaseName)
		}

		healthService.MongoStore.(*infrastructure.MockMongoStore).AssertExpectations(t)
	})
}
