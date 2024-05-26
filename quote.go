package article

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log/slog"
)

// Quote represents a quote from social media in the article.
type Quote struct {

	// ID is the unique identifier of the quote.
	// It is stable enough to be used as a key in a storage system.
	ID string `json:"id" validate:"required,max=36"`

	// Text is the text of the quote.
	// This field is required.
	Text string `json:"text" validate:"required,max=65000"`

	// Author is the author of the quote.
	// This field is required and should be between 1 and 255 characters long.
	Author string `json:"author" validate:"max=255"`

	// Source is the source URL of the quote.
	// This field is required and should be a valid URL.
	Source string `json:"source" validate:"url,max=4096"`

	// Platform is the platform of the quote (e.g., Twitter, Facebook).
	// This field is required and should be between 1 and 50 characters long.
	Platform string `json:"platform" validate:"max=255"`
}

func NewQuote(text string) *Quote {
	return &Quote{
		ID:   uuid.New().String(),
		Text: text,
	}
}

// Normalize validates and trims the fields of the Quote.
func (q *Quote) Normalize() {

	if q.ID == "" {
		q.ID = uuid.New().String()
	}

	q.ID = TrimToMaxLen(q.ID, 36)
	q.Text = TrimToMaxLen(q.Text, 65000)
	q.Author = TrimToMaxLen(q.Author, 255)
	q.Source = TrimToMaxLen(q.Source, 4096)
	q.Platform = TrimToMaxLen(q.Platform, 255)

	err := validate.Struct(q)
	if err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			slog.Debug("Validation error in Quote", slog.String("field", fieldErr.Namespace()), slog.String("error", fieldErr.Tag()))
			*q = Quote{}
		}
	}
}

// Map converts the Quote struct to a map[string]any.
func (q *Quote) Map() map[string]any {
	return map[string]any{
		"id":       q.ID,
		"text":     q.Text,
		"author":   q.Author,
		"source":   q.Source,
		"platform": q.Platform,
	}
}

// NewQuoteFromMap creates a Quote from a map[string]any, validates it, and returns a pointer to the Quote or an error.
func NewQuoteFromMap(m map[string]any) (*Quote, error) {
	quote := &Quote{
		ID:       StringFromMap(m, "id"),
		Text:     StringFromMap(m, "text"),
		Author:   StringFromMap(m, "author"),
		Source:   StringFromMap(m, "source"),
		Platform: StringFromMap(m, "platform"),
	}

	err := validate.Struct(quote)
	if err != nil {
		return nil, err
	}

	return quote, nil
}
