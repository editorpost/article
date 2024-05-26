package article

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log/slog"
)

// Video represents a video in the article.
type Video struct {

	// ID is the unique identifier of the video.
	// It is stable enough to be used as a key in a storage system.
	ID string `json:"id" validate:"required,max=36"`

	// Title is the title for the video.
	// This field is optional.
	Title string `json:"title,omitempty" validate:"max=500"`

	// URL is the URL of the video.
	// This field is required and should be a valid URL.
	URL string `json:"url" validate:"required,url,max=4096"`

	// Embed is the embed code for the video.
	// This field is optional.
	Embed string `json:"embed,omitempty" validate:"max=65000"`
}

// NewVideo creates a new Video with a random UUID.
//
//goland:noinspection GoUnusedExportedFunction
func NewVideo(url string) *Video {
	return &Video{
		ID:  uuid.New().String(),
		URL: url,
	}
}

// Normalize validates and trims the fields of the Video.
func (v *Video) Normalize() {

	if v.ID == "" {
		v.ID = uuid.New().String()
	}

	v.ID = TrimToMaxLen(v.ID, 36)
	v.URL = TrimToMaxLen(v.URL, 4096)
	v.Embed = TrimToMaxLen(v.Embed, 65000)
	v.Title = TrimToMaxLen(v.Title, 500)

	err := validate.Struct(v)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			slog.Debug("Validation error in Video", slog.String("field", err.Namespace()), slog.String("error", err.Tag()))
			*v = Video{}
		}
	}
}

// Map converts the Video struct to a map[string]any.
func (v *Video) Map() map[string]any {
	return map[string]any{
		"id":    v.ID,
		"url":   v.URL,
		"embed": v.Embed,
		"title": v.Title,
	}
}

// NewVideoFromMap creates a Video from a map[string]any, validates it, and returns a pointer to the Video or an error.
func NewVideoFromMap(m map[string]any) (*Video, error) {
	video := &Video{
		ID:    StringFromMap(m, "id"),
		URL:   StringFromMap(m, "url"),
		Embed: StringFromMap(m, "embed"),
		Title: StringFromMap(m, "title"),
	}

	err := validate.Struct(video)
	if err != nil {
		return nil, err
	}

	return video, nil
}
