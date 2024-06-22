package errors

const (
	ErrInvalidInput               = "error_invalid_input"
	ErrInternalServer             = "internal_server_error"
	ErrServiceNotFound            = "service_not_found"
	ErrServiceFoundWithSameName   = "service_found_with_the_same_name"
	ErrServiceWithVersionNotFound = "service_with_provided_name_and_version_not_found"
	ErrHealthcheckDbFailed        = "healthcheck_failed_db_connection_issue"
)
