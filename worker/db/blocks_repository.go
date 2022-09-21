package db

import (
	"github.com/Hendryboyz/eth-synchronizer/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BlocksRepository struct {
	DB *gorm.DB
}

func (b *BlocksRepository) CreateBlock(block models.Blocks) {
	b.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "number"}},
		UpdateAll: true,
	}).Create(&block)
}
