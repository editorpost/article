package article_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/editorpost/article"
)

func init() {
	gofakeit.Seed(0)
}

func TestMinimalInvariantArticle(t *testing.T) {

	expected := article.NewArticle()

	// required fields
	expected.Title = gofakeit.Sentence(3)
	expected.HTML = gofakeit.Paragraph(1, 5, 10, " ")
	expected.TextContent = gofakeit.Paragraph(1, 5, 10, " ")
	expected.Published = time.Now()

	got, err := article.NewArticleFromMap(expected.Map())
	require.NoError(t, err)

	assert.Equal(t, expected, got)
}

func TestFullInvariantArticle(t *testing.T) {

	expected := article.NewArticle()

	// Required fields
	expected.Title = gofakeit.Sentence(3)
	expected.HTML = gofakeit.Paragraph(1, 5, 10, " ")
	expected.TextContent = gofakeit.Paragraph(1, 5, 10, " ")
	expected.Published = time.Now()
	expected.Modified = time.Now()

	// Optional fields
	expected.Byline = gofakeit.Name()
	expected.Excerpt = gofakeit.Sentence(10)

	expected.Images = article.NewImages(&article.Image{
		ID:      gofakeit.UUID(),
		URL:     gofakeit.URL(),
		AltText: gofakeit.Sentence(5),
		Width:   gofakeit.Number(800, 1920),
		Height:  gofakeit.Number(600, 1080),
		Caption: gofakeit.Sentence(10),
	})

	expected.Videos = article.NewVideos(&article.Video{
		ID:        gofakeit.UUID(),
		URL:       gofakeit.URL(),
		EmbedCode: "<iframe src='" + gofakeit.URL() + "'></iframe>",
		Caption:   gofakeit.Sentence(10),
	})

	expected.Quotes = article.NewQuotes(&article.Quote{
		ID:       gofakeit.UUID(),
		Text:     gofakeit.Sentence(15),
		Author:   gofakeit.Name(),
		Source:   gofakeit.URL(),
		Platform: "Twitter",
	})

	expected.Tags = article.NewTags("travel", "Phuket", "Thailand")
	expected.Source = gofakeit.URL()
	expected.Language = "en"
	expected.Category = "Travel"
	expected.SiteName = "Example Travel Blog"

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
	expected.HTML = gofakeit.Paragraph(1, 5, 10, " ")
	expected.TextContent = gofakeit.Paragraph(1, 5, 10, " ")
	expected.Published = time.Now()
	expected.Modified = time.Now()

	// Optional fields with invalid nested structure
	expected.Byline = gofakeit.Name()
	expected.Excerpt = gofakeit.Sentence(10)

	expected.Images = article.NewImages(&article.Image{
		URL:     "invalid-url",
		AltText: gofakeit.Sentence(5),
		Width:   gofakeit.Number(800, 1920),
		Height:  gofakeit.Number(600, 1080),
		Caption: gofakeit.Sentence(10),
	})

	expected.Tags = article.NewTags("travel", "Phuket", "Thailand")
	expected.Source = gofakeit.URL()
	expected.Language = "en"
	expected.Category = "Travel"
	expected.SiteName = "Example Travel Blog"

	// Convert expected Article to map and then back to Article to simulate input processing
	inputMap := expected.Map()

	// Expect the images to be nil due to invalid URL in the nested structure
	expected.Images = article.NewImages()

	got, err := article.NewArticleFromMap(inputMap)
	require.NoError(t, err)

	// To compare Published and Modified separately due to possible time differences
	assert.Equal(t, expected.ID, got.ID)
	assert.Equal(t, expected.Title, got.Title)
	assert.Equal(t, expected.Byline, got.Byline)
	assert.Equal(t, expected.HTML, got.HTML)
	assert.Equal(t, expected.TextContent, got.TextContent)
	assert.Equal(t, expected.Excerpt, got.Excerpt)
	assert.Equal(t, expected.Images, got.Images)
	assert.WithinDuration(t, expected.Published, got.Published, time.Second)
	assert.WithinDuration(t, expected.Modified, got.Modified, time.Second)
	assert.Equal(t, expected.Tags, got.Tags)
	assert.Equal(t, expected.Source, got.Source)
	assert.Equal(t, expected.Language, got.Language)
	assert.Equal(t, expected.Category, got.Category)
	assert.Equal(t, expected.SiteName, got.SiteName)
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
	expected.HTML = gofakeit.Paragraph(1, 5, 10, " ")
	expected.TextContent = gofakeit.Paragraph(1, 5, 10, " ")
	expected.Published = time.Now()
	expected.Modified = time.Now()

	// Optional fields with some invalid data
	expected.Byline = gofakeit.Name()
	expected.Excerpt = gofakeit.Sentence(10)

	expected.Images = article.NewImages(&article.Image{
		URL:     "invalid-url",
		AltText: gofakeit.Sentence(5),
		Width:   gofakeit.Number(800, 1920),
		Height:  gofakeit.Number(600, 1080),
		Caption: gofakeit.Sentence(10),
	})

	expected.Videos = article.NewVideos(&article.Video{
		URL:       "invalid-url",
		EmbedCode: "<iframe src='invalid-url'></iframe>",
		Caption:   gofakeit.Sentence(10),
	})

	expected.Quotes = article.NewQuotes(&article.Quote{
		Text:     "",
		Author:   gofakeit.Name(),
		Source:   "invalid-url",
		Platform: "Twitter",
	})

	expected.Tags = article.NewTags("travel", "Phuket", "Thailand")
	expected.Source = "invalid-url"
	expected.Language = "en"
	expected.Category = "Travel"
	expected.SiteName = "Example Travel Blog"

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
	assert.Equal(t, "", expected.Source)
}

