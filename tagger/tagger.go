package main

import (
	"regexp"

	"github.com/guad/paperless2/backend/model"
)

func tag(content string, tags []model.Tag) []string {
	applied := []string{}

	bytes := []byte(content)

	for _, tag := range tags {
		r := regexp.MustCompile("(?i)" + tag.Regex)

		if r.Match(bytes) {
			applied = append(applied, tag.Name)
		}
	}

	return applied
}
