/*
 * Created on 15/09/23 02.30
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package participant

import (
	"bumn-sembako-be/helper"
	"bumn-sembako-be/model"
	"bumn-sembako-be/request"
	"fmt"
	"gorm.io/gorm"
)

type Service interface {
	ReadAllBy(criteria map[string]interface{}, search string, page, size int) (*[]model.Participant, error)
	ReadById(id int) (*model.Participant, error)
	Count(criteria map[string]interface{}) int64
	Update(id int, participant request.UpdateParticipant) (*model.Participant, error)
	UpdateStatus(id int, status *request.PartialDone) (*model.Participant, error)
	Create(participant *model.Participant) (*model.Participant, error)
	ReadAllReport(criteria map[string]interface{}) (*[]model.Report, error)
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

func (e *service) Update(id int, participant request.UpdateParticipant) (*model.Participant, error) {
	var upParticipant = model.Participant{}
	err := e.db.Table("participants").Where("id = ?", id).First(&upParticipant).Updates(&participant).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[participant.service.Update] error execute query %v \n", err)
		return nil, fmt.Errorf("failed update data")
	}
	return &upParticipant, nil
}

func (e *service) UpdateStatus(id int, status *request.PartialDone) (*model.Participant, error) {
	var upParticipant = model.Participant{}
	err := e.db.Table("participants").Where("id = ?", id).First(&upParticipant).Updates(&status).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[participant.service.UpdateStatus] error execute query %v \n", err)
		return nil, fmt.Errorf("failed update data")
	}
	return &upParticipant, nil
}

func (e *service) Create(participant *model.Participant) (*model.Participant, error) {
	tx := e.db.Begin()
	defer tx.Rollback()

	err := tx.Save(&participant).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[participant.service.Create] error execute query %v \n", err)
		return nil, fmt.Errorf("failed insert data")
	}

	helper.CommonLogger().Error(err)

	tx.Commit()

	return participant, nil
}

func (s *service) ReadAllReport(criteria map[string]interface{}) (*[]model.Report, error) {
	var reports []model.Report

	query := s.db.Table("participants").Where(criteria)
	err := query.Order("updated_at ASC").Find(&reports).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[participant.service.ReadAllReport] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return &reports, nil
}
