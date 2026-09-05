package source

import (
	"backend/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type TenraiSource struct {
	baseURL string
	client  *http.Client
}

func NewTenraiSource(baseURL string, client *http.Client) ExternalSource {
	if baseURL == "" || client == nil {
		return nil
	}

	return &TenraiSource{baseURL: baseURL, client: client}
}

// FetchDetails implements [ExternalSource].
func (t *TenraiSource) FetchDetails(title string) (*models.ShowDetail, error) {
	u, err := url.Parse(t.baseURL)
	if err != nil {
		return nil, err
	}
	u = u.JoinPath("anime")
	q := u.Query()
	q.Set("q", title)
	q.Set("limit", "1")
	u.RawQuery = q.Encode()

	resp, err := t.client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("tenrai search failed: %s", resp.Status)
	}

	var results tenraiSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}
	if len(results.Data) == 0 {
		return nil, models.ErrInternalServerError
	}
	result := results.Data[0]
	image := result.Images.JPG.LargeImageURL
	if image == "" {
		image = result.Images.WEBP.LargeImageURL
	}
	if image == "" {
		image = result.Images.JPG.ImageURL
	}

	return &models.ShowDetail{
		ID:           int64(result.MALID),
		Title:        result.Title,
		Alternative:  result.TitleEnglish,
		Image:        image,
		Synopsis:     result.Synopsis,
		Japanese:     result.TitleJapanese,
		EpisodeCount: result.Episodes,
	}, nil
}

type tenraiSearchResponse struct {
	Data []tenraiAnime `json:"data"`
}

type tenraiAnime struct {
	MALID         int    `json:"mal_id"`
	Title         string `json:"title"`
	TitleEnglish  string `json:"title_english"`
	TitleJapanese string `json:"title_japanese"`
	Synopsis      string `json:"synopsis"`
	Episodes      int    `json:"episodes"`
	Images        struct {
		JPG struct {
			ImageURL      string `json:"image_url"`
			LargeImageURL string `json:"large_image_url"`
		} `json:"jpg"`
		WEBP struct {
			ImageURL      string `json:"image_url"`
			LargeImageURL string `json:"large_image_url"`
		} `json:"webp"`
	} `json:"images"`
}
