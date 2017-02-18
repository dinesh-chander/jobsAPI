package jobs

import "github.com/jinzhu/gorm"
import _ "github.com/jinzhu/gorm/dialects/mysql"

type Job struct {
	gorm.Model
	Apply          string `gorm:"size:2048"`
	Title          string `gorm:"size:200"`
	Address        string `gorm:"size:500"`
	Channel_Name   string `gorm:"size:20"`
	City           string `gorm:"size:100"`
	Company        string `gorm:"size:100"`
	Country        string `gorm:"size:100"`
	Job_Type       string `gorm:"size:20"`
	Description    string `gorm:"type:text"`
	Published_Date uint64 `gorm:"type:bigint"`
	Is_Remote      bool
	Source         string `gorm:"size:200"`
	Source_Name    string `gorm:"size:20"`
	Source_Id      string `gorm:"size:50"`
	Tags           string `gorm:"size:300"`
	Approved       bool
}