func TestArticleNormalizeFieldClearing(t *testing.T) {

	invalid := article.NewArticle()

	// Set required fields with valid data
	invalid.Title = gofakeit.Sentence(3)
	invalid.HTML = gofakeit.Paragraph(1, 5, 10, " ")
	invalid.TextContent = gofakeit.Paragraph(1, 5, 10, " ")
	invalid.Published = time.Now()
	invalid.Modified = time.Now()

	// Set invalid data for optional fields
	invalid.ID = "invalid-uuid"
	invalid.Byline = gofakeit.Name()
	invalid.Excerpt = gofakeit.Sentence(10)
	invalid.Source = "invalid-url"
	invalid.Language = "inglese" // should be a valid ISO 639-1 language code
	invalid.Category = gofakeit.Sentence(2)
	invalid.SiteName = gofakeit.Sentence(2)

	valid := *invalid
	(&valid).Normalize()

	// Verify that invalid fields are cleared
	assert.Equal(t, "", valid.ID)
	assert.Equal(t, invalid.Byline, valid.Byline)   // should not be cleared since it's not required
	assert.Equal(t, invalid.Excerpt, valid.Excerpt) // should not be cleared since it's not required
	assert.Equal(t, "", valid.Source)
	assert.Equal(t, "inglese", valid.Language)
	assert.Equal(t, invalid.Category, valid.Category)
	assert.Equal(t, invalid.SiteName, valid.SiteName)
}

// TestGetStringSlice tests the GetStringSlice function in case of empty map and missing key:
func TestGetStringSlice(t *testing.T) {
	m := map[string]interface{}{}
	key := "key"
	assert.Equal(t, []string{}, article.GetStringSlice(m, key))
}

func TestTrimToMaxLen(t *testing.T) {
	s := "This is a test string with more than twenty characters."
	trimmed := article.TrimToMaxLen(s, 20)
	assert.Equal(t, "This is a test strin", trimmed)

	s = "Short string"
	trimmed = article.TrimToMaxLen(s, 20)
	assert.Equal(t, s, trimmed)
}

func TestUnmarshal(t *testing.T) {

	js := `{
		  "article__id": "123e4567-e89b-12d3-a456-426614174000",
		  "article__title": "The Rise of AI",
		  "article__byline": "By John Doe",
		  "article__html": "<p>Artificial Intelligence is transforming the world.</p>",
		  "article__text": "Artificial Intelligence is transforming the world.",
		  "article__excerpt": "An overview of how AI is changing various industries.",
		  "article__published": "2024-05-27T10:00:00Z",
		  "article__modified": "2024-05-28T12:00:00Z",
		  "article__images": [
			  {
				"id": "img-001",
				"url": "https://example.com/image1.jpg",
				"alt_text": "AI Illustration",
				"width": 800,
				"height": 600,
				"caption": "An illustration representing AI."
			  }
          ],
		  "article__videos": [
			  {
				"id": "vid-001",
				"url": "https://example.com/video1.mp4",
				"embed_code": "<iframe src='https://example.com/video1.mp4'></iframe>",
				"caption": "A video explaining AI."
			  }
		  ],
		  "article__quotes": [
			  {
				"id": "quote-001",
				"text": "AI is the future of technology.",
				"author": "Jane Smith",
				"source": "https://twitter.com/janesmith/status/123",
				"platform": "Twitter"
			  }
		  ],
		  "article__tags": ["AI", "Technology", "Future"],
		  "article__source": "https://example.com",
		  "article__language": "en",
		  "article__category": "Technology",
		  "article__site": "Tech News",
		  "article__socials": [
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
	assert.Equal(t, "By John Doe", art.Byline)
	assert.Equal(t, "<p>Artificial Intelligence is transforming the world.</p>", art.HTML)
	assert.Equal(t, "Artificial Intelligence is transforming the world.", art.TextContent)
	assert.Equal(t, "An overview of how AI is changing various industries.", art.Excerpt)
	assert.Equal(t, "2024-05-27T10:00:00Z", art.Published.Format(time.RFC3339))
	assert.Equal(t, "2024-05-28T12:00:00Z", art.Modified.Format(time.RFC3339))

	// check the values of the nested structures
	assert.Equal(t, 1, art.Images.Len())
	assert.Equal(t, 1, art.Videos.Len())
	assert.Equal(t, 1, art.Quotes.Len())
	assert.Equal(t, 1, art.Socials.Len())
	assert.Equal(t, 3, art.Tags.Len())
}
