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
	ID      int
	Sources []EpisodeSource
}

type ShowDetail struct {
	ID           int64
	Title        string
	Alternative  string
	Image        string
	Tags         []string
	Progress     string
	Percent      int
	Synopsis     string
	Japanese     string
	EpisodeCount int
	Seasons      []SeasonOption
	Season       Season
}

type SeasonOption struct {
	Index   int
	Name    string
	Current bool
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
	Path  string
	Label string
}
