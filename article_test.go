package article_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/editorpost/article"

	"github.com/samber/lo"
)

func init() {
	gofakeit.Seed(0)
}

func TestMinimalInvariantArticle(t *testing.T) {

	expected := article.NewArticle()

	// required fields
	expected.Title = gofakeit.Sentence(3)
	expected.Markup = gofakeit.Paragraph(1, 5, 10, " ")
	expected.Text = gofakeit.Paragraph(1, 5, 10, " ")
	expected.Published = time.Now()

	got, err := article.NewArticleFromMap(expected.Map())
	require.NoError(t, err)

	assert.Equal(t, expected, got)
}

func TestFullInvariantArticle(t *testing.T) {

	expected := article.NewArticle()

	// Required fields
	expected.Title = gofakeit.Sentence(3)
	expected.Markup = gofakeit.Paragraph(1, 5, 10, " ")
	expected.Text = gofakeit.Paragraph(1, 5, 10, " ")
	expected.Published = time.Now()
	expected.Modified = time.Now()

	// Optional fields
	expected.Summary = gofakeit.Name()
	expected.Genre = gofakeit.Sentence(10)

	expected.Images = article.NewImages(&article.Image{
		ID:     gofakeit.UUID(),
		URL:    gofakeit.URL(),
		Alt:    gofakeit.Sentence(5),
		Width:  gofakeit.Number(800, 1920),
		Height: gofakeit.Number(600, 1080),
		Title:  gofakeit.Sentence(10),
	})

	expected.Videos = article.NewVideos(&article.Video{
		ID:    gofakeit.UUID(),
		URL:   gofakeit.URL(),
		Embed: "<iframe src='" + gofakeit.URL() + "'></iframe>",
		Title: gofakeit.Sentence(10),
	})

	expected.Quotes = article.NewQuotes(&article.Quote{
		ID:        gofakeit.UUID(),
		Text:      gofakeit.Sentence(15),
		Author:    gofakeit.Name(),
		SourceURL: gofakeit.URL(),
		Platform:  "Twitter",
	})

	expected.Tags = article.NewTags("travel", "Phuket", "Thailand")
	expected.SourceURL = gofakeit.URL()
	expected.Language = "en"
	expected.Category = "Travel"
	expected.SourceName = "Example Travel Blog"

	expected.Socials = article.NewSocials(&article.Social{
		Platform: "Twitter",
		URL:      gofakeit.URL(),
	})

	got, err := article.NewArticleFromMap(expected.Map())
	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestInvalidNestedStructureArticle(t *testing.T) {
	expected := article.NewArticle()

	// Required fields
	expected.Title = gofakeit.Sentence(3)
	expected.Markup = gofakeit.Paragraph(1, 5, 10, " ")
	expected.Text = gofakeit.Paragraph(1, 5, 10, " ")
	expected.Published = time.Now()
	expected.Modified = time.Now()

	// Optional fields with invalid nested structure
	expected.Summary = gofakeit.Name()
	expected.Genre = gofakeit.Sentence(10)

	expected.Images = article.NewImages(&article.Image{
		URL:    "invalid-url",
		Alt:    gofakeit.Sentence(5),
		Width:  gofakeit.Number(800, 1920),
		Height: gofakeit.Number(600, 1080),
		Title:  gofakeit.Sentence(10),
	})

	expected.Tags = article.NewTags("travel", "Phuket", "Thailand")
	expected.SourceURL = gofakeit.URL()
	expected.Language = "en"
	expected.Category = "Travel"
	expected.SourceName = "Example Travel Blog"

	// Convert expected Article to map and then back to Article to simulate input processing
	inputMap := expected.Map()

	// Expect the items to be nil due to invalid URL in the nested structure
	expected.Images = article.NewImages()

	got, err := article.NewArticleFromMap(inputMap)
	require.NoError(t, err)

	// To compare Published and Modified separately due to possible time differences
	assert.Equal(t, expected.ID, got.ID)
	assert.Equal(t, expected.Title, got.Title)
	assert.Equal(t, expected.Summary, got.Summary)
	assert.Equal(t, expected.Markup, got.Markup)
	assert.Equal(t, expected.Text, got.Text)
	assert.Equal(t, expected.Genre, got.Genre)
	assert.Equal(t, expected.Images, got.Images)
	assert.WithinDuration(t, expected.Published, got.Published, time.Second)
	assert.WithinDuration(t, expected.Modified, got.Modified, time.Second)
	assert.Equal(t, expected.Tags, got.Tags)
	assert.Equal(t, expected.SourceURL, got.SourceURL)
	assert.Equal(t, expected.Language, got.Language)
	assert.Equal(t, expected.Category, got.Category)
	assert.Equal(t, expected.SourceName, got.SourceName)
}

func TestMissingRequiredFieldsArticle(t *testing.T) {
	art, err := article.NewArticleFromMap(article.NewArticle().Map())
	require.Error(t, err)
	assert.Nil(t, art)
}

func TestArticleNormalize(t *testing.T) {
	expected := article.NewArticle()

	// Required fields
	expected.Title = gofakeit.Sentence(3)
	expected.Markup = gofakeit.Paragraph(1, 5, 10, " ")
	expected.Text = gofakeit.Paragraph(1, 5, 10, " ")
	expected.Published = time.Now()
	expected.Modified = time.Now()

	// Optional fields with some invalid data
	expected.Summary = gofakeit.Name()
	expected.Genre = gofakeit.Sentence(10)

	expected.Images = article.NewImages(&article.Image{
		URL:    "invalid-url",
		Alt:    gofakeit.Sentence(5),
		Width:  gofakeit.Number(800, 1920),
		Height: gofakeit.Number(600, 1080),
		Title:  gofakeit.Sentence(10),
	})

	expected.Videos = article.NewVideos(&article.Video{
		URL:   "invalid-url",
		Embed: "<iframe src='invalid-url'></iframe>",
		Title: gofakeit.Sentence(10),
	})

	expected.Quotes = article.NewQuotes(&article.Quote{
		Text:      "",
		Author:    gofakeit.Name(),
		SourceURL: "invalid-url",
		Platform:  "Twitter",
	})

	expected.Tags = article.NewTags("travel", "Phuket", "Thailand")
	expected.SourceURL = "invalid-url"
	expected.Language = "en"
	expected.Category = "Travel"
	expected.SourceName = "Example Travel Blog"

	expected.Socials = article.NewSocials(&article.Social{
		Platform: "Twitter",
		URL:      "invalid-url",
	})

	expected.Normalize()

	// Verify that invalid fields are cleared
	assert.Empty(t, expected.Images.Slice())
	assert.Empty(t, expected.Videos.Slice())
	assert.Empty(t, expected.Quotes.Slice())
	assert.Empty(t, "", expected.Socials.Slice())
	assert.Equal(t, "", expected.SourceURL)
}

func TestArticleNormalizeFieldClearing(t *testing.T) {

	invalid := article.NewArticle()

	// Set required fields with valid data
	invalid.Title = gofakeit.Sentence(3)
	invalid.Markup = gofakeit.Paragraph(1, 5, 10, " ")
	invalid.Text = gofakeit.Paragraph(1, 5, 10, " ")
	invalid.Published = time.Now()
	invalid.Modified = time.Now()

	// Set invalid data for optional fields
	invalid.ID = "invalid-uuid"
	invalid.Summary = gofakeit.Name()
	invalid.Genre = gofakeit.Sentence(10)
	invalid.SourceURL = "invalid-url"
	invalid.Language = "english" // should be a valid ISO 639-1 language code
	invalid.Category = gofakeit.Sentence(2)
	invalid.SourceName = gofakeit.Sentence(2)

	valid := *invalid
	(&valid).Normalize()

	// Verify that invalid fields are cleared
	assert.Equal(t, "", valid.ID)
	assert.Equal(t, invalid.Summary, valid.Summary) // should not be cleared since it's not required
	assert.Equal(t, invalid.Genre, valid.Genre)     // should not be cleared since it's not required
	assert.Equal(t, "", valid.SourceURL)
	assert.Equal(t, "english", valid.Language)
	assert.Equal(t, invalid.Category, valid.Category)
	assert.Equal(t, invalid.SourceName, valid.SourceName)
}

// TestGetStringSlice tests the GetStringSlice function in case of empty map and missing key:
func TestGetStringSlice(t *testing.T) {
	m := map[string]interface{}{}
	key := "key"
	assert.Equal(t, []string{}, article.GetStringSlice(m, key))
}

func TestTrimToMaxLen(t *testing.T) {
	s := "This is a test string with more than twenty characters."
	trimmed := article.TrimToMaxLen(s, 21)
	assert.Equal(t, "This is a test string", trimmed)

	s = "Short string"
	trimmed = article.TrimToMaxLen(s, 20)
	assert.Equal(t, s, trimmed)
}

func TestUnmarshal(t *testing.T) {

	js := `{
		  "id": "123e4567-e89b-12d3-a456-426614174000",
		  "title": "The Rise of AI",
		  "summary": "By John Doe",
		  "markup": "<p>Artificial Intelligence is transforming the world.</p>",
		  "text": "Artificial Intelligence is transforming the world.",
		  "genre": "An overview of how AI is changing various industries.",
		  "published": "2024-05-27T10:00:00Z",
		  "modified": "2024-05-28T12:00:00Z",
		  "images": [
			  {
				"id": "img-001",
				"url": "https://example.com/image1.jpg",
				"alt": "AI Illustration",
				"width": 800,
				"height": 600,
				"title": "An illustration representing AI."
			  }
          ],
		  "videos": [
			  {
				"id": "vid-001",
				"url": "https://example.com/video1.mp4",
				"embed": "<iframe src='https://example.com/video1.mp4'></iframe>",
				"title": "A video explaining AI."
			  }
		  ],
		  "quotes": [
			  {
				"id": "quote-001",
				"text": "AI is the future of technology.",
				"author": "Jane Smith",
				"source_url": "https://twitter.com/janesmith/status/123",
				"platform": "Twitter"
			  }
		  ],
		  "tags": ["AI", "Technology", "Future"],
		  "source_url": "https://example.com",
		  "language": "en",
		  "category": "Technology",
		  "source_name": "Tech News",
		  "socials": [
			  {
				"id": "sp-001",
				"platform": "Twitter",
				"url": "https://twitter.com/johndoe"
			  }
		  ]
	}`

	// use encoding/json to unmarshal the JSON string into Article

	art := article.Article{}
	require.NoError(t, json.Unmarshal([]byte(js), &art))

	// check the values of the Article fields
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", art.ID)
	assert.Equal(t, "The Rise of AI", art.Title)
	assert.Equal(t, "By John Doe", art.Summary)
	assert.Equal(t, "<p>Artificial Intelligence is transforming the world.</p>", art.Markup)
	assert.Equal(t, "Artificial Intelligence is transforming the world.", art.Text)
	assert.Equal(t, "An overview of how AI is changing various industries.", art.Genre)
	assert.Equal(t, "2024-05-27T10:00:00Z", art.Published.Format(time.RFC3339))
	assert.Equal(t, "2024-05-28T12:00:00Z", art.Modified.Format(time.RFC3339))

	// check the values of the nested structures
	assert.Equal(t, 1, art.Images.Len())
	assert.Equal(t, 1, art.Videos.Len())
	assert.Equal(t, 1, art.Quotes.Len())
	assert.Equal(t, 1, art.Socials.Len())
	assert.Equal(t, 3, art.Tags.Len())
}

func TestArticle_ReplaceURLs_Empty(t *testing.T) {

	a := article.NewArticle()
	count := 2
	WithImages(a, GenerateURLs("example.com", count)...)

	// empty map
	failed := a.ReplaceURLs(map[string]string{})
	// all images ids as failed to replace
	assert.Equal(t, count, len(failed))
	// no images removed
	assert.Equal(t, count, a.Images.Len())

	// remove failed to replace
	failed = a.ReplaceOrRemoveURLs(map[string]string{})
	assert.Equal(t, count, len(failed))
	// all images removed
	assert.Zero(t, a.Images.Len())
}

func TestArticle_ReplaceURLs_All(t *testing.T) {

	count := 5
	replace := ReplacementURLs(count)
	src := lo.Keys(replace)

	a := article.NewArticle()
	WithImages(a, src...)

	failed := a.ReplaceURLs(replace)
	// no images ids as failed to replace
	assert.Zero(t, 0, len(failed))
	// no images removed
	assert.Equal(t, count, a.Images.Len())
}

func TestArticle_ReplaceOrRemoveURLs_All(t *testing.T) {

	count := 5
	replace := ReplacementURLs(count)
	src := lo.Keys(replace)

	a := article.NewArticle()
	WithImages(a, src...)

	// remove all images
	failed := a.ReplaceOrRemoveURLs(replace)
	assert.Zero(t, 0, len(failed))
	// no images removed
	assert.Equal(t, count, a.Images.Len())
}

func TestArticle_ReplaceURLs_Partial(t *testing.T) {

	count := 5
	replace := ReplacementURLs(count)
	src := lo.Keys(replace)

	a := article.NewArticle()
	WithImages(a, src...)

	// remove some urls
	for i := 0; i < count/2; i++ {
		delete(replace, src[i])
	}

	failed := a.ReplaceURLs(replace)
	// some images ids as failed to replace
	assert.Equal(t, count/2, len(failed))
	// no images removed
	assert.Equal(t, count, a.Images.Len())
}

func TestArticle_ReplaceOrRemoveURLs_Partial(t *testing.T) {

	count := 5
	replace := ReplacementURLs(count)
	src := lo.Keys(replace)

	a := article.NewArticle()
	WithImages(a, src...)

	// remove some urls
	for i := 0; i < count/2; i++ {
		delete(replace, src[i])
	}

	failed := a.ReplaceOrRemoveURLs(replace)
	// some images ids as failed to replace
	assert.Equal(t, count/2, len(failed))
	// failed images removed
	assert.Equal(t, count-(len(failed)), a.Images.Len())
}

func WithImages(a *article.Article, urls ...string) {
	for _, u := range urls {
		a.Images.Add(&article.Image{
			ID:     gofakeit.UUID(),
			URL:    u,
			Alt:    gofakeit.Sentence(5),
			Width:  gofakeit.Number(800, 1920),
			Height: gofakeit.Number(600, 1080),
			Title:  gofakeit.Sentence(10),
		})
	}
}

func GenerateURLs(domain string, count int) []string {
	urls := make([]string, count)
	for i := 0; i < count; i++ {
		urls[i] = fmt.Sprintf("https://%s/images/%s.jpg", domain, gofakeit.UUID())
	}
	return urls
}

func ReplacementURLs(count int) map[string]string {

	old := GenerateURLs("old.com", count)
	news := GenerateURLs("new.com", count)

	replace := make(map[string]string)
	for i, url := range old {
		replace[url] = news[i]
	}

	return replace
}
