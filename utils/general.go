package utils

import "fmt"

func FlattenMap(m map[string]any) map[string]string {
	out := make(map[string]string)

	var flatten func(prefix string, value any)

	flatten = func(prefix string, value any) {
		switch v := value.(type) {
		case map[string]any:
			for key, value := range v {
				newKey := key
				if prefix != "" {
					newKey = prefix + "_" + key
				}

				flatten(newKey, value)
			}

		default:
			out[prefix] = fmt.Sprint(v)
		}
	}

	for key, value := range m {
		flatten(key, value)
	}

	return out
}
