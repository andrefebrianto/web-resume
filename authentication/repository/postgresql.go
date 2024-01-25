package repository

import (
	"devoratio.dev/web-resume/model"
	"gorm.io/gorm"
)

type PostgreSQLDatabase struct {
	db *gorm.DB
}

func NewPostgreSQL(db *gorm.DB) *PostgreSQLDatabase {
	return &PostgreSQLDatabase{
		db: db,
	}
}

func (p *PostgreSQLDatabase) GetOwnerByUsernameOrEmail(identifier string) (*model.OwnerAccount, error) {
	var owner *model.OwnerAccount
	result := p.db.Where("username = ? OR email = ?", identifier, identifier).First(owner)

	return owner, result.Error
}
