package models

import "time"

type Blocks struct {
	Number     int64  `gorm:"primaryKey"`
	Hash       string `gorm:"index:block_hash,unique"`
	ParentHash string
	Time       time.Time
}
