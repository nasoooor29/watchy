package web

import "fmt"

func dummyEpisodes() []Episode {
	episodes := make([]Episode, 55)
	for i := range episodes {
		number := 18 - i
		episodes[i] = Episode{
			Name:    fmt.Sprintf("E%02d - Honzuki No Gekokujou S4", number),
			Watched: number < 18,
			Sources: dummySources(),
		}
	}
	return episodes
}
