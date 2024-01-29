package model

import "time"

type Owner struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	FistName string `gorm:"not null"`
	LastName string `gorm:"not null"`
	Email    string `gorm:"not null"`
}

func (o *Owner) FullName() string {
	return o.FistName + " " + o.LastName
}

type OwnerAccount struct {
	Owner
	Password  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}
