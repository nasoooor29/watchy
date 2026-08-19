package source

import (
	"backend/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type KitsuSource struct {
	baseURL string
	client  *http.Client
}

func NewKitsuSource(baseURL string, client *http.Client) ExternalSource {
	if baseURL == "" || client == nil {
		return nil
	}

	return &KitsuSource{
		baseURL: baseURL,
		client:  client,
	}
}

// FetchDetails implements [ExternalSource].
func (k *KitsuSource) FetchDetails(title string) (*models.ShowDetail, error) {
	u, err := url.Parse(k.baseURL)
	if err != nil {
		return nil, err
	}
	u = u.JoinPath("anime")
	q := u.Query()
	q.Set("filter[text]", title)
	u.RawQuery = q.Encode()

	resp, err := k.client.Get(u.String())
	if err != nil {
		return nil, err
	}
	details := kitsuDetails{}
	err = json.NewDecoder(resp.Body).Decode(&details)
	if err != nil {
		return nil, err
	}
	fmt.Printf("details: %v\n", details)
	if len(details.Data) < 1 {
		return nil, models.ErrInternalServerError
	}
	data := details.Data[0]
	tags, err := k.fetchGenres(data.Relationships.Genres.Links.Related)
	if err != nil {
		return nil, err
	}

	return &models.ShowDetail{
		Title:        title,
		Alternative:  data.Attributes.Titles.En,
		Image:        data.Attributes.PosterImage.Original,
		Tags:         tags,
		Synopsis:     data.Attributes.Synopsis,
		Japanese:     data.Attributes.Titles.JaJp,
		EpisodeCount: data.Attributes.EpisodeCount,
	}, nil
}

func (k *KitsuSource) fetchGenres(u string) ([]string, error) {
	genres := []string{}
	resp, err := k.client.Get(u)
	if err != nil {
		return nil, err
	}
	details := kitsuGenres{}
	err = json.NewDecoder(resp.Body).Decode(&details)
	if err != nil {
		return nil, err
	}
	for _, v := range details.Data {
		genres = append(genres, v.Attributes.Name)
	}
	return genres, nil
}

