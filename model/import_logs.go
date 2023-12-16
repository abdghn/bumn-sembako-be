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
	Reference   string     `json:"reference" gorm:"type:varchar(255)"`
	UploadedBy  string     `json:"uploaded_by" gorm:"type:varchar(255)"`
	Type        string     `json:"type" gorm:"type:varchar(100)"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `sql:"index" json:"deleted_at"`
}
