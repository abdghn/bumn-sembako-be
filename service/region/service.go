/*
 * Created on 15/09/23 16.40
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package region

import (
	"bumn-sembako-be/helper"
	"bumn-sembako-be/model"
	"fmt"
	"gorm.io/gorm"
)

type Service interface {
	ReadAllProvinceBy(criteria map[string]interface{}, search string) (*[]model.Province, error)
	ReadAllRegencyBy(criteria map[string]interface{}, search string) (*[]model.Regency, error)
	ReadAllDistrictBy(criteria map[string]interface{}, search string) (*[]model.District, error)
	ReadAllVillageBy(criteria map[string]interface{}, search string) (*[]model.Village, error)
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}

}

func (s *service) ReadAllProvinceBy(criteria map[string]interface{}, search string) (*[]model.Province, error) {
	var provincies []model.Province

	query := s.db.Where(criteria)

	if search != "" {
		query.Where("name LIKE ?", search+"%")
	}

	err := query.Order("created_at ASC").Find(&provincies).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[region.service.ReadAllProvinceBy] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return &provincies, nil
}

func (s *service) ReadAllRegencyBy(criteria map[string]interface{}, search string) (*[]model.Regency, error) {
	var regencies []model.Regency

	query := s.db.Where(criteria)

	if search != "" {
		query.Where("name LIKE ?", search+"%")
	}

	err := query.Order("created_at ASC").Find(&regencies).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[region.service.ReadAllRegencyBy] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return &regencies, nil
}

func (s *service) ReadAllDistrictBy(criteria map[string]interface{}, search string) (*[]model.District, error) {
	var districts []model.District

	query := s.db.Where(criteria)

	if search != "" {
		query.Where("name LIKE ?", search+"%")
	}

	err := query.Order("created_at ASC").Find(&districts).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[region.service.ReadAllDistrictBy] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return &districts, nil
}

func (s *service) ReadAllVillageBy(criteria map[string]interface{}, search string) (*[]model.Village, error) {
	var villages []model.Village

	query := s.db.Where(criteria)

	if search != "" {
		query.Where("name LIKE ?", search+"%")
	}

	err := query.Order("created_at ASC").Find(&villages).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[region.service.ReadAllVillageBy] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return &villages, nil
}
