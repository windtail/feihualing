package main

import (
	"strconv"
	"strings"
)

type Search struct {
	Id        int64
	Title     string
	Dynasty   string
	Author    string
	Content   []string
	FavorOnly bool
}

func EmptySearch() *Search {
	return &Search{
		Id:        0,
		Title:     "",
		Dynasty:   "",
		Author:    "",
		Content:   make([]string, 0),
		FavorOnly: false,
	}
}

func NewSearch(rule string, favorOnly bool) *Search {
	s := EmptySearch()

	rule = strings.TrimSpace(rule)
	parts := strings.Split(rule, " ")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		if strings.HasPrefix(part, "t") {
			s.Title = part[1:]
		} else if strings.HasPrefix(part, "d") {
			s.Dynasty = part[1:]
		} else if strings.HasPrefix(part, "a") {
			s.Author = part[1:]
		} else {
			if i, err := strconv.ParseUint(part, 10, 64); err != nil {
				s.Content = append(s.Content, part)
			} else {
				s.Id = int64(i)
			}
		}
	}

	s.FavorOnly = favorOnly

	return s
}