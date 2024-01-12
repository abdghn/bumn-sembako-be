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
	ReadBy(criteria map[string]interface{}) ([]*model.Participant, error)
	Count(criteria map[string]interface{}, search string) int64
	CountNotInId(criteria map[string]interface{}, id int) int64
	CountLogs(criteria map[string]interface{}) int64
	CountByDate(criteria map[string]interface{}, date time.Time) int64
	CountByRangeDate(criteria map[string]interface{}, startDate, endDate time.Time) int64
	Update(id int, participant *request.ParticipantEditInput) (*model.Participant, error)
	UpdateStatus(id int, status *request.PartialDone) (*model.Participant, error)
	Create(participant *model.Participant) (*model.Participant, error)
	ReadAllReport(criteria map[string]interface{}, date time.Time) ([]*model.Report, error)
	ReadAllReportByRangeDate(criteria map[string]interface{}, startDate, endDate time.Time) ([]*model.Report, error)
	ReadAllReportByRangeDateV2(criteria map[string]interface{}, startDate, endDate time.Time, page, size int) ([]*model.Report, error)
	ReadAllReportByRangeDateV3(criteria map[string]interface{}, startDate, endDate time.Time, page, size int) ([]*model.Report, error)
	GetQuota(criteria map[string]interface{}) (*model.Quota, error)
	CreateLog(m *model.ImportLog) (*model.ImportLog, error)
	CountAllStatus(criteria map[string]interface{}) (*model.TotalParticipantResponse, error)
	CountAllStatusGroup(criteria map[string]interface{}) ([]*model.TotalParticipantListResponse, error)
	Reset(id int) (*model.Participant, error)
	Delete(id int) error
	DeleteBy(criteria map[string]interface{}) error
	ReadAllDuplicates() ([]*model.Participant, error)
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}

}

func (e *service) ReadAllBy(criteria map[string]interface{}, search string, page, size int) (*[]model.Participant, error) {
	var participants []model.Participant

	query := e.db

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
	err := query.Where(criteria).Where("deleted_at IS NULL").Offset(offset).Order("created_at ASC").Limit(limit).Find(&participants).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[participant.service.ReadAllBy] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return &participants, nil
}

func (e *service) ReadAllLogBy(criteria map[string]interface{}, search string, page, size int) ([]model.ImportLog, error) {
	var logs []model.ImportLog

	query := e.db

	if search != "" {
		query.Where("name LIKE ?", search+"%")
	}

	limit, offset := helper.GetLimitOffset(page, size)
	err := query.Where(criteria).Where("deleted_at IS NULL").Offset(offset).Order("created_at DESC").Limit(limit).Find(&logs).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[participant.service.ReadAllLogBy] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return logs, nil
}

func (s *service) ReadById(id int) (*model.Participant, error) {
	var participant = model.Participant{}
	err := s.db.Table("participants").Where("id = ?", id).Where("deleted_at IS NULL").First(&participant).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[participant.service.ReadById] error execute query %v \n", err)
		return nil, fmt.Errorf("id is not exists")
	}
	return &participant, nil
}

func (s *service) ReadBy(criteria map[string]interface{}) ([]*model.Participant, error) {
	var participants []*model.Participant
	err := s.db.Table("participants").Where(criteria).Where("deleted_at IS NULL").First(&participants).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[participant.service.ReadBy] error execute query %v \n", err)
		return nil, fmt.Errorf("id is not exists")
	}
	return participants, nil
}

func (s *service) Count(criteria map[string]interface{}, search string) int64 {
	var result int64
	query := s.db.Table("participants")
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

	err := query.Where(criteria).Where("deleted_at IS NULL").Count(&result).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		return 0
	}
	return result
}

func (s *service) CountNotInId(criteria map[string]interface{}, id int) int64 {
	var result int64
	query := s.db.Table("participants").Where(criteria).Where("id NOT IN (?)", id).Where("deleted_at IS NULL")
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
	query := s.db.Table("participants").Where(criteria).Where("deleted_at IS NULL")
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

func (e *service) Update(id int, participant *request.ParticipantEditInput) (*model.Participant, error) {
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

	query := s.db.Table("participants").Select("ROW_NUMBER() OVER (ORDER BY id) AS No", "nik as NIK", "name AS Name", "image_penerima as Image", "phone AS Phone", "address AS Address", "COUNT(*) OVER() AS Total").Where("deleted_at IS NULL").Where(criteria)
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

	query := s.db.Table("participants").Select("ROW_NUMBER() OVER (ORDER BY id) AS No", "nik as NIK", "name AS Name", "SUBSTRING(image_penerima, 7) as Image", "phone AS Phone", "address AS Address", "COUNT(*) OVER() AS Total").Where("deleted_at IS NULL").Where(criteria)
	if !startDate.IsZero() && !endDate.IsZero() {
		query.Where("updated_at <= ? AND updated_at >=  ?", endDate, startDate)

	}

	err := query.Order("updated_at ASC").Find(&reports).Updates(map[string]interface{}{
		"has_printed": true,
	}).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[participant.service.ReadAllReportByRangeDate] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return reports, nil
}

func (s *service) ReadAllReportByRangeDateV2(criteria map[string]interface{}, startDate, endDate time.Time, page, size int) ([]*model.Report, error) {
	var reports []*model.Report

	query := s.db.Table("participants").Select("ROW_NUMBER() OVER (ORDER BY id) AS No", "id AS ID", "nik as NIK", "name AS Name", "SUBSTRING(image_penerima, 7) as Image", "phone AS Phone", "address AS Address", "COUNT(*) OVER() AS Total").Where("deleted_at IS NULL").Where(criteria)
	if !startDate.IsZero() && !endDate.IsZero() {
		query.Where("updated_at <= ? AND updated_at >=  ?", endDate, startDate)

	}

	query.Where("image_penerima IS NOT NULL")

	limit, offset := helper.GetLimitOffset(page, size)
	err := query.Offset(offset).Limit(limit).Order("updated_at ASC").Find(&reports).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[participant.service.ReadAllReportByRangeDate] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}
	return reports, nil
}

