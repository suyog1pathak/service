package service

import (
	"errors"
	apiv1 "github.com/suyog1pathak/services/api/v1/response"
	customerrors "github.com/suyog1pathak/services/pkg/errors/service"
	"github.com/suyog1pathak/services/pkg/model"
	"math"
)

func Create(service *model.Service) (apiv1.Service, error) {
	var response apiv1.Service
	service.Version = 1
	_, err := FetchByName(service.Name)
	if err != nil {
		if err.Error() == customerrors.ErrServiceNotFound {
			err = service.Add()
			if err != nil {
				return apiv1.Service{}, err
			}
			response = apiv1.Service{
				TotalVersions:  1,
				CurrentVersion: 1,
				Service:        service,
			}
			return response, nil
		}
		return apiv1.Service{}, err
	}
	return response, errors.New(customerrors.ErrServiceFoundWithSameName)
}

func CreateVersion(service *model.Service) (apiv1.Service, error) {
	var response apiv1.Service
	oldVersions, err := FetchByName(service.Name)
	if err != nil {
		if err.Error() == customerrors.ErrServiceNotFound {
			return apiv1.Service{}, err
		}
		return apiv1.Service{}, err
	}
	service.Version = oldVersions[len(oldVersions)-1].Version + 1
	err = service.Add()
	if err != nil {
		return apiv1.Service{}, err
	}
	response = apiv1.Service{
		TotalVersions:  len(oldVersions) + 1,
		CurrentVersion: service.Version,
		Service:        service,
	}
	return response, nil
}

func UpdateVersion(service *model.Service) (*model.Service, error) {
	_, err := FetchByVersionAndName(service.Name, service.Version)
	if err != nil {
		if err.Error() == customerrors.ErrServiceWithVersionNotFound {
			return service, err
		}
		return service, err
	}
	err = service.UpdateByNameAndVersion()
	if err != nil {
		return service, err
	}
	return service, nil
}

func Delete(name string) error {
	_, err := FetchByName(name)
	if err != nil {
		return err
	}
	service := &model.Service{}
	service.Name = name
	err = service.DeleteByName()
	if err != nil {
		return err
	}
	return nil
}

func FetchByName(name string) ([]model.Service, error) {
	service := &model.Service{
		Name: name,
	}
	services, err := service.GetByName()
	if err != nil {
		return []model.Service{}, err
	}
	return services, nil
}

func FetchByVersionAndName(name string, version int) (model.Service, error) {
	service := &model.Service{
		Name:    name,
		Version: version,
	}
	serviceFetched, err := service.GetByNameAndVersion()
	if err != nil {
		return model.Service{}, err
	}
	return serviceFetched, nil
}

func FetchAll() ([]model.Service, error) {
	service := &model.Service{}
	servicesFetched, err := service.List()
	if err != nil {
		return []model.Service{}, err
	}
	return servicesFetched, nil
}

func SearchAndSort(query, sort, dir string, page, pageSize int) (apiv1.ServicePagination, error) {
	service := model.Service{}
	var serviceDetailsHolder []apiv1.Service
	serviceData, totalCount, err := service.GetServiceAndVersionCounts(query, pageSize, page, sort, dir)
	if err != nil {
		return apiv1.ServicePagination{}, err
	}
	for _, s := range serviceData {
		details, err := serviceDetails(s.Name)
		if err != nil {
			return apiv1.ServicePagination{}, err
		}
		serviceDetailsHolder = append(serviceDetailsHolder, details)
	}
	response := apiv1.ServicePagination{
		Meta: apiv1.Meta{
			Page:         page,
			PageSize:     pageSize,
			TotalResults: int(totalCount),
			TotalPages:   int(math.Ceil(float64(totalCount) / float64(pageSize))),
		},
		Data: serviceDetailsHolder,
	}
	return response, err
}

func serviceDetails(name string) (apiv1.Service, error) {
	var response apiv1.Service
	Versions, err := FetchByName(name)
	if err != nil {
		return apiv1.Service{}, err
	}
	service, err := FetchByVersionAndName(name, len(Versions))
	if err != nil {
		return apiv1.Service{}, nil
	}
	response = apiv1.Service{
		TotalVersions:  len(Versions),
		CurrentVersion: len(Versions),
		Service:        &service,
	}
	return response, nil
}
