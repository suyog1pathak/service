package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/suyog1pathak/services/pkg/errors/service"
	log "github.com/suyog1pathak/services/pkg/logger"
)

func ServiceErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err != nil {
			response, responseCode := errors.ServiceErrorHandler(err.Error())
			c.IndentedJSON(responseCode, response)
			c.Abort()
		} else {
			log.Debug("no error at middleware")
		}
	}
}
