package middleware

import (
	"github.com/gin-gonic/gin"
	customerror "github.com/suyog1pathak/services/pkg/errors/service"
)

func HealthcheckCatchErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		response, responseCode := customerror.HealthcheckErrorHandler(err.Error())
		c.IndentedJSON(responseCode, response)
	}
}
