/*
 * Created on 15/09/23 15.28
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package model

import (
	"gorm.io/gorm"
	"time"
)

type Province struct {
	gorm.Model `json:"-"`
	ID         int        `json:"id" gorm:"primary_key"`
	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `sql:"index" json:"deleted_at"`
}

type Regency struct {
	gorm.Model `json:"-"`
	ID         int        `json:"id" gorm:"primary_key"`
	ProvinceID uint       `json:"province_id" gorm:"column:province_id"`
	Province   Province   `json:"province" gorm:"foreignKey:ProvinceID"`
	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `sql:"index" json:"deleted_at"`
}

type District struct {
	gorm.Model `json:"-"`
	ID         int        `json:"id" gorm:"primary_key"`
	RegencyID  uint       `json:"regency_id" gorm:"column:regency_id"`
	Regency    Regency    `json:"regency" gorm:"foreignKey:RegencyID"`
	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `sql:"index" json:"deleted_at"`
}

type Village struct {
	gorm.Model `json:"-"`
	ID         int        `json:"id" gorm:"primary_key"`
	DistrictID uint       `json:"district_id" gorm:"column:district_id"`
	District   District   `json:"district" gorm:"foreignKey:DistrictID"`
	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `sql:"index" json:"deleted_at"`
}
