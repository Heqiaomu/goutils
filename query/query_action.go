package query

import (
	"regexp"
	"strings"
)

type Action string

var likeRegex = regexp.MustCompile(`\((.*)\)`)

func (a Action) parseValue(s string) string {
	switch a {
	case IN:
		s := strings.ReplaceAll(s, "@", ",")
		sub := likeRegex.FindStringSubmatch(s)
		return sub[1]
	case LIKE:
		s := strings.ReplaceAll(s, "#", "%")
		return s
	default:
		return s
	}
}

var ActionMap = map[string]Action{
	"eq":   EQ,
	"ne":   NE,
	"gt":   GT,
	"gte":  GE,
	"lt":   LT,
	"let":  LE,
	"in":   IN,
	"like": LIKE,
}

const (
	EQ   Action = "="
	NE   Action = "<>"
	LT   Action = "<"
	LE   Action = "<="
	GT   Action = ">"
	GE   Action = ">="
	IN   Action = "in"
	LIKE Action = "like"
)