type kitsuGenres struct {
	Data []struct {
		ID    string `json:"id"`
		Type  string `json:"type"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
		Attributes struct {
			CreatedAt   time.Time `json:"createdAt"`
			UpdatedAt   time.Time `json:"updatedAt"`
			Name        string    `json:"name"`
			Slug        string    `json:"slug"`
			Description string    `json:"description"`
		} `json:"attributes"`
	} `json:"data"`
	Meta struct {
		Count int `json:"count"`
	} `json:"meta"`
	Links struct {
		First string `json:"first"`
		Last  string `json:"last"`
	} `json:"links"`
}

type kitsuDetails struct {
	Data []struct {
		ID    string `json:"id"`
		Type  string `json:"type"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
		Attributes struct {
			CreatedAt           time.Time `json:"createdAt"`
			UpdatedAt           time.Time `json:"updatedAt"`
			Slug                string    `json:"slug"`
			Synopsis            string    `json:"synopsis"`
			Description         string    `json:"description"`
			CoverImageTopOffset int       `json:"coverImageTopOffset"`
			Titles              struct {
				En   string `json:"en"`
				EnJp string `json:"en_jp"`
				JaJp string `json:"ja_jp"`
			} `json:"titles"`
			CanonicalTitle    string   `json:"canonicalTitle"`
			AbbreviatedTitles []string `json:"abbreviatedTitles"`
			AverageRating     string   `json:"averageRating"`
			RatingFrequencies struct {
				Num2  string `json:"2"`
				Num3  string `json:"3"`
				Num4  string `json:"4"`
				Num5  string `json:"5"`
				Num6  string `json:"6"`
				Num7  string `json:"7"`
				Num8  string `json:"8"`
				Num9  string `json:"9"`
				Num10 string `json:"10"`
				Num11 string `json:"11"`
				Num12 string `json:"12"`
				Num13 string `json:"13"`
				Num14 string `json:"14"`
				Num15 string `json:"15"`
				Num16 string `json:"16"`
				Num17 string `json:"17"`
				Num18 string `json:"18"`
				Num19 string `json:"19"`
				Num20 string `json:"20"`
			} `json:"ratingFrequencies"`
			UserCount      int         `json:"userCount"`
			FavoritesCount int         `json:"favoritesCount"`
			StartDate      string      `json:"startDate"`
			EndDate        string      `json:"endDate"`
			NextRelease    interface{} `json:"nextRelease"`
			PopularityRank int         `json:"popularityRank"`
			RatingRank     int         `json:"ratingRank"`
			AgeRating      string      `json:"ageRating"`
			AgeRatingGuide string      `json:"ageRatingGuide"`
			Subtype        string      `json:"subtype"`
			Status         string      `json:"status"`
			Tba            interface{} `json:"tba"`
			PosterImage    struct {
				Tiny     string `json:"tiny"`
				Large    string `json:"large"`
				Small    string `json:"small"`
				Medium   string `json:"medium"`
				Original string `json:"original"`
				Meta     struct {
					Dimensions struct {
						Tiny struct {
							Width  int `json:"width"`
							Height int `json:"height"`
						} `json:"tiny"`
						Large struct {
							Width  int `json:"width"`
							Height int `json:"height"`
						} `json:"large"`
						Small struct {
							Width  int `json:"width"`
							Height int `json:"height"`
						} `json:"small"`
						Medium struct {
							Width  int `json:"width"`
							Height int `json:"height"`
						} `json:"medium"`
					} `json:"dimensions"`
				} `json:"meta"`
			} `json:"posterImage"`
			CoverImage struct {
				Tiny     string `json:"tiny"`
				Large    string `json:"large"`
				Small    string `json:"small"`
				Original string `json:"original"`
				Meta     struct {
					Dimensions struct {
						Tiny struct {
							Width  int `json:"width"`
							Height int `json:"height"`
						} `json:"tiny"`
						Large struct {
							Width  int `json:"width"`
							Height int `json:"height"`
						} `json:"large"`
						Small struct {
							Width  int `json:"width"`
							Height int `json:"height"`
						} `json:"small"`
					} `json:"dimensions"`
				} `json:"meta"`
			} `json:"coverImage"`
			EpisodeCount   int    `json:"episodeCount"`
			EpisodeLength  int    `json:"episodeLength"`
			TotalLength    int    `json:"totalLength"`
			YoutubeVideoID string `json:"youtubeVideoId"`
			ShowType       string `json:"showType"`
			Nsfw           bool   `json:"nsfw"`
		} `json:"attributes"`
		Relationships struct {
			Genres struct {
				Links struct {
					Self    string `json:"self"`
					Related string `json:"related"`
				} `json:"links"`
			} `json:"genres"`
			Categories struct {
				Links struct {
					Self    string `json:"self"`
					Related string `json:"related"`
				} `json:"links"`
			} `json:"categories"`
			Castings struct {
				Links struct {
					Self    string `json:"self"`
					Related string `json:"related"`
				} `json:"links"`
			} `json:"castings"`
			Installments struct {
				Links struct {
					Self    string `json:"self"`
					Related string `json:"related"`
				} `json:"links"`
			} `json:"installments"`
			Mappings struct {
				Links struct {
					Self    string `json:"self"`
					Related string `json:"related"`
				} `json:"links"`
			} `json:"mappings"`
			Reviews struct {
				Links struct {
					Self    string `json:"self"`
					Related string `json:"related"`
				} `json:"links"`
			} `json:"reviews"`
			MediaRelationships struct {
				Links struct {
					Self    string `json:"self"`
					Related string `json:"related"`
				} `json:"links"`
			} `json:"mediaRelationships"`
			Characters struct {
				Links struct {
					Self    string `json:"self"`
					Related string `json:"related"`
				} `json:"links"`
			} `json:"characters"`
			Staff struct {
				Links struct {
					Self    string `json:"self"`
					Related string `json:"related"`
				} `json:"links"`
			} `json:"staff"`
			Productions struct {
				Links struct {
					Self    string `json:"self"`
					Related string `json:"related"`
				} `json:"links"`
			} `json:"productions"`
			Quotes struct {
				Links struct {
					Self    string `json:"self"`
					Related string `json:"related"`
				} `json:"links"`
			} `json:"quotes"`
			Episodes struct {
				Links struct {
					Self    string `json:"self"`
					Related string `json:"related"`
				} `json:"links"`
			} `json:"episodes"`
			StreamingLinks struct {
				Links struct {
					Self    string `json:"self"`
					Related string `json:"related"`
				} `json:"links"`
			} `json:"streamingLinks"`
			AnimeProductions struct {
				Links struct {
					Self    string `json:"self"`
					Related string `json:"related"`
				} `json:"links"`
			} `json:"animeProductions"`
			AnimeCharacters struct {
				Links struct {
					Self    string `json:"self"`
					Related string `json:"related"`
				} `json:"links"`
			} `json:"animeCharacters"`
			AnimeStaff struct {
				Links struct {
					Self    string `json:"self"`
					Related string `json:"related"`
				} `json:"links"`
			} `json:"animeStaff"`
		} `json:"relationships"`
	} `json:"data"`
	Meta struct {
		Count int `json:"count"`
	} `json:"meta"`
	Links struct {
		First string `json:"first"`
		Last  string `json:"last"`
	} `json:"links"`
}
