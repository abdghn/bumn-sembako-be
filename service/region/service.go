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
	ReadAllProvinceBy(criteria map[string]interface{}, search string) ([]*model.Province, error)
	ReadProvinceBy(criteria map[string]interface{}) (*model.Province, error)
	ReadAllRegencyBy(criteria map[string]interface{}, search string) ([]*model.Regency, error)
	ReadRegencyBy(criteria map[string]interface{}) (*model.Regency, error)
	ReadAllDistrictBy(criteria map[string]interface{}, search string) ([]*model.District, error)
	ReadDistrictBy(criteria map[string]interface{}) (*model.District, error)
	ReadAllVillageBy(criteria map[string]interface{}, search string) ([]*model.Village, error)
	ReadVillageBy(criteria map[string]interface{}) (*model.Village, error)
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}

}

func (s *service) ReadAllProvinceBy(criteria map[string]interface{}, search string) ([]*model.Province, error) {
	var provincies []*model.Province

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
	return provincies, nil
}

func (s *service) ReadProvinceBy(criteria map[string]interface{}) (*model.Province, error) {
	var province = model.Province{}
	err := s.db.Table("provinces").Where(criteria).First(&province).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[region.service.ReadProvinceBy] error execute query %v \n", err)
		return nil, fmt.Errorf("province is not exists")
	}
	return &province, nil
}

func (s *service) ReadAllRegencyBy(criteria map[string]interface{}, search string) ([]*model.Regency, error) {
	var regencies []*model.Regency

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
	return regencies, nil
}

func (s *service) ReadRegencyBy(criteria map[string]interface{}) (*model.Regency, error) {
	var regency = model.Regency{}
	err := s.db.Table("regencies").Where(criteria).First(&regency).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[region.service.ReadRegencyBy] error execute query %v \n", err)
		return nil, fmt.Errorf("city is not exists")
	}
	return &regency, nil
}

func (s *service) ReadAllDistrictBy(criteria map[string]interface{}, search string) ([]*model.District, error) {
	var districts []*model.District

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
	return districts, nil
}

func (s *service) ReadDistrictBy(criteria map[string]interface{}) (*model.District, error) {
	var district = model.District{}
	err := s.db.Table("districts").Where(criteria).First(&district).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[region.service.ReadDistrictBy] error execute query %v \n", err)
		return nil, fmt.Errorf("district is not exists")
	}
	return &district, nil
}

func (s *service) ReadAllVillageBy(criteria map[string]interface{}, search string) ([]*model.Village, error) {
	var villages []*model.Village

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
	return villages, nil
}

func (s *service) ReadVillageBy(criteria map[string]interface{}) (*model.Village, error) {
	var village = model.Village{}
	err := s.db.Table("villages").Where(criteria).First(&village).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[region.service.ReadVillageBy] error execute query %v \n", err)
		return nil, fmt.Errorf("village is not exists")
	}
	return &village, nil
}
