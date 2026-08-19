package source

import (
	"backend/models"
	"net/http"
)

type JikanSource struct {
	baseURL string
	client  *http.Client
}

func NewJikanSource(baseURL string, client *http.Client) ExternalSource {
	if baseURL == "" || client == nil {
		return nil
	}

	return &JikanSource{
		baseURL: baseURL,
		client:  client,
	}
}

// FetchDetails implements [ExternalSource].
func (j *JikanSource) FetchDetails(title string) (*models.ShowDetail, error) {
	panic("unimplemented")
}
