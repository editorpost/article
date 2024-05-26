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
	ID        string    `json:"id" validate:"required,uuid4,max=36"`
	Title     string    `json:"title" validate:"required,max=255"`
	Summary   string    `json:"summary" validate:"max=255"`
	HTML      string    `json:"html" validate:"required,max=65000"`
	Text      string    `json:"text" validate:"required,max=65000"`
	Excerpt   string    `json:"excerpt" validate:"max=500"`
	Source    string    `json:"source" validate:"omitempty,url,max=4096"`
	Language  string    `json:"language" validate:"max=255"`
	Category  string    `json:"category" validate:"max=255"`
	SiteName  string    `json:"site" validate:"max=255"`
	Published time.Time `json:"published" validate:"required"`
	Modified  time.Time `json:"modified"`
	Images    *Images   `json:"images"`
	Videos    *Videos   `json:"videos"`
	Quotes    *Quotes   `json:"quotes"`
	Tags      *Tags     `json:"tags"`
	Socials   *Socials  `json:"socials"`
}

// Normalize validates the Article and its nested structures, logs any validation errors, and clears invalid fields.
func (a *Article) Normalize() {

	a.ID = TrimToMaxLen(a.ID, 36)
	a.Title = TrimToMaxLen(a.Title, 255)
	a.Summary = TrimToMaxLen(a.Summary, 255)
	a.HTML = TrimToMaxLen(a.HTML, 65000)
	a.Text = TrimToMaxLen(a.Text, 65000)
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
			case "Article.Summary":
				a.Summary = ""
			case "Article.HTML":
				a.HTML = ""
			case "Article.Text":
				a.Text = ""
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
		"id":                              a.ID,
		"title":                           a.Title,
		"summary":                         a.Summary,
		"html":                            a.HTML,
		"text":                            a.Text,
		"excerpt":                         a.Excerpt,
		"images":                          images,
		"videos":                          videos,
		"quotes":                          quotes,
		"published":                       a.Published,
		"modified":                        a.Modified,
		"tags":                            a.Tags.Slice(),
		"source":                          a.Source,
		"language":                        a.Language,
		"category":                        a.Category,
		"site":                            a.SiteName,
		"article__author_social_profiles": socialProfiles,
	}
}

// NewArticleFromMap creates an Article from a map[string]any, validates it, and returns a pointer to the Article or an error.
func NewArticleFromMap(m map[string]any) (*Article, error) {

	images := NewImages()
	if imgMaps, ok := m["images"].([]map[string]any); ok {
		for _, imgMap := range imgMaps {
			if img, err := NewImageFromMap(imgMap); err == nil {
				images.Add(img)
			}
		}
	}

	videos := NewVideos()
	if vidMaps, ok := m["videos"].([]map[string]any); ok {
		for _, vidMap := range vidMaps {
			if vid, err := NewVideoFromMap(vidMap); err == nil {
				videos.Add(vid)
			}
		}
	}

	quotes := NewQuotes()
	if quoteMaps, ok := m["quotes"].([]map[string]any); ok {
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

	publishDate, _ := m["published"].(time.Time)
	modifiedDate, _ := m["modified"].(time.Time)

	article := &Article{
		ID:        StringFromMap(m, "id"),
		Title:     StringFromMap(m, "title"),
		Summary:   StringFromMap(m, "summary"),
		HTML:      StringFromMap(m, "html"),
		Text:      StringFromMap(m, "text"),
		Excerpt:   StringFromMap(m, "excerpt"),
		Images:    images,
		Videos:    videos,
		Quotes:    quotes,
		Published: publishDate,
		Modified:  modifiedDate,
		Tags:      NewTags(GetStringSlice(m, "tags")...),
		Source:    StringFromMap(m, "source"),
		Language:  StringFromMap(m, "language"),
		Category:  StringFromMap(m, "category"),
		SiteName:  StringFromMap(m, "site"),
		Socials:   social,
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
