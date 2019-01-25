package model

type HealthStatus struct {
	Status      string `json:"status,omitempty"  example:"ok"`
	Environment string `json:"environment"   example:"database name"`
}
