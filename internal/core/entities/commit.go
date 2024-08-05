package entities

import (
	"time"

	"gorm.io/gorm"
)

type Commit struct {
	gorm.Model
	SHA             string    `json:"sha"`
	Author          string    `json:"author"`
	Message         string    `json:"message"`
	Date            time.Time `json:"date"`
	URL             string    `json:"url"`
	OwnerRepository string    `json:"owner_repository,omitempty"`
}

type CommitCount struct {
	Author      string
	CommitCount int
}
