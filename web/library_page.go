package web

type LibraryPage struct {
	TotalSeries int
	Shows       []LibraryShow
	Selected    ShowDetail
	Stats       []LibraryStat
}

type LibraryShow struct {
	Title    string
	Year     int
	Score    float64
	Progress string
	Percent  int
	Image    string
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
	Episodes []Episode
}

type Episode struct {
	Name    string
	Watched bool
	Sources []EpisodeSource
}

type EpisodeSource struct {
	Label string
}

func DummyLibraryPage() LibraryPage {
	shows := []LibraryShow{
		{Title: "Honzuki no Gekokujou Shisho ni Naru Tame ni wa Shudan wo...", Year: 2026, Score: 7.98, Progress: "54/55 watched", Percent: 98, Image: "https://images.unsplash.com/photo-1515879218367-8466d910aaa4?auto=format&fit=crop&w=180&q=80"},
		{Title: "Nige Jouzu no Wakagimi", Year: 2024, Score: 7.78, Progress: "13/18 watched", Percent: 72},
		{Title: "Buchigire Reijou wa Houfuku wo Chikaimashita", Year: 2026, Score: 7.08, Progress: "0/0 watched", Percent: 0},
		{Title: "Ryoumin 0-nin Start no Henkyou Ryoushu-sama", Year: 2026, Score: 7.04, Progress: "7/7 watched", Percent: 100},
		{Title: "Tensei shitara Slime Datta Ken", Year: 2018, Score: 8.13, Progress: "95/95 watched", Percent: 100},
		{Title: "Koko wa Ore ni Makasete Saki ni Ike", Year: 2026, Score: 6.78, Progress: "7/7 watched", Percent: 100},
		{Title: "Hell Mode: Yarikomizuki no Gamer", Year: 2026, Score: 7.2, Progress: "30/30 watched", Percent: 100},
		{Title: "Otome Kaijuu Carameliser", Year: 2026, Score: 7.67, Progress: "7/7 watched", Percent: 100},
	}

	return LibraryPage{
		TotalSeries: 177,
		Shows:       shows[:6],
		Selected: ShowDetail{
			Title:        "Honzuki no Gekokujou Shisho ni Naru Tame ni wa Shudan wo Erandeiraremasen",
			Alternative:  "Ascendance of a Bookworm: Adopted Daughter of an Archduke",
			Image:        shows[0].Image,
			Tags:         []string{"TV", "Currently Airing", "Spring 2026", "Score 7.98", "Wit Studio", "Fantasy"},
			Progress:     "54/55 watched",
			Percent:      98,
			Synopsis:     "Anime adaptation of part three of the Honzuki no Gekokujou light novel.",
			Japanese:     "本好きの下剋上 ～司書になるためには手段を選んでいられません～",
			EpisodeCount: 55,
			Season: Season{Name: "Season 04", Episodes: []Episode{
				{Name: "E18 - Honzuki No Gekokujou S4", Sources: dummySources()},
				{Name: "E17 - Honzuki No Gekokujou S4", Watched: true, Sources: dummySources()},
				{Name: "E16 - Honzuki No Gekokujou S4", Watched: true, Sources: dummySources()},
				{Name: "E15 - Honzuki No Gekokujou S4", Watched: true, Sources: dummySources()},
			}},
		},
		Stats: []LibraryStat{{Label: "Watching", Value: "0"}, {Label: "Missing eps", Value: "0"}, {Label: "Missing seasons", Value: "0"}, {Label: "MAL eps", Value: "-"}},
	}
}

func dummySources() []EpisodeSource {
	return []EpisodeSource{{Label: "1080p · WEBRip · HEVC · AAC"}, {Label: "1080p · WEB-DL · AVC · AAC"}}
}
