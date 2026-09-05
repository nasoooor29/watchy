package models

type Show struct {
	Path     string
	Seasons  []Season
	Metadata Metadata
}

type Season struct {
	Name     string
	Episodes []Episode
}

type Metadata struct {
	Title       string
	Alternative string
	Image       string
	Tags        []string
	Progress    string
	Percent     int
	Synopsis    string
	Japanese    string
	LatestEp    int64
}

type Episode struct {
	Index   int
	Watched bool
	Paths   []string
}
