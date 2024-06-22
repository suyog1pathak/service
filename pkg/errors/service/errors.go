package errors

import (
	apiv1generic "github.com/suyog1pathak/services/api/v1/generic"
	apiv1healthcheck "github.com/suyog1pathak/services/api/v1/healthcheck"
	"net/http"
	"time"
)

func ServiceErrorHandler(errorString string) (apiv1generic.ErrorResponse, int) {
	switch errorString {
	case ErrInvalidInput:
		response := apiv1generic.ErrorResponse{
			Message: ErrInvalidInput,
			Error:   ErrInvalidInput,
		}
		return response, http.StatusInternalServerError
	case ErrServiceNotFound:
		response := apiv1generic.ErrorResponse{
			Message: "service not found.",
			Error:   ErrServiceNotFound,
		}
		return response, http.StatusNotFound
	case ErrServiceFoundWithSameName:
		response := apiv1generic.ErrorResponse{
			Message: "service with same name found, please use update api.",
			Error:   ErrServiceFoundWithSameName,
		}
		return response, http.StatusBadRequest
	case ErrServiceWithVersionNotFound:
		response := apiv1generic.ErrorResponse{
			Message: "service with name and version not found.",
			Error:   ErrServiceWithVersionNotFound,
		}
		return response, http.StatusNotFound
	}

	// default
	return apiv1generic.ErrorResponse{
		Message: ErrInternalServer,
		Error:   ErrInternalServer,
	}, http.StatusInternalServerError
}

func HealthcheckErrorHandler(errorString string) (apiv1healthcheck.Response, int) {
	status := apiv1healthcheck.Response{
		Status:     "healthy",
		StatusCode: http.StatusOK,
		Components: apiv1healthcheck.Components{
			Datastore: true,
		},
		Timestamp: time.Now(),
	}
	switch errorString {
	case ErrHealthcheckDbFailed:
		status.Components.Datastore = false
		status.Status = "unhealthy"
		status.StatusCode = http.StatusInternalServerError
	}
	return status, status.StatusCode
}
