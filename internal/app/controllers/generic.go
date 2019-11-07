package controllers

import (
	"github.com/nathanhnew/bowl-backend/internal/app/db"
	"regexp"
	"strconv"
	"strings"
)

func newSlug(slugType string, name string) (string, error) {
	var existingSlug string
	var err error
	slugBase := slugify(name)
	if slugType == "league" {
		existingSlug, err = db.GetLeagueSlugByBase(slugBase)
	} else if slugType == "bowl" {
		existingSlug, err = db.GetBowlSlugByBase(slugBase)
	}
	if err != nil {
		return "", err
	}
	splitSlug := strings.Split(existingSlug, "-")
	slugIndex := splitSlug[len(splitSlug)-1]
	nextIndex, err := strconv.Atoi(slugIndex)
	if err != nil {
		if existingSlug == "" {
			nextIndex = -1
		} else {
			nextIndex = 0
		}
	}
	nextIndex += 1
	if nextIndex == 0 {
		return slugBase, nil
	} else {
		return slugBase + "-" + strconv.Itoa(nextIndex), nil
	}
}

func slugify(name string) string {
	slug := strings.ToLower(name)
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	spaces := regexp.MustCompile(`\s+`)
	slug = reg.ReplaceAllString(slug, "")
	slug = spaces.ReplaceAllString(slug, " ")
	slug = strings.TrimSpace(slug)
	if len(slug) > 32 {
		slug = slug[:32]
		slug = strings.TrimSpace(slug)
	}
	slug = strings.Replace(slug, " ", "-", -1)
	return slug
}
