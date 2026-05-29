package storage

import (
	"errors"
)

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
)

type Link struct {
	id           string // `json:"id"`
	Short_code   string `json:"short_code,omitempty"`
	Original_url string `json:"url,omitempty"`
	Created_at   string `json:"created_at,omitempty"`
	Visits       int    `json:"visits,omitempty"`
}
