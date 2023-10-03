/*
 * Created on 15/09/23 14.32
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package request

import (
	"html/template"
	"mime/multipart"
	"time"
)

type ParticipantPaged struct {
	Search    string `form:"search"`
	Page      int    `form:"page"`
	Size      int    `form:"size"`
	Provinsi  string `form:"provinsi"`
	Kota      string `form:"kota"`
	Kecamatan string `form:"kecamatan"`
	Kelurahan string `form:"kelurahan"`
	Status    string `form:"status"`
}

type ParticipantFilter struct {
	Provinsi string    `form:"provinsi"`
	Kota     string    `form:"kota"`
	Date     time.Time `form:"date" time_format:"2006-01-02"`
}

type UpdateParticipant struct {
	Name               string                `json:"name" form:"name"`
	NIK                string                `json:"nik" form:"nik" `
	Gender             string                `json:"gender" form:"gender"`
	Phone              string                `json:"phone" form:"phone"`
	Address            string                `json:"address" form:"address"`
	RT                 string                `json:"rt" form:"rt"`
	RW                 string                `json:"rw" form: "rw"`
	Provinsi           string                `json:"provinsi" form:"provinsi"`
	Kota               string                `json:"kota" form:"kota"`
	Kecamatan          string                `json:"kecamatan" form:"kecamatan"`
	Kelurahan          string                `json:"kelurahan" form: "kelurahan"`
	KodePOS            string                `json:"kode_pos" form:"kode_pos"`
	ResidenceAddress   string                `json:"residence_address" form:"residence_address"`
	ResidenceRT        string                `json:"residence_rt" form:"residence_rt"`
	ResidenceRW        string                `json:"residence_rw" form:"residence_rw"`
	ResidenceProvinsi  string                `json:"residence_provinsi" form:"residence_provinsi"`
	ResidenceKota      string                `json:"residence_kota" form:"residence_kota"`
	ResidenceKecamatan string                `json:"residence_kecamatan" form:"residence_kecamatan"`
	ResidenceKelurahan string                `json:"residence_kelurahan" form:"residence_kelurahan"`
	ResidenceKodePOS   string                `json:"residence_kode_pos" form:"residence_kode_pos"`
	Status             string                `json:"status" form:"status"`
	Image              string                `json:"image" form:"image"`
	ImagePenerima      string                `json:"image_penerima" form:"image_penerima"`
	File               *multipart.FileHeader `json:"-" form:"file"`
	FilePenerima       *multipart.FileHeader `json:"-" form:"file_penerima"`
	UpdatedBy          string                `json:"updated_by" form:"updated_by"`
}

type ParticipantDone struct {
	Status        string                `json:"status" form:"status"`
	File          *multipart.FileHeader `json:"-" form:"file"`
	Image         string                `json:"image" form:"image"`
	ImagePenerima string                `json:"image_penerima" form:"image_penerima"`
	UpdatedBy     string                `json:"updated_by" form:"updated_by"`
}

type PartialDone struct {
	Status        string `json:"status" form:"status"`
	Image         string `json:"image" form:"image"`
	ImagePenerima string `json:"image_penerima" form:"image_penerima"`
	UpdatedBy     string `json:"updated_by" form:"updated_by"`
}

type Report struct {
	Provinsi string        `json:"provinsi"`
	Kota     string        `json:"kota"`
	Date     string        `json:"date"`
	Jam      template.HTML `json:"jam"`
	Evaluasi template.HTML `json:"evaluasi"`
	Solusi   template.HTML `json:"solusi"`
	Url      string
}
