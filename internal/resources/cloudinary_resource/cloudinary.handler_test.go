package cloudinaryresource

import (
	"strings"
	"testing"
)

// TestGenerateUniqueFilenameBoundsLength guards against the regression where an
// over-long original filename produced a Cloudinary public_id exceeding the
// 255-character limit, causing the upload to be silently rejected.
func TestGenerateUniqueFilenameBoundsLength(t *testing.T) {
	h := &CloudinaryHandlerImpl{}

	longBase := strings.Repeat("a", 300)
	got := h.generateUniqueFilename(longBase + ".png")

	if !strings.HasSuffix(got, ".png") {
		t.Fatalf("expected generated name to keep the .png extension, got %q", got)
	}

	// "uploads/" (8) + name must stay comfortably under Cloudinary's 255 limit.
	if len(got) > maxBaseNameLength+32 {
		t.Fatalf("generated name too long: %d chars (%q)", len(got), got)
	}
}

func TestGenerateUniqueFilenameKeepsShortNames(t *testing.T) {
	h := &CloudinaryHandlerImpl{}

	got := h.generateUniqueFilename("photo.png")

	if !strings.HasPrefix(got, "photo_") {
		t.Fatalf("expected short name to be preserved, got %q", got)
	}
	if !strings.HasSuffix(got, ".png") {
		t.Fatalf("expected .png extension, got %q", got)
	}
}
