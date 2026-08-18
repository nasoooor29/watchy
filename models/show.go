package models

import "github.com/nssteinbrenner/anitogo"

type Show struct {
	Title string
	Image string // it's optional but if it's empty it will be replaced with a placeholder image

	Episodes []Episode

	Metadata map[string]string
}

type Episode struct {
	Name        string
	Number      int
	Watched     bool
	Path        string
	anitomyInfo *anitogo.Elements
}

func GroupEpisodes(episodes []Episode) map[int][]Episode {
	grouped := make(map[int][]Episode)
	for _, episode := range episodes {
		grouped[episode.Number] = append(grouped[episode.Number], episode)
	}
	return grouped
}
