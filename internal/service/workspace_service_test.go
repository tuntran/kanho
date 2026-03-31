package service_test

import (
	"testing"

	"github.com/tungtran/kanho/internal/domain"
)

func TestCreateWorkspace_AutoGeneratesSlug(t *testing.T) {
	// UniqueSlug is the core logic used by WorkspaceService.Create
	slug, err := domain.UniqueSlug("My New Workspace", func(slug string) (bool, error) {
		return false, nil // no conflicts
	})
	if err != nil {
		t.Fatalf("UniqueSlug: %v", err)
	}

	if slug == "" {
		t.Error("slug should not be empty")
	}
	if slug == "My New Workspace" {
		t.Error("slug should be transformed, not raw name")
	}
}

func TestCreateWorkspace_SlugUnique(t *testing.T) {
	callCount := 0
	slug, err := domain.UniqueSlug("Test Workspace", func(slug string) (bool, error) {
		callCount++
		// First slug attempt conflicts, second does not
		return callCount == 1, nil
	})
	if err != nil {
		t.Fatalf("UniqueSlug: %v", err)
	}

	if callCount < 2 {
		t.Error("should have checked slug existence at least twice due to conflict")
	}
	if slug == "" {
		t.Error("should produce a valid slug after retry")
	}
}
