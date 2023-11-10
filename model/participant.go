/*
 * Created on 15/09/23 02.10
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package model

import (
	"html/template"
	"time"
)

type Participant struct {
	ID                 int        `json:"id" gorm:"primary_key"`
	Name               string     `json:"name"`
	NIK                string     `json:"nik"`
	Gender             string     `json:"gender"`
	Phone              string     `json:"phone"`
	Address            string     `json:"address"`
	RT                 string     `json:"rt"`
	RW                 string     `json:"rw"`
	Provinsi           string     `json:"provinsi"`
	Kota               string     `json:"kota"`
	Kecamatan          string     `json:"kecamatan"`
	Kelurahan          string     `json:"kelurahan"`
	KodePOS            string     `json:"kode_pos"`
	ResidenceAddress   string     `json:"residence_address"`
	ResidenceRT        string     `json:"residence_rt"`
	ResidenceRW        string     `json:"residence_rw"`
	ResidenceProvinsi  string     `json:"residence_provinsi"`
	ResidenceKota      string     `json:"residence_kota"`
	ResidenceKecamatan string     `json:"residence_kecamatan"`
	ResidenceKelurahan string     `json:"residence_kelurahan"`
	ResidenceKodePOS   string     `json:"residence_kode_pos"`
	Status             string     `json:"status"`
	Image              string     `json:"image" gorm:"type: text"`
	ImagePenerima      string     `json:"image_penerima" gorm:"type: text"`
	IsRepresented      bool       `json:"is_represented" gorm:"default:false"`
	UpdatedBy          string     `json:"updated_by"`
	Reference          string     `json:"reference" gorm:"type:varchar(255)"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `sql:"index" json:"deleted_at"`
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
	Provinsi           string `json:"provinsi"`
	Kota               string `json:"kota"`
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
