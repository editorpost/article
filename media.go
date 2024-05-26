package article

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log/slog"
)

// Media represents a media in the article.
type Media struct {

	// ID is the unique identifier of the media.
	// It is stable enough to be used as a key in a storage system.
	ID string `json:"id" validate:"required,max=36"`

	// Author is the author of the media.
	// This field is optional and could be between 1 and 255 characters long.
	Author string `json:"author,omitempty" validate:"max=255"`

	// URL is the URL of the media.
	// This field is required and should be a valid URL.
	URL string `json:"url" validate:"url,max=4096"`

	// Title is the alternative text for the media.
	// This field is required and should be between 1 and 255 characters long.
	Title string `json:"title" validate:"max=255"`

	// Description is the description for the media.
	// This field is optional.
	Description string `json:"description,omitempty" validate:"max=500"`

	// Width is the width of the media in pixels.
	// This field is optional.
	Width int `json:"width" validate:"min=0"`

	// Height is the height of the media in pixels.
	// This field is optional.
	Height int `json:"height,omitempty" validate:"min=0"`

	// Size is the size of the media in bytes.
	// This field is optional. Populate by loaders.
	Size int `json:"size,omitempty" validate:"min=0"`
}

// NewMedia creates a new Media with a random UUID.
//
//goland:noinspection GoUnusedExportedFunction
func NewMedia(url string) *Media {
	return &Media{
		ID:  uuid.New().String(),
		URL: url,
	}
}

// Normalize validates and trims the fields of the Media.
func (i *Media) Normalize() {

	if i.ID == "" {
		i.ID = uuid.New().String()
	}

	i.ID = TrimToMaxLen(i.ID, 36)
	i.Author = TrimToMaxLen(i.Author, 255)
	i.URL = TrimToMaxLen(i.URL, 4096)
	i.Title = TrimToMaxLen(i.Title, 255)
	i.Description = TrimToMaxLen(i.Description, 500)

	err := validate.Struct(i)
	if err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			slog.Debug("Validation error in Media", slog.String("field", fieldErr.Namespace()), slog.String("error", fieldErr.Tag()))
			*i = Media{}
		}
	}
}

// Map converts the Media struct to a map[string]any.
func (i *Media) Map() map[string]any {
	return map[string]any{
		"id":          i.ID,
		"author":      i.Author,
		"url":         i.URL,
		"title":       i.Title,
		"description": i.Description,
		"width":       i.Width,
		"height":      i.Height,
		"size":        i.Size,
	}
}

// NewMediaFromMap creates a Media from a map[string]any, validates it, and returns a pointer to the Media or an error.
func NewMediaFromMap(m map[string]any) (*Media, error) {
	img := &Media{
		ID:          StringFromMap(m, "id"),
		Author:      StringFromMap(m, "author"),
		URL:         StringFromMap(m, "url"),
		Title:       StringFromMap(m, "title"),
		Description: StringFromMap(m, "description"),
		Width:       IntFromMap(m, "width"),
		Height:      IntFromMap(m, "height"),
		Size:        IntFromMap(m, "size"),
	}

	err := validate.Struct(img)
	if err != nil {
		return nil, err
	}

	return img, nil
}
