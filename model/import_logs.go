package model

import "time"

type ImportLog struct {
	ID          int        `json:"id" gorm:"primary_key"`
	FileName    string     `json:"file_name" gorm:"type:varchar(255)"`
	Status      string     `json:"status"`
	TotalRows   int        `json:"total_rows"`
	SuccessRows int        `json:"success_rows"`
	FailedRows  int        `json:"failed_rows"`
	Path        string     `json:"path" gorm:"type: text"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `sql:"index" json:"deleted_at"`
}
