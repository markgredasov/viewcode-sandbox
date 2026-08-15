package utils_string

import (
	"encoding/json"
	"strings"
)

func CleanString(s string) string {
	var decoded string
	if err := json.Unmarshal([]byte(`"`+s+`"`), &decoded); err == nil {
		s = decoded
	}

	var builder strings.Builder
	for _, r := range s {
		if r == '\n' || r == '\t' || (r >= 32 && r != '\\') {
			builder.WriteRune(r)
		}
	}
	s = builder.String()

	lines := removeSpecialSymbols(strings.Split(s, "\n"))

	return strings.Join(lines, "\n")
}

func removeSpecialSymbols(lines []string) []string {
	for i, line := range lines {
		if len(line) > 2 && line[0] >= 'A' && line[0] <= 'Z' && line[1] == ' ' {
			lines[i] = line[2:]
		}
		if len(line) > 0 && strings.ContainsAny(string(line[0]), "+#@%&!?<>;0123456789") {
			idx := strings.IndexFunc(line, func(r rune) bool {
				return r == '=' || r == '-'
			})
			if idx > 0 {
				lines[i] = line[idx:]
			}
		}
		lines[i] = strings.TrimSpace(lines[i])
	}

	return lines
}
