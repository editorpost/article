package article_test

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/editorpost/article"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMediaConversions is a table-driven test for the Media struct.
// It verifies the conversion of map data to Media struct, the validation process, and handling of zero-value fields.
//
// Explanation of test cases:
// - Valid Media: Ensures that valid data is correctly converted into an Media struct without errors.
// - Invalid URL: Ensures that an invalid URL triggers a validation error.
// - Missing Required Fields: Ensures that missing mandatory fields trigger a validation error.
// - Zero Value Fields: Ensures that empty field values are handled correctly and trigger a validation error.
func TestMediaConversions(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		expected *article.Media
		err      bool
	}{
		{
			name: "Valid Media",
			input: map[string]any{
				"id":          "123e4567-e89b-12d3-a456-426614174000",
				"url":         "https://example.com/media.jpg",
				"title":       "An example media",
				"description": "An example description",
				"width":       800,
				"height":      600,
				"size":        1024,
				"author":      "",
			},
			expected: &article.Media{
				ID:          "123e4567-e89b-12d3-a456-426614174000",
				URL:         "https://example.com/media.jpg",
				Title:       "An example media",
				Description: "An example description",
				Width:       800,
				Height:      600,
				Size:        1024,
			},
			err: false,
		},
		{
			name: "Invalid URL",
			input: map[string]any{
				"url":         "invalid-url",
				"title":       "An example media",
				"width":       800,
				"height":      600,
				"description": "An example description",
			},
			expected: nil,
			err:      true,
		},
		{
			name: "Missing Required Fields",
			input: map[string]any{
				"url":         "",
				"title":       "An example media",
				"width":       800,
				"height":      600,
				"description": "An example description",
			},
			expected: nil,
			err:      true,
		},
		{
			name: "Zero Value Fields",
			input: map[string]any{
				"url":         "",
				"title":       "",
				"width":       0,
				"height":      0,
				"description": "",
			},
			expected: &article.Media{
				URL:         "",
				Title:       "",
				Width:       0,
				Height:      0,
				Description: "",
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img, err := article.NewMediaFromMap(tt.input)
			if tt.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, img)
				assert.Equal(t, tt.input, img.Map())
			}
		})
	}
}

func TestMediaNormalize(t *testing.T) {
	img := &article.Media{
		URL:         "  " + gofakeit.URL() + "  ",
		Title:       "  " + gofakeit.Sentence(5) + "  ",
		Width:       gofakeit.Number(800, 1920),
		Height:      gofakeit.Number(600, 1080),
		Description: "  " + gofakeit.Sentence(10) + "  ",
	}

	img.Normalize()

	assert.NotEmpty(t, img.URL)
	assert.Equal(t, strings.TrimSpace(img.URL), img.URL)
	assert.Equal(t, strings.TrimSpace(img.Title), img.Title)
	assert.Equal(t, strings.TrimSpace(img.Description), img.Description)
}
