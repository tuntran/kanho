package domain

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

var (
	nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)
	leadTrailDash   = regexp.MustCompile(`^-+|-+$`)
)

// GenerateSlug converts a name into a URL-friendly slug.
func GenerateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = nonAlphanumeric.ReplaceAllString(slug, "-")
	slug = leadTrailDash.ReplaceAllString(slug, "")
	if len(slug) > 100 {
		slug = slug[:100]
	}
	if slug == "" {
		slug = "untitled"
	}
	return slug
}

// UniqueSlug generates a unique slug by appending suffixes if needed.
// exists is called to check if a slug is already taken.
func UniqueSlug(base string, exists func(string) (bool, error)) (string, error) {
	slug := GenerateSlug(base)

	taken, err := exists(slug)
	if err != nil {
		return "", err
	}
	if !taken {
		return slug, nil
	}

	for i := 2; i <= 99; i++ {
		candidate := fmt.Sprintf("%s-%d", slug, i)
		taken, err = exists(candidate)
		if err != nil {
			return "", err
		}
		if !taken {
			return candidate, nil
		}
	}

	return slug + "-" + uuid.New().String()[:8], nil
}
