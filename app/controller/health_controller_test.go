package controller

import (
	"bytes"
	"encoding/json"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant/ecodes"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/model"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/service"
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/errors"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/logging"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthController_HealthCheck(t *testing.T) {
	t.Run("Failed to call healthService.HealthCheck", func(t *testing.T) {
		expected := map[string]interface{}{
				"code":    ecodes.InternalServerError,
				"message": constant.ServerIssue,
				"success": false,
		}

		payload := []byte(``)
		req, _ := http.NewRequest("GET", "/health", bytes.NewBuffer(payload))
		req = req.WithContext(context.WithValue(req.Context(), logging.DiscardLoggingKey{}, true))
		recorder := httptest.NewRecorder()

		generalConfig := &config.GeneralConfig{}

		healthService := &service.MockHealthServiceInterface{}
		healthService.On("HealthCheck", req.Context()).Return(errors.New(ecodes.InternalServerError, constant.ServerIssue))

		logger := logging.InitializeLogger(
			generalConfig.Log.LogLevel,
			generalConfig.Log.LogFilePath,
			generalConfig.Log.JsonLogFormat,
			generalConfig.Log.LogRotation)

		baseController := BaseController{Config: generalConfig, Logging: logger}
		healthController := HealthController{
			BaseController: baseController,
			HealthService: healthService,
		}

		healthController.HealthCheck(recorder, req)
		code := recorder.Result().StatusCode
		if code == http.StatusOK {
			t.Errorf("Expected not to return response code 200. Got '%v'", code)
		}

		var errResp model.ErrorResponse
		json.Unmarshal(recorder.Body.Bytes(), &errResp)

		success := expected["success"].(bool)
		if errResp.Success == true{
			t.Errorf("Expected success attribute of response body to be '%v'. Got '%v'", success, errResp.Success )
		}

		errorCode := expected["code"].(uint32)
		if errResp.Error.Code != errorCode {
			t.Errorf("Expected error code to be '%v'. Got '%v'",errorCode, errResp.Error.Code)
		}

		if errResp.Error.Message != expected["message"] {
			t.Errorf("Expected error message to be '%v'. Got '%v'", expected["message"], errResp.Error.Message)
		}

		healthController.HealthService.(*service.MockHealthServiceInterface).AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		expected := map[string]interface{}{
			"code":    ecodes.InternalServerError,
			"message": constant.ServerIssue,
			"success": true,
		}

		payload := []byte(``)
		req, _ := http.NewRequest("GET", "/health", bytes.NewBuffer(payload))
		req = req.WithContext(context.WithValue(req.Context(), logging.DiscardLoggingKey{}, true))
		recorder := httptest.NewRecorder()

		generalConfig := &config.GeneralConfig{}


		healthService := &service.MockHealthServiceInterface{}
		healthService.On("HealthCheck", req.Context()).Return(nil)

		logger := logging.InitializeLogger(generalConfig)
		baseController := BaseController{Config: generalConfig, Logging: logger}
		healthController := HealthController{
			BaseController: baseController,
			HealthService: healthService,
		}

		healthController.HealthCheck(recorder, req)
		code := recorder.Result().StatusCode
		if code != http.StatusOK {
			t.Errorf("Expected response code 200. Got '%v'", code)
		}

		var resp model.SuccessResponse
		json.Unmarshal(recorder.Body.Bytes(), &resp)

		success := expected["success"].(bool)
		if resp.Success != success {
			t.Errorf("Expected to be '%v'. Got '%v'", success, resp.Success)
		}

		if resp.Data == nil{
			t.Errorf("Expected not to be '%v'. Got '%v'", nil, resp.Data)
		}

		healthController.HealthService.(*service.MockHealthServiceInterface).AssertExpectations(t)
	})
}

func TestHealthController_DBHealthCheckon(t *testing.T) {
	t.Run("Failed to call healthService.DBHealthCheck", func(t *testing.T) {
		expected := map[string]interface{}{
			"code":    ecodes.InternalServerError,
			"message": constant.ServerIssue,
			"success": false,
		}

		payload := []byte(``)
		req, _ := http.NewRequest("GET", "/health/db", bytes.NewBuffer(payload))
		req = req.WithContext(context.WithValue(req.Context(), logging.DiscardLoggingKey{}, true))
		recorder := httptest.NewRecorder()

		generalConfig := &config.GeneralConfig{}

		healthService := &service.MockHealthServiceInterface{}
		healthService.On("DBHealthCheck", req.Context()).Return("", errors.New(ecodes.InternalServerError, constant.ServerIssue))

		logger := logging.InitializeLogger(generalConfig)
		baseController := BaseController{Config: generalConfig, Logging: logger}
		healthController := HealthController{
			BaseController: baseController,
			HealthService: healthService,
		}

		healthController.DBHealthCheck(recorder, req)
		code := recorder.Result().StatusCode
		if code == http.StatusOK {
			t.Errorf("Expected not to return response code 200. Got '%v'", code)
		}

		var errResp model.ErrorResponse
		json.Unmarshal(recorder.Body.Bytes(), &errResp)

		success := expected["success"].(bool)
		if errResp.Success == true{
			t.Errorf("Expected success attribute of response body to be '%v'. Got '%v'", success, errResp.Success )
		}

		errorCode := expected["code"].(uint32)
		if errResp.Error.Code != errorCode {
			t.Errorf("Expected error code to be '%v'. Got '%v'",errorCode, errResp.Error.Code)
		}

		if errResp.Error.Message != expected["message"] {
			t.Errorf("Expected error message to be '%v'. Got '%v'", expected["message"], errResp.Error.Message)
		}

		healthController.HealthService.(*service.MockHealthServiceInterface).AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		expected := map[string]interface{}{
			"success": true,
			"database_name": "api_backend_dev",
		}

		payload := []byte(``)
		req, _ := http.NewRequest("GET", "/health/db", bytes.NewBuffer(payload))
		req = req.WithContext(context.WithValue(req.Context(), logging.DiscardLoggingKey{}, true))
		recorder := httptest.NewRecorder()

		generalConfig := &config.GeneralConfig{}


		healthService := &service.MockHealthServiceInterface{}
		healthService.On("DBHealthCheck", req.Context()).Return(expected["database_name"], nil)

		logger := logging.InitializeLogger(generalConfig)
		baseController := BaseController{Config: generalConfig, Logging: logger}
		healthController := HealthController{
			BaseController: baseController,
			HealthService: healthService,
		}

		healthController.DBHealthCheck(recorder, req)
		code := recorder.Result().StatusCode
		if code != http.StatusOK {
			t.Errorf("Expected response code 200. Got '%v'", code)
		}

		var resp model.SuccessResponse
		json.Unmarshal(recorder.Body.Bytes(), &resp)

		success := expected["success"].(bool)
		if resp.Success != success {
			t.Errorf("Expected to be '%v'. Got '%v'", success, resp.Success)
		}

		if resp.Data == nil{
			t.Errorf("Expected not to be '%v'. Got '%v'", nil, resp.Data)
		}

		var healthStatus model.HealthStatus
		data, _ := json.Marshal(resp.Data)
		json.Unmarshal(data, &healthStatus)

		if healthStatus.Environment != expected["database_name"]{
			t.Errorf("Expected to be '%v'. Got '%v'", expected["database_name"], healthStatus.Environment)
		}

		healthController.HealthService.(*service.MockHealthServiceInterface).AssertExpectations(t)
	})
}