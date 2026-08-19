package source

import "backend/models"

type Platform int

const (
	PlatformInternal Platform = iota
	PlatformKitsu
	PlatformMAL
)

type ID struct {
	id       int
	Platform Platform
}

type ExternalSource interface {
	FetchDetails(id ID) (*models.ShowDetail, error)
}
