/*
 * Created on 15/09/23 02.10
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package model

import "time"

type Participant struct {
	ID                 int        `json:"id" gorm:"primary_key"`
	Name               string     `json:"name"`
	NIK                string     `json:"nik"  gorm:"unique"`
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
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `sql:"index" json:"deleted_at"`
}