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
	"time"

	"gorm.io/gorm"
)

type Service interface {
	ReadAllBy(criteria map[string]interface{}, search string, page, size int) (*[]model.Participant, error)
	ReadAllLogBy(criteria map[string]interface{}, search string, page, size int) ([]model.ImportLog, error)
	ReadById(id int) (*model.Participant, error)
	Count(criteria map[string]interface{}, search string) int64
	CountLogs(criteria map[string]interface{}) int64
	CountByDate(criteria map[string]interface{}, date time.Time) int64
	CountByRangeDate(criteria map[string]interface{}, startDate, endDate time.Time) int64
	Update(id int, participant request.UpdateParticipant) (*model.Participant, error)
	UpdateStatus(id int, status *request.PartialDone) (*model.Participant, error)
	Create(participant *model.Participant) (*model.Participant, error)
	ReadAllReport(criteria map[string]interface{}, date time.Time) ([]*model.Report, error)
	ReadAllReportByRangeDate(criteria map[string]interface{}, startDate, endDate time.Time) ([]*model.Report, error)
	GetQuota(criteria map[string]interface{}) (*model.Quota, error)
	CreateLog(m *model.ImportLog) (*model.ImportLog, error)
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
		query.Where("name LIKE ?", search+"%").
			Or("nik LIKE ?", search+"%").
			Or("phone LIKE ?", search+"%").
			Or("address LIKE ?", search+"%").
			Or("provinsi LIKE ?", search+"%").
			Or("kota LIKE ?", search+"%").
			Or("kecamatan LIKE ?", search+"%").
			Or("kelurahan LIKE ?", search+"%").
			Or("kode_pos LIKE ?", search+"%").
			Or("residence_address LIKE ?", search+"%").
			Or("residence_provinsi LIKE ?", search+"%").
			Or("residence_kota LIKE ?", search+"%").
			Or("residence_kecamatan LIKE ?", search+"%").
			Or("residence_kelurahan LIKE ?", search+"%").
			Or("residence_kode_pos LIKE ?", search+"%")
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

func (e *service) ReadAllLogBy(criteria map[string]interface{}, search string, page, size int) ([]model.ImportLog, error) {
	var logs []model.ImportLog

	query := e.db.Where(criteria)

	if search != "" {
		query.Where("name LIKE ?", search+"%")
	}

	limit, offset := helper.GetLimitOffset(page, size)
	err := query.Offset(offset).Order("created_at DESC").Limit(limit).Find(&logs).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[participant.service.ReadAllLogBy] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return logs, nil
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

func (s *service) Count(criteria map[string]interface{}, search string) int64 {
	var result int64
	query := s.db.Table("participants").Where(criteria)
	if search != "" {
		query.Where("name LIKE ?", search+"%").
			Or("nik LIKE ?", search+"%").
			Or("phone LIKE ?", search+"%").
			Or("address LIKE ?", search+"%").
			Or("provinsi LIKE ?", search+"%").
			Or("kota LIKE ?", search+"%").
			Or("kecamatan LIKE ?", search+"%").
			Or("kelurahan LIKE ?", search+"%").
			Or("kode_pos LIKE ?", search+"%").
			Or("residence_address LIKE ?", search+"%").
			Or("residence_provinsi LIKE ?", search+"%").
			Or("residence_kota LIKE ?", search+"%").
			Or("residence_kecamatan LIKE ?", search+"%").
			Or("residence_kelurahan LIKE ?", search+"%").
			Or("residence_kode_pos LIKE ?", search+"%")
	}

	err := query.Count(&result).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		return 0
	}
	return result
}

func (s *service) CountLogs(criteria map[string]interface{}) int64 {
	var result int64
	err := s.db.Table("import_logs").Where(criteria).Count(&result).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		return 0
	}
	return result
}

func (s *service) CountByDate(criteria map[string]interface{}, date time.Time) int64 {
	var result int64
	query := s.db.Table("participants").Where(criteria)
	if !date.IsZero() {
		query.Where("updated_at < ?", date)

	}
	err := query.Count(&result).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		return 0
	}
	return result
}

func (s *service) CountByRangeDate(criteria map[string]interface{}, startDate, endDate time.Time) int64 {
	var result int64
	query := s.db.Table("participants").Where(criteria)
	if !startDate.IsZero() && !endDate.IsZero() {
		query.Where("updated_at <= ? AND updated_at >=  ?", endDate, startDate)

	}
	err := query.Count(&result).Error
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

	tx.Commit()

	return participant, nil
}

func (s *service) ReadAllReport(criteria map[string]interface{}, date time.Time) ([]*model.Report, error) {
	var reports []*model.Report

	query := s.db.Table("participants").Select("ROW_NUMBER() OVER (ORDER BY id) AS No", "nik as NIK", "name AS Name", "image_penerima as Image", "phone AS Phone", "address AS Address", "COUNT(*) OVER() AS Total").Where(criteria)
	if !date.IsZero() {
		query.Where("updated_at < ?", date)

	}
	err := query.Order("updated_at ASC").Find(&reports).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[participant.service.ReadAllReport] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return reports, nil
}

func (s *service) ReadAllReportByRangeDate(criteria map[string]interface{}, startDate, endDate time.Time) ([]*model.Report, error) {
	var reports []*model.Report

	query := s.db.Table("participants").Select("ROW_NUMBER() OVER (ORDER BY id) AS No", "nik as NIK", "name AS Name", "image_penerima as Image", "phone AS Phone", "address AS Address", "COUNT(*) OVER() AS Total").Where(criteria)
	if !startDate.IsZero() && !endDate.IsZero() {
		query.Where("updated_at <= ? AND updated_at >=  ?", endDate, startDate)

	}
	err := query.Order("updated_at ASC").Find(&reports).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[participant.service.ReadAllReportByRangeDate] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return reports, nil
}

func (s *service) GetQuota(criteria map[string]interface{}) (*model.Quota, error) {

	var quota = model.Quota{}
	err := s.db.Table("quotas").Where(criteria).Find(&quota).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[participant.service.GetQuota] error execute query %v \n", err)
		return nil, err
	}
	return &quota, nil
}

func (s *service) CreateLog(m *model.ImportLog) (*model.ImportLog, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	err := tx.Save(&m).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[user.service.CreateLog] error execute query %v \n", err)
		return nil, fmt.Errorf("failed insert data")
	}

	tx.Commit()

	return m, nil
}
