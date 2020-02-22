package main

import (
	"regexp"

	"github.com/guad/paperless2/backend/model"
)

func tag(content string, tags []model.Tag) []string {
	applied := []string{}
	appliedSet := map[string]struct{}{}

	bytes := []byte(content)

	tagSet := map[string]model.Tag{}

	for _, t := range tags {
		tagSet[t.Name] = t
	}

	for _, tag := range tags {
		if tag.Regex == "" {
			continue
		}

		r := regexp.MustCompile("(?i)" + tag.Regex)

		if r.Match(bytes) {
			applied = append(applied, tag.Name)
			appliedSet[tag.Name] = struct{}{}
		}
	}

	// Propagate implies

	changed := true

	for changed {
		changed = false

		for _, t := range applied {
			tag := tagSet[t]

			if len(tag.Implies) > 0 {
				for _, implied := range tag.Implies {
					if _, ok := appliedSet[implied]; !ok {
						appliedSet[implied] = struct{}{}
						applied = append(applied, implied)
						changed = true
					}
				}
			}
		}

	}

	return applied
}
