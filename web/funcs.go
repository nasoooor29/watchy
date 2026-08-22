package web

import (
	"strings"

	"github.com/nssteinbrenner/anitogo"
)

// create custom template functions for the templates
var funcs = map[string]any{
	"formatEpisodeName": func(name string) string {
		parsed := anitogo.Parse(name, anitogo.DefaultOptions)

		parts := make([]string, 0, 4)

		if parsed.VideoResolution != "" {
			parts = append(parts, parsed.VideoResolution)
		}

		parts = append(parts, parsed.VideoTerm...)
		parts = append(parts, parsed.AudioTerm...)
		parts = append(parts, parsed.Subtitles...)

		// Things like MultiSub are usually left as unknown terms.
		for _, term := range parsed.Unknown {
			lower := strings.ToLower(term)

			if strings.Contains(lower, "sub") {
				parts = append(parts, term)
			}
		}

		if len(parts) == 0 {
			return name
		}

		return strings.Join(parts, " • ")
	},
}
