package article

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log/slog"
	"strings"
	"time"
	"unicode/utf8"
)

// init go validator instance
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// NewArticle creates a new Article with the provided data and returns a pointer to the Article.
func NewArticle() *Article {
	return &Article{
		ID:       uuid.New().String(),
		Language: "en",
		Tags:     NewTags(),
		Images:   NewImages(),
		Videos:   NewVideos(),
		Quotes:   NewQuotes(),
		Socials:  NewSocials(),
	}
}

// Article represents a news article with various types of content.
// This structure provides a flexible and universal foundation for storing and working with various types of content,
// allowing for easy creation and modification of articles, as well as integration of media and social elements.
type Article struct {
	ID          string    `json:"article__id" validate:"required,uuid4,max=36"`
	Title       string    `json:"article__title" validate:"required,max=255"`
	Byline      string    `json:"article__byline" validate:"max=255"`
	HTML        string    `json:"article__html" validate:"required,max=65000"`
	TextContent string    `json:"article__text" validate:"required,max=65000"`
	Excerpt     string    `json:"article__excerpt" validate:"max=500"`
	Published   time.Time `json:"article__published" validate:"required"`
	Modified    time.Time `json:"article__modified"`
	Images      *Images   `json:"article__images"`
	Videos      *Videos   `json:"article__videos"`
	Quotes      *Quotes   `json:"article__quotes"`
	Tags        *Tags     `json:"article__tags"`
	Socials     *Socials  `json:"article__socials"`
	Source      string    `json:"article__source" validate:"omitempty,url,max=4096"`
	Language    string    `json:"article__language" validate:"max=255"`
	Category    string    `json:"article__category" validate:"max=255"`
	SiteName    string    `json:"article__site" validate:"max=255"`
}

// Normalize validates the Article and its nested structures, logs any validation errors, and clears invalid fields.
func (a *Article) Normalize() {

	a.ID = TrimToMaxLen(a.ID, 36)
	a.Title = TrimToMaxLen(a.Title, 255)
	a.Byline = TrimToMaxLen(a.Byline, 255)
	a.HTML = TrimToMaxLen(a.HTML, 65000)
	a.TextContent = TrimToMaxLen(a.TextContent, 65000)
	a.Excerpt = TrimToMaxLen(a.Excerpt, 500)
	a.Source = TrimToMaxLen(a.Source, 4096)
	a.Language = TrimToMaxLen(a.Language, 255)
	a.Category = TrimToMaxLen(a.Category, 255)
	a.SiteName = TrimToMaxLen(a.SiteName, 255)

	err := validate.Struct(a)
	if err != nil {
		for _, invalid := range err.(validator.ValidationErrors) {
			slog.Debug("Validation error", slog.String("field", invalid.Namespace()), slog.String("error", invalid.Tag()))

			// Clear invalid fields
			switch invalid.Namespace() {
			case "Article.ID":
				a.ID = ""
			case "Article.Title":
				a.Title = ""
			case "Article.Byline":
				a.Byline = ""
			case "Article.HTML":
				a.HTML = ""
			case "Article.TextContent":
				a.TextContent = ""
			case "Article.Excerpt":
				a.Excerpt = ""
			case "Article.Published":
				a.Published = time.Time{}
			case "Article.Modified":
				a.Modified = time.Time{}
			case "Article.Source":
				a.Source = ""
			case "Article.Language":
				a.Language = ""
			case "Article.Category":
				a.Category = ""
			case "Article.SiteName":
				a.SiteName = ""
			}
		}
	}

	// Normalize nested structures
	a.Images.Normalize()
	a.Videos.Normalize()
	a.Quotes.Normalize()
	a.Socials.Normalize()
}

