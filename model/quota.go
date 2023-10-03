/*
 * Created on 22/09/23 00.43
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package model

import "time"

type Quota struct {
	ID        int        `json:"id" gorm:"primary_key"`
	Total     int64      `json:"total"`
	Provinsi  string     `json:"provinsi"`
	Kota      string     `json:"kota"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func (i *Quota) TableName() string {
	return "quotas"
}
