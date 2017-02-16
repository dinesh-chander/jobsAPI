package jobs

import "github.com/jinzhu/gorm"

type Job struct {
	gorm.Model
	Apply          string
	Title          string
	Address        string
	Channel_Name   string
	City           string
	Company        string
	Country        string
	Job_Type       string
	Description    string
	Published_Date uint64
	Is_Remote      bool
	Source         string
	Source_Name    string
	Source_Id      string
	Tags           string
	Approved       bool
}