// Map converts the Article struct to a map[string]any, including nested structures.
func (a *Article) Map() map[string]any {

	images := make([]map[string]any, a.Images.Len())
	for i, image := range a.Images.Slice() {
		images[i] = image.Map()
	}

	videos := make([]map[string]any, a.Videos.Len())
	for i, video := range a.Videos.Slice() {
		videos[i] = video.Map()
	}

	quotes := make([]map[string]any, a.Quotes.Len())
	for i, quote := range a.Quotes.Slice() {
		quotes[i] = quote.Map()
	}

	socialProfiles := make([]map[string]any, a.Socials.Len())
	for i, profile := range a.Socials.Slice() {
		socialProfiles[i] = profile.Map()
	}

	return map[string]any{
		"article__id":                     a.ID,
		"article__title":                  a.Title,
		"article__byline":                 a.Byline,
		"article__html":                   a.HTML,
		"article__text":                   a.TextContent,
		"article__excerpt":                a.Excerpt,
		"article__images":                 images,
		"article__videos":                 videos,
		"article__quotes":                 quotes,
		"article__published":              a.Published,
		"article__modified":               a.Modified,
		"article__tags":                   a.Tags.Slice(),
		"article__source":                 a.Source,
		"article__language":               a.Language,
		"article__category":               a.Category,
		"article__site":                   a.SiteName,
		"article__author_social_profiles": socialProfiles,
	}
}

// NewArticleFromMap creates an Article from a map[string]any, validates it, and returns a pointer to the Article or an error.
func NewArticleFromMap(m map[string]any) (*Article, error) {

	images := NewImages()
	if imgMaps, ok := m["article__images"].([]map[string]any); ok {
		for _, imgMap := range imgMaps {
			if img, err := NewImageFromMap(imgMap); err == nil {
				images.Add(img)
			}
		}
	}

	videos := NewVideos()
	if vidMaps, ok := m["article__videos"].([]map[string]any); ok {
		for _, vidMap := range vidMaps {
			if vid, err := NewVideoFromMap(vidMap); err == nil {
				videos.Add(vid)
			}
		}
	}

	quotes := NewQuotes()
	if quoteMaps, ok := m["article__quotes"].([]map[string]any); ok {
		for _, quoteMap := range quoteMaps {
			if quote, err := NewQuoteFromMap(quoteMap); err == nil {
				quotes.Add(quote)
			}
		}
	}

	social := NewSocials()
	if profileMaps, ok := m["article__author_social_profiles"].([]map[string]any); ok {
		for _, profileMap := range profileMaps {
			if profile, err := NewSocialProfileFromMap(profileMap); err == nil {
				social.Add(profile)
			}
		}
	}

	publishDate, _ := m["article__published"].(time.Time)
	modifiedDate, _ := m["article__modified"].(time.Time)

	article := &Article{
		ID:          StringFromMap(m, "article__id"),
		Title:       StringFromMap(m, "article__title"),
		Byline:      StringFromMap(m, "article__byline"),
		HTML:        StringFromMap(m, "article__html"),
		TextContent: StringFromMap(m, "article__text"),
		Excerpt:     StringFromMap(m, "article__excerpt"),
		Images:      images,
		Videos:      videos,
		Quotes:      quotes,
		Published:   publishDate,
		Modified:    modifiedDate,
		Tags:        NewTags(GetStringSlice(m, "article__tags")...),
		Source:      StringFromMap(m, "article__source"),
		Language:    StringFromMap(m, "article__language"),
		Category:    StringFromMap(m, "article__category"),
		SiteName:    StringFromMap(m, "article__site"),
		Socials:     social,
	}

	err := validate.Struct(article)
	if err != nil {
		return nil, err
	}

	return article, nil
}

// GetStringSlice safely extracts a slice of strings from the map or returns a zero value.
func GetStringSlice(m map[string]any, key string) []string {
	if value, ok := m[key]; ok {
		if slice, ok := value.([]string); ok {
			return slice
		}
	}
	return []string{}
}

// TrimToMaxLen trims the input string to the specified maximum length, ensuring that it doesn't exceed the length in runes.
func TrimToMaxLen(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	if utf8.RuneCountInString(s) > maxLen {
		runeStr := []rune(s)
		return string(runeStr[:maxLen])
	}
	return s
}
