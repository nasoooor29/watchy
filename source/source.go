package source

import (
	"backend/models"
	"errors"
)

type ExternalSource interface {
	FetchDetails(title string) (*models.Metadata, error)
}

type Sources struct {
	TierList []ExternalSource
}

func NewSources(sources ...ExternalSource) *Sources {
	return &Sources{
		TierList: sources,
	}
}

func (m *Sources) GetMetadata(title string) (*models.Metadata, error) {
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
