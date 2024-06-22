package model

import (
	"errors"
	"fmt"
	"github.com/suyog1pathak/services/pkg/datastore"
	customerrors "github.com/suyog1pathak/services/pkg/errors/service"
	log "github.com/suyog1pathak/services/pkg/logger"
	"gorm.io/gorm"
	"sync"
)

var db *gorm.DB
var err error
var once sync.Once

func Setup() {
	once.Do(func() {
		db, err = datastore.GetDBConnection()
		return
	})
}

type ServiceCount struct {
	Name  string
	Count int64
}

type Service struct {
	gorm.Model  `swaggerignore:"true"`
	Name        string `json:"serviceName"`
	Description string `json:"describe"`
	Version     int    `json:"version" swaggerignore:"true"`
	IsActive    bool   `json:"isActive" swaggertype:"boolean"`
	Tags        string `json:"tags" `
} //@name ServiceModelDb

// only if you want to use table with custom name
//type Tabler interface {
//	TableName() string
//}
//
//func (Service) TableName() string {
//	return "service"
//}

//-----------------------------//

func (s *Service) Add() error {
	log.Debug("adding service", "service", s.Name)
	result := db.Create(s)
	if result.Error != nil {
		log.Error("error in adding services", "service", s.Name, "error", result.Error.Error())
		return result.Error
	}
	return nil
}

func (s *Service) List() ([]Service, error) {
	log.Debug("fetching all services")
	var output []Service
	result := db.Where("is_active = ?", true).Find(&output)
	if result.Error != nil {
		log.Error("error in listing services", "error", result.Error.Error())
		return output, result.Error
	}
	return output, nil
}

func (s *Service) GetByName() ([]Service, error) {
	log.Debug("fetching service by name", "service", s.Name)
	var output []Service
	result := db.Where("name = ?", s.Name).Find(&output)
	if result.Error != nil {
		log.Error("error in getting service by name", "name", s.Name, "error", result.Error.Error())
		return output, result.Error
	}
	if result.RowsAffected == 0 {
		log.Warn("service not found", "name", s.Name)
		return output, errors.New(customerrors.ErrServiceNotFound)
	}
	return output, nil
}

func (s *Service) GetByNameCount() (int64, error) {
	log.Debug("fetching count service by name", "service", s.Name)
	var output int64
	result := db.Model(&Service{}).Debug().Where("name = ?", s.Name).Count(&output)
	if result.Error != nil {
		log.Error("error in fetching count service by name", "name", s.Name, "error", result.Error.Error())
		return output, result.Error
	}
	if result.RowsAffected == 0 {
		log.Warn("service not found", "name", s.Name)
		return output, errors.New(customerrors.ErrServiceNotFound)
	}
	return output, nil
}

func (s *Service) GetServiceAndVersionCounts(query string, limit int, offset int, sort, dir string) ([]ServiceCount, int64, error) {
	log.Debug("fetching service and version counts", "query", query, "limit", limit,
		"offset", offset, "sort_by", sort, "direction", dir)
	var count int64
	/*
		SELECT COUNT(count) as totalCount
			FROM (
				SELECT name, COUNT(*) as count
				FROM `services` WHERE 'name' Like '%' AND `services`.`deleted_at` IS NULL
				GROUP BY `name`
			) as u
	*/
	totalCount := db.Table("(?) as u", db.Debug().Model(&Service{}).
		Select("name, COUNT(*) as count").
		Where("'name' Like ?", query).Group("name")).Select("COUNT(count) as totalCount").Scan(&count)

	if totalCount.Error != nil {
		log.Error("error in fetching service and version counts", "error", totalCount.Error.Error())
		return []ServiceCount{}, 0, totalCount.Error
	}

	var results []ServiceCount
	o := (offset - 1) * limit //from where we want to start
	result := db.Limit(limit).Offset(o).Model(&Service{}).
		Select("name, COUNT(*) as count"+fmt.Sprintf(", MAX(%s) as %s", sort, sort)).
		Where("'name' Like ?", query).
		Group("name").
		Order(fmt.Sprintf("%s %s", sort, dir)).
		Scan(&results)

	if result.Error != nil {
		log.Error("error in fetching service with pagination", "error", result.Error.Error())
		return []ServiceCount{}, 0, result.Error
	}

	if result.RowsAffected == 0 {
		return []ServiceCount{}, 0, errors.New(customerrors.ErrServiceNotFound)
	}
	return results, count, nil
}

func (s *Service) GetByNameAndVersion() (Service, error) {
	log.Debug("fetching service with name and version", "name", s.Name, "version", s.Version)
	var output Service
	result := db.Where("name = ? and version = ?", s.Name, s.Version).Find(&output)
	if result.Error != nil {
		log.Error("error in fetching service with name and version", "name", s.Name, "version", s.Version, "error", result.Error)
		return output, result.Error
	}
	if result.RowsAffected == 0 {
		log.Warn("service not found", "name", s.Name, "version", s.Version)
		return output, errors.New(customerrors.ErrServiceWithVersionNotFound)
	}
	return output, nil
}

func (s *Service) UpdateByNameAndVersion() error {
	log.Debug("updating service with name and version", "name", s.Name, "version", s.Version)
	result := db.Model(&s).Where("name = ? and version = ?", s.Name, s.Version).Updates(s)
	if result.Error != nil {
		log.Error("error in updating service with name and version", "name", s.Name, "version", s.Version, "error", result.Error.Error())
		return result.Error
	}
	return nil
}

func (s *Service) DeleteByName() error {
	log.Debug("deleting service", "name", s.Name)
	//add Unscoped() for hard delete
	result := db.Delete(s, "name = ?", s.Name)
	if result.Error != nil {
		log.Error("error in deleting service", "name", s.Name, "error", result.Error.Error())
		return result.Error
	}
	return nil
}
