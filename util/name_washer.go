package util

import (
	"log"
	"strings"
)

func WashName(name string) string {
	split := strings.Split(name, " ")
	switch len(split) {
	case 1:
		return strings.ReplaceAll(split[0], "\\'", "")
	case 2:
		second := split[1]

		titleLower := strings.ToLower(second)
		if titleLower == "jr." || titleLower == "sr." {
			return split[0]
		}
		break
	default:
		// TODO: Figure out the implications of this
		log.Println("This should not happen!! name=", name)
		return split[len(split)-1]
	}
	return name
}
