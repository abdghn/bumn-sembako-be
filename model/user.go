/*
 * Created on 15/09/23 01.59
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model     `json:"-"`
	ID             int          `json:"id" gorm:"primary_key"`
	Name           string       `json:"name"`
	Username       string       `json:"username"  gorm:"unique"`
	Password       string       `json:"-"`
	Role           string       `json:"role"`
	Provinsi       string       `json:"provinsi"`
	Kota           string       `json:"kota"`
	OrganizationID uint         `json:"organization_id" gorm:"column:organization_id"`
	Organization   Organization `json:"organization" gorm:"foreignKey:OrganizationID"`
	RetryAttempts  int64        `json:"retry_attempts"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	DeletedAt      *time.Time   `sql:"index" json:"deleted_at"`
}
