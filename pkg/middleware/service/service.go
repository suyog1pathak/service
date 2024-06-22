package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/suyog1pathak/services/api/v1/generic"
	"github.com/suyog1pathak/services/pkg/model"
	"github.com/suyog1pathak/services/pkg/util"
	"net/http"
)

func ServiceBodyValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody model.Service
		// schema validation
		if err := c.BindJSON(&requestBody); err != nil {
			res := generic.ErrorResponse{
				Message: "request body validation failed",
				Error:   err.Error(),
			}
			c.IndentedJSON(http.StatusBadRequest, res)
			c.Abort()
			return
		}

		c.Set("requestBody", &requestBody)
		c.Next()
	}
}

func ServiceQueryParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.RawQuery != "" {
			// :todo validations
			query := c.DefaultQuery("query", "%")
			sort := c.DefaultQuery("sort", "created_at")
			dir := c.DefaultQuery("dir", "desc")
			page, _ := util.StringToInt(c.DefaultQuery("page", "1"))
			pageSize, _ := util.StringToInt(c.DefaultQuery("pagesize", "10"))

			if page <= 0 {
				page = 1
			}

			switch {
			// max page size
			case pageSize > 100:
				pageSize = 100
			// min page size
			case pageSize <= 0:
				pageSize = 1
			}
			c.Set("query", query)
			c.Set("sortBy", sort)
			c.Set("dir", dir)
			c.Set("page", page)
			c.Set("pageSize", pageSize)
		}
		c.Next()
	}
}
