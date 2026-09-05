package models

type Show struct {
	Path     string
	Seasons  map[string][]Episode
	Metadata Metadata
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
}

type Episode struct {
	Index   int
	Watched bool
	Paths   []string
}
