package models

type LibraryPage struct {
	TotalSeries int
	Shows       []LibraryShow
	Selected    ShowDetail
	Stats       []LibraryStat
}

type LibraryShow struct {
	ID       int64
	Title    string
	Year     int
	Score    float64
	Progress string
	Percent  int
	Image    string
}

type LibraryEpisode struct {
	Name    string
	Watched bool
	Sources []EpisodeSource
}

type ShowDetail struct {
	Title        string
	Alternative  string
	Image        string
	Tags         []string
	Progress     string
	Percent      int
	Synopsis     string
	Japanese     string
	EpisodeCount int
	Season       Season
}

type LibraryStat struct {
	Label string
	Value string
}

type Season struct {
	Name     string
	Episodes []LibraryEpisode
}

type EpisodeSource struct {
	Label string
}
