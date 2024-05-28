package article

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/go-shiori/go-readability"
	distiller "github.com/markusmobius/go-domdistiller"
	"github.com/samber/lo"
	"net/url"
	"strings"
	"time"
)

//type Article struct {
//	ID    string `json:"id" validate:"required,uuid4,max=36"`
//	Title string `json:"title" validate:"required,max=255"`
//	// Summary is a short description of the article.
//	Summary string `json:"summary" validate:"max=500"`
//	// Markup is the raw HTML or Markdown content of the article.
//	Markup string `json:"markup" validate:"required,max=65000"`
//	// Text plain text content of the article.
//	Text string `json:"text" validate:"required,max=65000"`
//	// Genre is a short summary or preview of the article.
//	Genre   string    `json:"genre" validate:"max=500"`
//	SourceURL    string    `json:"source_url" validate:"omitempty,url,max=4096"`
//	Language  string    `json:"language" validate:"max=255"`
//	Category  string    `json:"category" validate:"max=255"`
//	SourceName  string    `json:"source_name" validate:"max=255"`
//	Published time.Time `json:"published" validate:"required"`
//	Modified  time.Time `json:"modified"`
//	Images    *Images   `json:"images"`
//	Videos    *Videos   `json:"videos"`
//	Quotes    *Quotes   `json:"quotes"`
//	Tags      *Tags     `json:"tags"`
//	Socials   *Socials  `json:"socials"`
//}

// Extract Article from HTML
func FromHTML(html string, resource *url.URL) (*Article, error) {

	a := NewArticle()
	a.SourceURL = resource.String()

	// readability: title, summary, text, html, language
	readabilityArticle(html, resource, a)

	// distiller: category, images, source name
	distillArticle(html, resource, a)

	// fallback: published
	a.Published = lo.Ternary(a.Published.IsZero(), fallbackPublished(html), a.Published)

	// nil article if it's invalid
	if err := a.Normalize(); err != nil {
		return nil, err
	}

	return a, nil
}

func readabilityArticle(html string, resource *url.URL, a *Article) {

	read, err := readability.FromReader(strings.NewReader(html), resource)
	if err != nil {
		return
	}

	// set the article fields
	a.Title = read.Title
	a.Summary = read.Excerpt
	a.Text = read.TextContent
	a.Language = read.Language

	a.Title = lo.Ternary(a.Title == "", read.Title, a.Title)
	a.Summary = lo.Ternary(a.Summary == "", read.Excerpt, a.Summary)
	a.Text = lo.Ternary(a.Text == "", read.TextContent, a.Text)
	a.Language = lo.Ternary(a.Language == "", read.Language, a.Language)

	a.Markup = read.Content

	if read.PublishedTime != nil {
		a.Published = *read.PublishedTime
	}
}

func distillArticle(html string, resource *url.URL, a *Article) {

	distill, err := distiller.ApplyForReader(strings.NewReader(html), &distiller.Options{
		OriginalURL: resource,
	})
	if err != nil {
		return
	}

	info := distill.MarkupInfo

	// set the article fields
	a.Category = info.Article.Section
	a.SourceName = info.Publisher
	a.Images = distillImages(distill, resource)

	// fallback fields applied only if the fields are empty
	a.Title = lo.Ternary(a.Title == "", distill.Title, a.Title)
	a.Summary = lo.Ternary(a.Summary == "", info.Description, a.Summary)
	a.Published = lo.Ternary(a.Published.IsZero(), distillPublished(distill), a.Published)
}

func distillPublished(distill *distiller.Result) time.Time {

	publishedStr := distill.MarkupInfo.Article.PublishedTime
	published, timeErr := time.Parse(time.RFC3339, publishedStr)
	if timeErr == nil {
		return time.Now()
	}
	return published
}

func distillImages(distill *distiller.Result, resource *url.URL) *Images {

	images := NewImages()

	for _, src := range distill.MarkupInfo.Images {
		image := NewImage(AbsoluteUrl(resource, src.URL))
		image.Width = src.Width
		image.Height = src.Height
		images.Add()
	}

	return images
}

func fallbackPublished(html string) time.Time {

	fallback := time.Now()

	q, readerErr := goquery.NewDocumentFromReader(strings.NewReader(html))
	if readerErr != nil {
		return fallback
	}

	// .field--name-created
	if el := q.Find(".field--name-created").Text(); len(el) > 0 {
		if published, err := time.Parse("Monday, 2 January 2006", el); err == nil {
			return published
		}
	}

	return fallback
}

func AbsoluteUrl(base *url.URL, href string) string {

	// parse the href
	rel, err := url.Parse(href)
	if err != nil {
		return ""
	}

	// resolve the base with the relative href
	abs := base.ResolveReference(rel)

	return abs.String()
}
