package response

import "github.com/suyog1pathak/services/pkg/model"

type Service struct {
	*model.Service
	CurrentVersion int `json:"currentVersion"`
	TotalVersions  int `json:"totalVersion"`
} //@name ServiceResponse

type ServicePagination struct {
	Meta Meta
	Data []Service
} //@name ServicePagination

type Meta struct {
	Page         int `json:"page"`
	TotalResults int `json:"totalResults"`
	TotalPages   int `json:"totalPages"`
	PageSize     int `json:"pageSize"`
} //@name Meta
