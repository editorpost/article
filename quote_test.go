package article_test

import (
	"github.com/brianvoe/gofakeit/v6"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/editorpost/article"
)

func TestQuoteNormalize(t *testing.T) {
	q := &article.Quote{
		Text:      "  " + gofakeit.Sentence(15) + "  ",
		Author:    "  " + gofakeit.Name() + "  ",
		SourceURL: "  " + gofakeit.URL() + "  ",
		Platform:  "  " + gofakeit.Word() + "  ",
	}

	q.Normalize()

	assert.NotEmpty(t, q.Text)
	assert.Equal(t, strings.TrimSpace(q.Text), q.Text)
	assert.Equal(t, strings.TrimSpace(q.Author), q.Author)
	assert.Equal(t, strings.TrimSpace(q.SourceURL), q.SourceURL)
	assert.Equal(t, strings.TrimSpace(q.Platform), q.Platform)
}

// TestQuoteConversions is a table-driven test for the Quote struct.
// It verifies the conversion of map data to Quote struct, the validation process, and handling of zero-value fields.
//
// Explanation of test cases:
// - Valid Quote: Ensures that valid data is correctly converted into a Quote struct without errors.
// - Invalid SourceURL URL: Ensures that an invalid source URL triggers a validation error.
// - Missing Required Fields: Ensures that missing mandatory fields trigger a validation error. Specifically, 'text', 'author', 'source', and 'platform' fields are required.
// - Zero Value Fields: Ensures that empty field values are handled correctly and trigger a validation error.
func TestQuoteConversions(t *testing.T) {
	tests := []struct {
		name          string
		inputMap      map[string]any
		expectedQuote *article.Quote
		expectError   bool
	}{
		{
			name: "Valid Quote",
			inputMap: map[string]any{
				"id":         "123e4567-e89b-12d3-a456-426614174000",
				"text":       "This is a quote",
				"author":     "John Doe",
				"source_url": "https://example.com",
				"platform":   "Twitter",
			},
			expectedQuote: &article.Quote{
				ID:        "123e4567-e89b-12d3-a456-426614174000",
				Text:      "This is a quote",
				Author:    "John Doe",
				SourceURL: "https://example.com",
				Platform:  "Twitter",
			},
			expectError: false,
		},
		{
			name: "Invalid SourceURL URL",
			inputMap: map[string]any{
				"text":       "This is a quote",
				"author":     "John Doe",
				"source_url": "invalid-url",
				"platform":   "Twitter",
			},
			expectedQuote: nil,
			expectError:   true,
		},
		{
			name: "Missing Required Fields",
			inputMap: map[string]any{
				"text":   "This is a quote",
				"author": "John Doe",
			},
			expectedQuote: nil,
			expectError:   true,
		},
		{
			name: "Zero Value Fields",
			inputMap: map[string]any{
				"text":       "",
				"author":     "",
				"source_url": "",
				"platform":   "",
			},
			expectedQuote: &article.Quote{
				Text:      "",
				Author:    "",
				SourceURL: "",
				Platform:  "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quote, err := article.NewQuoteFromMap(tt.inputMap)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedQuote, quote)
				assert.Equal(t, tt.inputMap, quote.Map())
			}
		})
	}
}
