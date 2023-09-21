/*
 * Created on 18/09/23 06.37
 *
 * Copyright (c) 2023 Abdul Ghani Abbasi
 */

package model

import "time"

type Organization struct {
	ID        int        `json:"id" gorm:"primary_key"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}
