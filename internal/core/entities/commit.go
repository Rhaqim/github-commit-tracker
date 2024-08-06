package entities

import (
	"time"

	"gorm.io/gorm"
)

// Commit represents the apllication commit entity.
type Commit struct {
	gorm.Model
	SHA             string    `json:"sha"`
	Author          string    `json:"author" gorm:"index"`
	Message         string    `json:"message"`
	Date            time.Time `json:"date"`
	URL             string    `json:"url"`
	OwnerRepository string    `json:"owner_repository,omitempty" gorm:"index"`
}

// CommitCount is a return type for the Top N Authors.
type CommitCount struct {
	Author      string
	CommitCount int
}
