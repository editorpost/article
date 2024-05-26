package article

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log/slog"
)

// SocialProfile represents a social media profile of an author.
type SocialProfile struct {

	// ID is the unique identifier of the quote.
	// It is stable enough to be used as a key in a storage system.
	ID string `json:"id" validate:"required,max=36"`

	// Platform is the platform of the social profile (e.g., Twitter, Facebook).
	// This field is required and should be between 1 and 50 characters long.
	Platform string `json:"platform" validate:"max=255"`

	// URL is the URL of the social profile.
	// This field is required and should be a valid URL.
	URL string `json:"url" validate:"required,url,max=4096"`
}

func NewSocialProfile(platform, url string) *SocialProfile {
	return &SocialProfile{
		ID:       uuid.New().String(),
		Platform: platform,
		URL:      url,
	}
}

// Normalize validates and trims the fields of the SocialProfile.
func (s *SocialProfile) Normalize() {

	if s.ID == "" {
		s.ID = uuid.New().String()
	}

	s.ID = TrimToMaxLen(s.ID, 36)
	s.Platform = TrimToMaxLen(s.Platform, 255)
	s.URL = TrimToMaxLen(s.URL, 4096)

	err := validate.Struct(s)
	if err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			slog.Debug("Validation error in SocialProfile", slog.String("field", fieldErr.Namespace()), slog.String("error", fieldErr.Tag()))
			*s = SocialProfile{}
		}
	}
}

// Map converts the SocialProfile struct to a map[string]any.
func (s *SocialProfile) Map() map[string]any {
	return map[string]any{
		"id":       s.ID,
		"platform": s.Platform,
		"url":      s.URL,
	}
}

// NewSocialProfileFromMap creates a SocialProfile from a map[string]any, validates it, and returns a pointer to the SocialProfile or an error.
func NewSocialProfileFromMap(m map[string]any) (*SocialProfile, error) {
	profile := &SocialProfile{
		ID:       StringFromMap(m, "id"),
		Platform: StringFromMap(m, "platform"),
		URL:      StringFromMap(m, "url"),
	}

	err := validate.Struct(profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
