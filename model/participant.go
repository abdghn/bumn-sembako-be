/*
 * Created on 15/09/23 02.10
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package model

import (
	"html/template"
	"time"

	"gorm.io/gorm"
)

type Participant struct {
	ID                 int            `json:"id" gorm:"primary_key"`
	Name               string         `json:"name" gorm:"type:varchar(255)"`
	NIK                string         `json:"nik" gorm:"type:varchar(255)"`
	Gender             string         `json:"gender" gorm:"type:varchar(255)"`
	Phone              string         `json:"phone" gorm:"type:varchar(255)"`
	Address            string         `json:"address" gorm:"type: text"`
	RT                 string         `json:"rt" gorm:"type:varchar(255)"`
	RW                 string         `json:"rw" gorm:"type:varchar(255)"`
	Provinsi           string         `sql:"index" json:"provinsi" gorm:"type:varchar(255);index"`
	Kota               string         `sql:"index" json:"kota" gorm:"type:varchar(255);index"`
	Kecamatan          string         `sql:"index" json:"kecamatan" gorm:"type:varchar(255);index"`
	Kelurahan          string         `sql:"index" json:"kelurahan" gorm:"type:varchar(255);index"`
	KodePOS            string         `json:"kode_pos" gorm:"type:varchar(255)"`
	ResidenceAddress   string         `json:"residence_address" gorm:"type: text"`
	ResidenceRT        string         `json:"residence_rt" gorm:"type:varchar(255)"`
	ResidenceRW        string         `json:"residence_rw" gorm:"type:varchar(255)"`
	ResidenceProvinsi  string         `sql:"index" json:"residence_provinsi" gorm:"type:varchar(255);index"`
	ResidenceKota      string         `sql:"index" json:"residence_kota" gorm:"type:varchar(255);index"`
	ResidenceKecamatan string         `json:"residence_kecamatan" gorm:"type:varchar(255)"`
	ResidenceKelurahan string         `json:"residence_kelurahan" gorm:"type:varchar(255)"`
	ResidenceKodePOS   string         `json:"residence_kode_pos" gorm:"type:varchar(255)"`
	Status             string         `json:"status" gorm:"type:varchar(255)"`
	Image              string         `json:"image" gorm:"type: text"`
	ImagePenerima      string         `json:"image_penerima" gorm:"type: text"`
	IsRepresented      bool           `json:"is_represented" gorm:"default:false"`
	HasPrinted         bool           `json:"has_represented" gorm:"default:false"`
	UpdatedBy          string         `json:"updated_by" gorm:"type:varchar(100)"`
	Reference          string         `json:"reference" gorm:"type: text"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `sql:"index" json:"deleted_at" gorm:"index"`
}

type TotalParticipantResponse struct {
	TotalPenerima      int64 `json:"total_penerima"`
	TotalSudahMenerima int64 `json:"total_sudah_menerima"`
	TotalPartialDone   int64 `json:"total_partial_done"`
	TotalBelumMenerima int64 `json:"total_belum_menerima"`
	TotalDataGugur     int64 `json:"total_data_gugur"`
	TotalQuota         int64 `json:"total_quota"`
}

type TotalParticipantListResponse struct {
	ResidenceProvinsi  string `json:"residence_provinsi"`
	ResidenceKota      string `json:"residence_kota"`
	TotalPenerima      int64  `json:"total_penerima"`
	TotalSudahMenerima int64  `json:"total_sudah_menerima"`
	TotalPartialDone   int64  `json:"total_partial_done"`
	TotalBelumMenerima int64  `json:"total_belum_menerima"`
	TotalDataGugur     int64  `json:"total_data_gugur"`
}

type Report struct {
	No       int
	NIK      string
	Name     string
	Phone    string
	Image    string
	Address  string
	ImageB64 template.URL
	Total    int
}

type ReportPerFile struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