func (s *service) ReadAllReportByRangeDateV3(criteria map[string]interface{}, startDate, endDate time.Time, page, size int) ([]*model.Report, error) {
	var ids []int
	var reports []*model.Report

	s.db.Table("participants").Select("id").Where("deleted_at IS NULL").Find(&ids)

	query := s.db.Table("participants").Select("ROW_NUMBER() OVER (ORDER BY id) AS No", "nik as NIK", "name AS Name", "SUBSTRING(image_penerima, 7) as Image", "phone AS Phone", "address AS Address", "COUNT(*) OVER() AS Total").Where("deleted_at IS NULL").Where(criteria)
	query.Where("id IN (?)", ids)
	if !startDate.IsZero() && !endDate.IsZero() {
		query.Where("updated_at <= ? AND updated_at >=  ?", endDate, startDate)

	}

	query.Where("image_penerima IS NOT NULL")

	limit, offset := helper.GetLimitOffset(page, size)
	err := query.Offset(offset).Limit(limit).Order("updated_at ASC").Updates(model.Participant{HasPrinted: true}).Find(&reports).Error
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

func (s *service) CountAllStatus(criteria map[string]interface{}) (*model.TotalParticipantResponse, error) {
	var totalData = model.TotalParticipantResponse{}

	query := `COUNT(*) AS total_penerima,
				SUM( status = "DONE" ) AS total_sudah_menerima,
				SUM( status = "PARTIAL_DONE" ) AS total_partial_done,
				SUM( status = "NOT DONE" ) AS total_belum_menerima,
				SUM( status = "REJECTED" ) AS total_data_gugur,
				0 AS total_quota`

	err := s.db.Table("participants").Select(query).Where("deleted_at IS NULL").Where(criteria).Find(&totalData).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[participant.service.CountAllStatus] error execute query %v \n", err)
		return nil, fmt.Errorf("failed count all status")
	}

	return &totalData, nil

}

func (s *service) CountAllStatusGroup(criteria map[string]interface{}) ([]*model.TotalParticipantListResponse, error) {
	var list []*model.TotalParticipantListResponse

	query := `residence_provinsi, residence_kota,
				COUNT(*) AS total_penerima,
				SUM( status = "DONE" ) AS total_sudah_menerima,
				SUM( status = "PARTIAL_DONE" ) AS total_partial_done,
				SUM( status = "NOT DONE" ) AS total_belum_menerima,
				SUM( status = "REJECTED" ) AS total_data_gugur`

	err := s.db.Table("participants").Select(query).Where("deleted_at IS NULL").Where(criteria).Group("residence_provinsi, residence_kota").Order("residence_provinsi ASC").Find(&list).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[participant.service.CountAllStatusGroup] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}

	return list, nil

}

func (s *service) Reset(id int) (*model.Participant, error) {
	var upParticipant = model.Participant{}
	err := s.db.Table("participants").Where("id = ?", id).First(&upParticipant).Updates(map[string]interface{}{
		"status": "NOT DONE",
	}).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[participant.service.Reset] error execute query %v \n", err)
		return nil, fmt.Errorf("failed update data")
	}
	return &upParticipant, nil
}

func (s *service) Delete(id int) error {
	var participant = model.Participant{}
	err := s.db.Table("participants").Where("id = ?", id).First(&participant).Delete(&participant).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[participant.service.Delete] error execute query %v \n", err)
		return fmt.Errorf("id is not exists")
	}
	return nil
}

func (s *service) DeleteBy(criteria map[string]interface{}) error {
	var participant = model.Participant{}
	err := s.db.Table("participants").Where(criteria).First(&participant).Delete(&participant).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[participant.service.Delete] error execute query %v \n", err)
		return fmt.Errorf("id is not exists")
	}
	return nil
}

func (s *service) ReadAllDuplicates() ([]*model.Participant, error) {
	var participants []*model.Participant
	query := `t.* 
				right join
				(select nik from participants group by 1 having count(*)>1) duplicates
				on duplicates.nik=t.nik
				order by nik`
	err := s.db.Table("participants t").Select(query).Where("deleted_at IS NULL").Find(&participants).Error
	if err != nil {
		helper.CommonLogger().Error(err)
		fmt.Printf("[participant.service.ReadAllDuplicates] error execute query %v \n", err)
		return nil, fmt.Errorf("failed view all data")
	}

	return participants, nil
}
