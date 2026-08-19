package source

import (
	"backend/models"
	"errors"
)

type Platform int

const (
	PlatformInternal Platform = iota
	PlatformKitsu
	PlatformMAL
)

type ExternalSource interface {
	FetchDetails(title string) (*models.ShowDetail, error)
}

type Sources struct {
	TierList []ExternalSource
}

func NewSources(kitsu ExternalSource) *Sources {
	return &Sources{
		TierList: []ExternalSource{
			kitsu,
			// mal,
		},
	}
}

func (m *Sources) GetMetadata(title string) (*models.ShowDetail, error) {
	var errs []error
	for _, src := range m.TierList {
		meta, err := src.FetchDetails(title)
		if err == nil && meta != nil {
			return meta, nil
		}
		if err != nil {
			errs = append(errs, err)
		}
	}

	return nil, errors.Join(errs...)
}
