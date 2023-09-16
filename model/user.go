/*
 * Created on 15/09/23 01.59
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package model

import "time"

type User struct {
	ID        int        `json:"id" gorm:"primary_key"`
	Name      string     `json:"name"`
	Username  string     `json:"username"  gorm:"unique"`
	Password  string     `json:"-"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}