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

	// URL is the URL of the video.
	// This field is required and should be a valid URL.
	URL string `json:"url" validate:"required,url,max=4096"`

	// EmbedCode is the embed code for the video.
	// This field is optional.
	EmbedCode string `json:"embed_code,omitempty" validate:"max=65000"`

	// Caption is the caption for the video.
	// This field is optional.
	Caption string `json:"caption,omitempty" validate:"max=500"`
}

// NewVideo creates a new Video with a random UUID.
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
	v.EmbedCode = TrimToMaxLen(v.EmbedCode, 65000)
	v.Caption = TrimToMaxLen(v.Caption, 500)

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
		"id":         v.ID,
		"url":        v.URL,
		"embed_code": v.EmbedCode,
		"caption":    v.Caption,
	}
}

// NewVideoFromMap creates a Video from a map[string]any, validates it, and returns a pointer to the Video or an error.
func NewVideoFromMap(m map[string]any) (*Video, error) {
	video := &Video{
		ID:        StringFromMap(m, "id"),
		URL:       StringFromMap(m, "url"),
		EmbedCode: StringFromMap(m, "embed_code"),
		Caption:   StringFromMap(m, "caption"),
	}

	err := validate.Struct(video)
	if err != nil {
		return nil, err
	}

	return video, nil
}
