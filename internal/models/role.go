package models

import (

)

type Role struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `gorm:"size:255;uniqueIndex"`
	Description string    `gorm:"size:255"`
	
	Admins []Admin `gorm:"foreignKey:RoleID"`
}