package article

import (
	"errors"
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
		ID:      uuid.New().String(),
		Tags:    NewTags(),
		Images:  NewImages(),
		Videos:  NewVideos(),
		Quotes:  NewQuotes(),
		Socials: NewSocials(),
	}
}

// Article represents a news article with various types of content.
// This structure provides a flexible and universal foundation for storing and working with various types of content,
// allowing for easy creation and modification of articles, as well as integration of media and social elements.
type Article struct {
	ID string `json:"id" validate:"required,uuid4,max=36"`
	// Genre of the article, e.g. news, opinion, review.
	Genre    string `json:"genre" validate:"max=500"`
	Category string `json:"category" validate:"max=255"`
	Author   string `json:"author" validate:"max=255"`
	// Title of the article.
	Title string `json:"title" validate:"required,max=255"`
	// Summary is a short description of the article.
	Summary string `json:"summary" validate:"max=500"`
	// Markup is the raw HTML or Markdown content of the article.
	Markup string `json:"markup" validate:"required,max=65000"`
	// Text plain text content of the article.
	Text string `json:"text" validate:"required,max=65000"`
	// SourceURL is the URL of the article.
	SourceURL string `json:"source_url" validate:"omitempty,url,max=4096"`
	// SourceName is the web resource name of the source, e.g. Washington Post.
	SourceName string `json:"source_name" validate:"max=255"`
	Language   string `json:"language" validate:"max=255"`
	// Published is the date and time when the article was published.
	Published time.Time `json:"published" validate:"required"`
	Modified  time.Time `json:"modified"`
	Images    *Images   `json:"images"`
	Videos    *Videos   `json:"videos"`
	Quotes    *Quotes   `json:"quotes"`
	Tags      *Tags     `json:"tags"`
	Socials   *Socials  `json:"socials"`
}

// Normalize validates the Article and its nested structures, logs any validation errors, and clears invalid fields.
func (a *Article) Normalize() error {

	a.ID = TrimToMaxLen(a.ID, 36)
	a.Genre = TrimToMaxLen(a.Genre, 500)
	a.Category = TrimToMaxLen(a.Category, 255)
	a.Author = TrimToMaxLen(a.Author, 255)
	a.Title = TrimToMaxLen(a.Title, 255)
	a.Summary = TrimToMaxLen(a.Summary, 500)
	a.Markup = TrimToMaxLen(a.Markup, 65000)
	a.Text = TrimToMaxLen(a.Text, 65000)
	a.SourceURL = TrimToMaxLen(a.SourceURL, 4096)
	a.SourceName = TrimToMaxLen(a.SourceName, 255)
	a.Language = TrimToMaxLen(a.Language, 255)

	// fallback to English if language is not set
	if len(a.Language) == 0 {
		a.Language = "en"
	}

	// fallback date to now if not set
	if a.Published.IsZero() {
		a.Published = time.Now()
	}

	err := validate.Struct(a)
	if err != nil {

		var invalids validator.ValidationErrors
		if !errors.As(err, &invalids) {
			return err
		}

		for _, invalid := range invalids {

			slog.Debug("Validation error", slog.String("field", invalid.Namespace()), slog.String("error", invalid.Tag()))

			if invalid.Tag() == "required" {
				return err
			}

			// Clear invalid fields
			switch invalid.Namespace() {
			case "Article.ID":
				a.ID = ""
			case "Article.Title":
				a.Title = ""
			case "Article.Summary":
				a.Summary = ""
			case "Article.Markup":
				a.Markup = ""
			case "Article.Text":
				a.Text = ""
			case "Article.Genre":
				a.Genre = ""
			case "Article.Published":
				a.Published = time.Time{}
			case "Article.Modified":
				a.Modified = time.Time{}
			case "Article.SourceURL":
				a.SourceURL = ""
			case "Article.Language":
				a.Language = ""
			case "Article.Category":
				a.Category = ""
			case "Article.SourceName":
				a.SourceName = ""
			}
		}
	}

	// Normalize nested structures
	a.Images.Normalize()
	a.Videos.Normalize()
	a.Quotes.Normalize()
	a.Socials.Normalize()

	return nil
}

func (a *Article) validate() error {

	err := validate.Struct(a)
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
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
		"markup":                          a.Markup,
		"text":                            a.Text,
		"genre":                           a.Genre,
		"images":                          images,
		"videos":                          videos,
		"quotes":                          quotes,
		"published":                       a.Published,
		"modified":                        a.Modified,
		"tags":                            a.Tags.Slice(),
		"source_url":                      a.SourceURL,
		"language":                        a.Language,
		"category":                        a.Category,
		"source_name":                     a.SourceName,
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
		ID:         StringFromMap(m, "id"),
		Title:      StringFromMap(m, "title"),
		Summary:    StringFromMap(m, "summary"),
		Markup:     StringFromMap(m, "markup"),
		Text:       StringFromMap(m, "text"),
		Genre:      StringFromMap(m, "genre"),
		Images:     images,
		Videos:     videos,
		Quotes:     quotes,
		Published:  publishDate,
		Modified:   modifiedDate,
		Tags:       NewTags(GetStringSlice(m, "tags")...),
		SourceURL:  StringFromMap(m, "source_url"),
		Language:   StringFromMap(m, "language"),
		Category:   StringFromMap(m, "category"),
		SourceName: StringFromMap(m, "source_name"),
		Socials:    social,
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
