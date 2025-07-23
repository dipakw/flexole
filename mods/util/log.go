package util

import "strings"

func LogKindsToFlag(kinds []string) uint8 {
	flag := uint8(0)

	for _, kind := range kinds {
		switch kind {
		case "info":
			flag |= 1
		case "warn":
			flag |= 2
		case "error":
			flag |= 4
		}
	}

	return flag
}

func LogShortToKinds(short string) []string {
	kinds := []string{}

	if short == "off" {
		return kinds
	}

	if short == "all" {
		return []string{"info", "warn", "error"}
	}

	if strings.Contains(short, "i") {
		kinds = append(kinds, "info")
	}

	if strings.Contains(short, "w") {
		kinds = append(kinds, "warn")
	}

	if strings.Contains(short, "e") {
		kinds = append(kinds, "error")
	}

	return kinds
}
