package web

import (
	"encoding/json"
	"strings"

	"github.com/nssteinbrenner/anitogo"
)

func metadataAttributes(metadata any) map[string]any {
	encoded, err := json.Marshal(metadata)
	if err != nil {
		return map[string]any{}
	}

	fields := make(map[string]any)
	if err := json.Unmarshal(encoded, &fields); err != nil {
		return map[string]any{}
	}

	attributes := make(map[string]any, len(fields))
	for name, value := range fields {
		attributes[strings.ToLower(name)] = value
	}

	return attributes
}

// create custom template functions for the templates
var funcs = map[string]any{
	"metadataAttributes": metadataAttributes,
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
