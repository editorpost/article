package article

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log/slog"
)

// Image represents an image in the article.
type Image struct {

	// ID is the unique identifier of the image.
	// It is stable enough to be used as a key in a storage system.
	ID string `json:"id" validate:"required,max=36"`

	// URL is the URL of the image.
	// This field is required and should be a valid URL.
	URL string `json:"url" validate:"url,max=4096"`

	// AltText is the alternative text for the image.
	// This field is required and should be between 1 and 255 characters long.
	AltText string `json:"alt_text" validate:"max=255"`

	// Width is the width of the image in pixels.
	// This field is optional.
	Width int `json:"width" validate:"min=0"`

	// Height is the height of the image in pixels.
	// This field is optional.
	Height int `json:"height,omitempty" validate:"min=0"`

	// Caption is the caption for the image.
	// This field is optional.
	Caption string `json:"caption,omitempty" validate:"max=500"`
}

// NewImage creates a new Image with a random UUID.
func NewImage(url string) *Image {
	return &Image{
		ID:  uuid.New().String(),
		URL: url,
	}
}

// Normalize validates and trims the fields of the Image.
func (i *Image) Normalize() {

	if i.ID == "" {
		i.ID = uuid.New().String()
	}

	i.ID = TrimToMaxLen(i.ID, 36)
	i.URL = TrimToMaxLen(i.URL, 4096)
	i.AltText = TrimToMaxLen(i.AltText, 255)
	i.Caption = TrimToMaxLen(i.Caption, 500)

	err := validate.Struct(i)
	if err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			slog.Debug("Validation error in Image", slog.String("field", fieldErr.Namespace()), slog.String("error", fieldErr.Tag()))
			*i = Image{}
		}
	}
}

// Map converts the Image struct to a map[string]any.
func (i *Image) Map() map[string]any {
	return map[string]any{
		"id":       i.ID,
		"url":      i.URL,
		"alt_text": i.AltText,
		"width":    i.Width,
		"height":   i.Height,
		"caption":  i.Caption,
	}
}

// NewImageFromMap creates an Image from a map[string]any, validates it, and returns a pointer to the Image or an error.
func NewImageFromMap(m map[string]any) (*Image, error) {
	img := &Image{
		ID:      StringFromMap(m, "id"),
		URL:     StringFromMap(m, "url"),
		AltText: StringFromMap(m, "alt_text"),
		Width:   IntFromMap(m, "width"),
		Height:  IntFromMap(m, "height"),
		Caption: StringFromMap(m, "caption"),
	}

	err := validate.Struct(img)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// StringFromMap safely extracts a string from the map or returns a zero value.
func StringFromMap(m map[string]any, key string) string {
	if value, exists := m[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// IntFromMap safely extracts an int from the map or returns a zero value.
func IntFromMap(m map[string]any, key string) int {
	if value, exists := m[key]; exists {
		if i, ok := value.(int); ok {
			return i
		}
	}
	return 0
}
