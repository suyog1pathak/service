package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/suyog1pathak/services/api/v1/generic"
	apiv1 "github.com/suyog1pathak/services/api/v1/response"
	"github.com/suyog1pathak/services/internal/service"
	log "github.com/suyog1pathak/services/pkg/logger"
	"github.com/suyog1pathak/services/pkg/model"
	"github.com/suyog1pathak/services/pkg/util"
	"net/http"
)

// CreateService
//
//	@BasePath		/api/v1/
//	@Summary		List services
//	@Description	List services
//	@Tags			services
//	@Accept			json
//	@Param			create	service	body	model.Service	true	"Add Service"
//	@Produce		application/json
//	@Success		201	{object}	apiv1.Service{}
//	@Failure		400	{object}	generic.ErrorResponse
//	@Failure		500	{object}	generic.ErrorResponse
//	@Router			/api/v1/services [post]
func CreateService(c *gin.Context) {
	_ = apiv1.Service{}
	reqBodyPtr, _ := c.Get("requestBody")
	//:TODO
	reqBody, _ := reqBodyPtr.(*model.Service)
	log.Info("received a request to create a service.", "body", util.StructToJson(reqBody))
	response, err := service.Create(reqBody)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusCreated, response)
}

// UpdateService
//
//	@BasePath		/api/v1/
//	@Summary		update service // create new version
//	@Description	update service // create new version
//	@Tags			services
//	@Accept			json
//	@Param			name	path	string	true			"service name"
//	@Param			update	service	body	model.Service	true	"update Service"
//	@Produce		application/json
//	@Success		201	{object}	apiv1.Service{}
//	@Failure		400	{object}	generic.ErrorResponse
//	@Failure		500	{object}	generic.ErrorResponse
//	@Router			/api/v1/services/{name} [patch]
func UpdateService(c *gin.Context) {
	reqBodyPtr, _ := c.Get("requestBody")
	reqBody, _ := reqBodyPtr.(*model.Service)
	name := c.Param("name")
	reqBody.Name = name
	response, err := service.CreateVersion(reqBody)
	log.Info("received a request to create a version for the service.", "body", util.StructToJson(reqBody))
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusCreated, response)
}

// UpdateServiceVersion
//
//	@BasePath		/api/v1/
//	@Summary		update service version
//	@Description	update service version
//	@Tags			services
//	@Accept			json
//	@Param			name	path	string	true			"service name"
//	@Param			version	path	string	true			"version"
//	@Param			update	service	body	model.Service	true	"update Service"
//	@Produce		application/json
//	@Success		201	{object}	model.Service
//	@Failure		400	{object}	generic.ErrorResponse
//	@Failure		500	{object}	generic.ErrorResponse
//	@Router			/api/v1/services/{name}/{version} [patch]
func UpdateServiceVersion(c *gin.Context) {
	reqBodyPtr, _ := c.Get("requestBody")
	reqBody, _ := reqBodyPtr.(*model.Service)
	name := c.Param("name")
	log.Info("received a request to update the existing version of the service.", "name", name)
	versionStr := c.Param("version")
	version, err := util.StringToInt(versionStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, generic.ErrorResponse{
			Message: "version is not a int",
			Error:   err.Error(),
		})
		return
	}
	reqBody.Name = name
	reqBody.Version = version
	response, err := service.UpdateVersion(reqBody)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusCreated, response)
}

// GetServiceByName
//
//	@BasePath		/api/v1/
//	@Summary		Get services by name
//	@Description	Get services by name
//	@Tags			services
//	@Accept			json
//	@Param			name	path	string	true	"service name"
//	@Produce		application/json
//	@Success		200	{object}	[]model.Service
//	@Failure		500	{object}	generic.ErrorResponse
//	@Failure		400	{object}	generic.ErrorResponse
//	@Failure		404	{object}	generic.ErrorResponse
//	@Router			/api/v1/services/{name} [get]
func GetServiceByName(c *gin.Context) {
	name := c.Param("name")
	log.Info("received a request to list all existing versions of the service.", "name", name)
	response, err := service.FetchByName(name)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, response)
}

// GetServiceNameAndVersion
//
//	@BasePath		/api/v1/
//	@Summary		Get service by version and name
//	@Description	Get service by version and name
//	@Tags			services
//	@Accept			json
//	@Param			name	path	string	true	"service name"
//	@Param			version	path	string	true	"version name"
//	@Produce		application/json
//	@Success		200	{object}	model.Service
//	@Failure		400	{object}	generic.ErrorResponse
//	@Failure		404	{object}	generic.ErrorResponse
//	@Failure		500	{object}	generic.ErrorResponse
//	@Router			/api/v1/services/{name}/{version} [get]
func GetServiceNameAndVersion(c *gin.Context) {
	name := c.Param("name")
	versionStr := c.Param("version")
	version, err := util.StringToInt(versionStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, generic.ErrorResponse{
			Message: "version is not a int",
			Error:   err.Error(),
		})
		return
	}
	log.Info("received a request to describe the service version.", "name", name, "version", version)
	response, err := service.FetchByVersionAndName(name, version)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, response)
}

// DeleteService
//
//	@BasePath		/api/v1/
//	@Summary		delete service
//	@Description	delete service
//	@Tags			services
//	@Accept			json
//	@Param			name	path	string	true	"service name"
//	@Produce		application/json
//	@Success		202	{object}	generic.Response
//	@Failure		400	{object}	generic.ErrorResponse
//	@Failure		404	{object}	generic.ErrorResponse
//	@Failure		500	{object}	generic.ErrorResponse
//	@Router			/api/v1/services/{name}/ [delete]
func DeleteService(c *gin.Context) {
	name := c.Param("name")
	log.Info("received a request to delete the service.", "name", name)
	err := service.Delete(name)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusAccepted, generic.Response{Message: fmt.Sprintf("service %s accepetd for deletion.", name)})
}

func GetAllServices(c *gin.Context) {
	if c.Request.URL.RawQuery == "" {
		getAllServices(c)
	} else {
		SearchAndSortServices(c)
	}
}

func getAllServices(c *gin.Context) {
	log.Info("received a request to get all services.")
	res, err := service.FetchAll()
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, res)
}

// SearchAndSortServices
//
//	@BasePath		/api/v1/
//	@Summary		list and filter services with pagination
//	@Description	list and filter services with pagination
//	@Tags			services
//	@Accept			json
//	@Produce		application/json
//	@Param			query		query		string	false	"query to filter, service name only"	example(Service1)
//	@Param			dir			query		string	false	"direction of sorting"	Enums(desc, asc)
//	@Param			page		query		int		false	"page no"				minimum(1)	maximum(1000)
//	@Param			sort		query		string	false	"sort by column name"	Enums(created_at, updated_at, version)
//	@Param			pagesize	query		int		false	"page size"				minimum(1)	maximum(10)
//	@Success		200			{object}	apiv1.ServicePagination
//	@Failure		500			{object}	generic.ErrorResponse
//	@Router			/api/v1/services [get]
func SearchAndSortServices(c *gin.Context) {
	log.Info("received a request to get all services with filters.")
	res, err := service.SearchAndSort(
		c.GetString("query"),
		c.GetString("sortBy"),
		c.GetString("dir"),
		c.GetInt("page"),
		c.GetInt("pageSize"))
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, res)
}
