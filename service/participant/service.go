/*
 * Created on 15/09/23 02.30
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package participant

import (
	"bumn-sembako-be/helper"
	"bumn-sembako-be/model"
	"fmt"
	"gorm.io/gorm"
)

type Service interface {
	ReadAllBy(criteria map[string]interface{}, search string, page, size int) (*[]model.Participant, error)
	ReadById(id int) (*model.Participant, error)
	Count(criteria map[string]interface{}) int64
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}

}

func (e *service) ReadAllBy(criteria map[string]interface{}, search string, page, size int) (*[]model.Participant, error) {
	var participants []model.Participant

	query := e.db.Where(criteria)

	if search != "" {
		query.Where("name LIKE ?", search+"%")
	}

	if page == 0 || size == 0 {
		page, size = -1, -1
	}

	limit, offset := helper.GetLimitOffset(page, size)
	err := query.Offset(offset).Order("created_at ASC").Limit(limit).Find(&participants).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[participant.service.ReadAllBy] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return &participants, nil
}

func (s *service) ReadById(id int) (*model.Participant, error) {
	var participant = model.Participant{}
	err := s.db.Table("participants").Where("id = ?", id).First(&participant).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[participant.service.ReadById] error execute query %v \n", err)
		return nil, fmt.Errorf("id is not exists")
	}
	return &participant, nil
}

func (s *service) Count(criteria map[string]interface{}) int64 {
	var result int64
	err := s.db.Table("participants").Where(criteria).Count(&result).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		return 0
	}
	return result
}
