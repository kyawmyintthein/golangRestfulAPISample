package controller

import (
	"bitbucket.org/libertywireless/go-restapi-boilerplate/app/model"
	"bitbucket.org/libertywireless/go-restapi-boilerplate/app/service"
	"net/http"
)

type HealthController struct{
	BaseController
	HealthService service.HealthServiceInterface
}


// HealthCheck godoc
// @Summary Check server's health
// @Description check health file exist or not
// @Tags  health
// @Accept  json
// @Produce  json
// @Success 200 {object} model.HealthStatus
// @Failure 500 {object} model.ErrorResponse
// @Router /health [get]
func (c *HealthController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	c.WriteJSON(r, w, http.StatusOK,
		model.HealthStatus{
		"ok",
		""})
	return
}

// DBHealthCheck godoc
// @Summary Check database's health
// @Description check database is working or not
// @Tags  health
// @Accept  json
// @Produce  json
// @Success 200 {object} model.HealthStatus
// @Failure 500 {object} model.ErrorResponse
// @Router /health/db [get]
func (c *HealthController) DBHealthCheck(w http.ResponseWriter, r *http.Request) {
	dbname, err := c.HealthService.DBHealthCheck(r.Context())
	if err != nil {
		c.WriteError(r, w, err)
		return
	}

	c.WriteJSON(r, w, http.StatusOK, model.HealthStatus{
		"ok",
		dbname})
	return
}
