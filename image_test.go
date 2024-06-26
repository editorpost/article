package article_test

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/editorpost/article"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestImageConversions is a table-driven test for the Image struct.
// It verifies the conversion of map data to Image struct, the validation process, and handling of zero-value fields.
//
// Explanation of test cases:
// - Valid Image: Ensures that valid data is correctly converted into an Image struct without errors.
// - Invalid URL: Ensures that an invalid URL triggers a validation error.
// - Missing Required Fields: Ensures that missing mandatory fields trigger a validation error.
// - Zero Value Fields: Ensures that empty field values are handled correctly and trigger a validation error.
func TestImageConversions(t *testing.T) {
	tests := []struct {
		name          string
		inputMap      map[string]any
		expectedImage *article.Image
		expectError   bool
	}{
		{
			name: "Valid Image",
			inputMap: map[string]any{
				"id":     "123e4567-e89b-12d3-a456-426614174000",
				"url":    "https://example.com/image.jpg",
				"alt":    "An example image",
				"width":  800,
				"height": 600,
				"title":  "An example title",
			},
			expectedImage: &article.Image{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				URL:    "https://example.com/image.jpg",
				Alt:    "An example image",
				Width:  800,
				Height: 600,
				Title:  "An example title",
			},
			expectError: false,
		},
		{
			name: "Invalid URL",
			inputMap: map[string]any{
				"url":    "invalid-url",
				"alt":    "An example image",
				"width":  800,
				"height": 600,
				"title":  "An example title",
			},
			expectedImage: nil,
			expectError:   true,
		},
		{
			name: "Missing Required Fields",
			inputMap: map[string]any{
				"url":    "",
				"alt":    "An example image",
				"width":  800,
				"height": 600,
				"title":  "An example title",
			},
			expectedImage: nil,
			expectError:   true,
		},
		{
			name: "Zero Value Fields",
			inputMap: map[string]any{
				"url":    "",
				"alt":    "",
				"width":  0,
				"height": 0,
				"title":  "",
			},
			expectedImage: &article.Image{
				URL:    "",
				Alt:    "",
				Width:  0,
				Height: 0,
				Title:  "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img, err := article.NewImageFromMap(tt.inputMap)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedImage, img)
				assert.Equal(t, tt.inputMap, img.Map())
			}
		})
	}
}

func TestImageNormalize(t *testing.T) {
	img := &article.Image{
		URL:    "  " + gofakeit.URL() + "  ",
		Alt:    "  " + gofakeit.Sentence(5) + "  ",
		Width:  gofakeit.Number(800, 1920),
		Height: gofakeit.Number(600, 1080),
		Title:  "  " + gofakeit.Sentence(10) + "  ",
	}

	img.Normalize()

	assert.NotEmpty(t, img.URL)
	assert.Equal(t, strings.TrimSpace(img.URL), img.URL)
	assert.Equal(t, strings.TrimSpace(img.Alt), img.Alt)
	assert.Equal(t, strings.TrimSpace(img.Title), img.Title)
}
