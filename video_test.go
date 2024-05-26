package article_test

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/editorpost/article"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVideoNormalize(t *testing.T) {
	v := &article.Video{
		URL:   "  " + gofakeit.URL() + "  ",
		Embed: "  " + gofakeit.Sentence(100) + "  ",
		Title: "  " + gofakeit.Sentence(10) + "  ",
	}

	v.Normalize()

	assert.NotEmpty(t, v.URL)
	assert.Equal(t, strings.TrimSpace(v.URL), v.URL)
	assert.Equal(t, strings.TrimSpace(v.Embed), v.Embed)
	assert.Equal(t, strings.TrimSpace(v.Title), v.Title)
}

// TestVideoConversions is a table-driven test for the Video struct.
// It verifies the conversion of map data to Video struct, the validation process, and handling of zero-value fields.
//
// Explanation of test cases:
// - Valid Video: Ensures that valid data is correctly converted into a Video struct without errors.
// - Invalid URL: Ensures that an invalid URL triggers a validation error.
// - Missing Required Fields: Ensures that missing mandatory fields trigger a validation error. Specifically, the 'url' field is required.
// - Zero Value Fields: Ensures that empty field values are handled correctly and trigger a validation error.
func TestVideoConversions(t *testing.T) {
	tests := []struct {
		name          string
		inputMap      map[string]any
		expectedVideo *article.Video
		expectError   bool
	}{
		{
			name: "Valid Video",
			inputMap: map[string]any{
				"id":    "123e4567-e89b-12d3-a456-426614174000",
				"url":   "https://example.com/video.mp4",
				"embed": "<iframe src='https://example.com/video'></iframe>",
				"title": "An example title",
			},
			expectedVideo: &article.Video{
				ID:    "123e4567-e89b-12d3-a456-426614174000",
				URL:   "https://example.com/video.mp4",
				Embed: "<iframe src='https://example.com/video'></iframe>",
				Title: "An example title",
			},
			expectError: false,
		},
		{
			name: "Invalid URL",
			inputMap: map[string]any{
				"id":    "123e4567-e89b-12d3-a456-426614174000",
				"url":   "invalid-url",
				"embed": "<iframe src='https://example.com/video'></iframe>",
				"title": "An example title",
			},
			expectedVideo: nil,
			expectError:   true,
		},
		{
			name: "Missing Required Fields",
			inputMap: map[string]any{
				"title": "Some title, but no URL",
			},
			expectedVideo: nil,
			expectError:   true,
		},
		{
			name: "Zero Value Fields",
			inputMap: map[string]any{
				"url":   "",
				"embed": "",
				"title": "",
			},
			expectedVideo: &article.Video{
				URL:   "",
				Embed: "",
				Title: "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vid, err := article.NewVideoFromMap(tt.inputMap)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedVideo, vid)
				assert.Equal(t, tt.inputMap, vid.Map())
			}
		})
	}
}
