package controllers

import (
	"github.com/gin-gonic/gin"
	apiv1healthcheck "github.com/suyog1pathak/services/api/v1/healthcheck"
	"github.com/suyog1pathak/services/internal/healthcheck"
)

// Healthcheck
//
//	@Summary		healthcheck
//	@Description	healthcheck
//	@Tags			healthcheck
//	@Accept			json
//	@Produce		application/json
//	@Success		200	{object}	apiv1healthcheck.Response
//	@Failure		500	{object}	apiv1healthcheck.Response
//	@Router			/healthcheck [get]
func Healthcheck(c *gin.Context) {
	_ = apiv1healthcheck.Response{}
	err := healthcheck.CheckHealthCheckStatus()
	c.Error(err)
}

// LivenessCheck
//
//	@Summary		liveness
//	@Description	liveness
//	@Tags			healthcheck
//	@Accept			json
//	@Produce		application/json
//	@Success		200	{object}	apiv1healthcheck.Response
//	@Failure		500	{object}	apiv1healthcheck.Response
//	@Router			/liveness [get]
func LivenessCheck(c *gin.Context) {
	_ = apiv1healthcheck.Response{}
	err := healthcheck.CheckHealthCheckStatus()
	c.Error(err)
}

// ReadinessCheck
// @Summary		ReadinessCheck
// @Description	ReadinessCheck
// @Tags			healthcheck
// @Accept			json
// @Produce		application/json
// @Success		200	{object}	apiv1healthcheck.Response
// @Failure		500	{object}	apiv1healthcheck.Response
// @Router			/readiness [get]
func ReadinessCheck(c *gin.Context) {
	_ = apiv1healthcheck.Response{}
	err := healthcheck.CheckHealthCheckStatus()
	c.Error(err)
}
