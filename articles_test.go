package article_test

import (
	"github.com/editorpost/article"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Articles represents a collection of Article pointers
type Articles struct {
	items []*article.Article
}

func NewArticles(articles ...*article.Article) *Articles {

	var valid []*article.Article

	for _, item := range articles {
		if item != nil {
			valid = append(valid, item)
		}
	}

	return &Articles{items: valid}
}

// Get returns the article by ID
func (list *Articles) Get(id string) (*article.Article, bool) {
	for _, item := range list.items {
		if item.ID == id {
			return item, true
		}
	}
	return nil, false
}

// Slice returns a slice of all articles
func (list *Articles) Slice() []*article.Article {
	return list.items
}

// Add adds articles to the collection
func (list *Articles) Add(articles ...*article.Article) *Articles {
	for _, item := range articles {
		if item != nil {
			list.items = append(list.items, item)
		}
	}
	return list
}

// Remove removes articles from the collection
func (list *Articles) Remove(ids ...string) *Articles {
	idSet := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}

	var filtered []*article.Article
	for _, item := range list.items {
		if _, found := idSet[item.ID]; !found {
			filtered = append(filtered, item)
		}
	}

	list.items = filtered
	return list
}

// IDs returns a slice of all article IDs
func (list *Articles) IDs() []string {
	ids := make([]string, len(list.items))
	for i, item := range list.items {
		ids[i] = item.ID
	}
	return ids
}

// Normalize nested structures
func (list *Articles) Normalize() {
	for _, item := range list.items {
		item.Normalize()
	}
}

// Maps converts the Articles struct to a []map[string]any, including nested structures.
func (list *Articles) Maps() []map[string]any {

	var result []map[string]any

	for _, item := range list.items {
		result = append(result, item.Map())
	}

	return result
}

// FilterFn is a function to filter articles
type FilterFn func(article *article.Article) bool

// Filter filters articles using the provided function
func (list *Articles) Filter(fn FilterFn) *Articles {
	var filtered []*article.Article
	for _, item := range list.items {
		if fn(item) {
			filtered = append(filtered, item)
		}
	}
	return &Articles{items: filtered}
}

func (list *Articles) Images() *article.Images {
	images := article.NewImages()
	for _, item := range list.items {
		images.Add(item.Images.Slice()...)
	}
	return images
}

func (list *Articles) Videos() *article.Videos {
	videos := article.NewVideos()
	for _, item := range list.items {
		videos.Add(item.Videos.Slice()...)
	}
	return videos
}

func (list *Articles) Quotes() *article.Quotes {
	quotes := article.NewQuotes()
	for _, item := range list.items {
		quotes.Add(item.Quotes.Slice()...)
	}
	return quotes
}

func (list *Articles) Socials() *article.Socials {
	socials := article.NewSocials()
	for _, item := range list.items {
		socials.Add(item.Socials.Slice()...)
	}
	return socials
}

func (list *Articles) Tags() *article.Tags {
	tags := article.NewTags()
	for _, item := range list.items {
		tags.Add(item.Tags.Slice()...)
	}
	return tags
}

func (list *Articles) Len() int {
	return len(list.items)
}

func TestArticles_Add(t *testing.T) {

	articles := NewArticles()

	assert.Equal(t, articles.Len(), 0)

	article1 := article.NewArticle()
	article2 := article.NewArticle()

	articles.Add(article1, article2)
	assert.Equal(t, articles.Len(), 2)
}

func TestArticles_Remove(t *testing.T) {

	articles := NewArticles()

	assert.Equal(t, articles.Len(), 0)

	article1 := article.NewArticle()
	article2 := article.NewArticle()

	articles.Add(article1, article2)
	assert.Equal(t, articles.Len(), 2)

	articles.Remove(article1.ID)
	assert.Equal(t, articles.Len(), 1)
}

func TestArticles_Get(t *testing.T) {

	articles := NewArticles()

	assert.Equal(t, articles.Len(), 0)

	article1 := article.NewArticle()
	article2 := article.NewArticle()

	articles.Add(article1, article2)
	assert.Equal(t, articles.Len(), 2)

	article, found := articles.Get(article1.ID)
	assert.True(t, found)
	assert.Equal(t, article, article1)
}

func TestArticles_Slice(t *testing.T) {

	articles := NewArticles()

	assert.Equal(t, articles.Len(), 0)

	article1 := article.NewArticle()
	article2 := article.NewArticle()

	articles.Add(article1, article2)
	assert.Equal(t, articles.Len(), 2)

	slice := articles.Slice()
	assert.Equal(t, len(slice), 2)
}

func TestArticles_IDs(t *testing.T) {

	articles := NewArticles()

	assert.Equal(t, articles.Len(), 0)

	article1 := article.NewArticle()
	article2 := article.NewArticle()

	articles.Add(article1, article2)
	assert.Equal(t, articles.Len(), 2)

	ids := articles.IDs()
	assert.Equal(t, len(ids), 2)
}

func TestArticles_Filter(t *testing.T) {

	articles := NewArticles()

	assert.Equal(t, articles.Len(), 0)

	article1 := article.NewArticle()
	article2 := article.NewArticle()

	articles.Add(article1, article2)
	assert.Equal(t, articles.Len(), 2)

	filtered := articles.Filter(func(article *article.Article) bool {
		return article.ID == article1.ID
	})
	assert.Equal(t, filtered.Len(), 1)
}

func TestArticles_Normalize(t *testing.T) {

	articles := NewArticles()

	assert.Equal(t, articles.Len(), 0)

	article1 := article.NewArticle()
	article2 := article.NewArticle()

	articles.Add(article1, article2)
	assert.Equal(t, articles.Len(), 2)

	articles.Normalize()
}

func TestArticles_Maps(t *testing.T) {

	articles := NewArticles()

	assert.Equal(t, articles.Len(), 0)

	article1 := article.NewArticle()
	article2 := article.NewArticle()

	articles.Add(article1, article2)
	assert.Equal(t, articles.Len(), 2)

	maps := articles.Maps()
	assert.Equal(t, len(maps), 2)
}

func TestArticles_Images(t *testing.T) {

	articles := NewArticles()

	assert.Equal(t, articles.Len(), 0)

	article1 := article.NewArticle()
	article2 := article.NewArticle()

	articles.Add(article1, article2)
	assert.Equal(t, articles.Len(), 2)

	images := articles.Images()
	assert.NotNil(t, images)
}

func TestArticles_Videos(t *testing.T) {

	articles := NewArticles()

	assert.Equal(t, articles.Len(), 0)

	article1 := article.NewArticle()
	article2 := article.NewArticle()

	articles.Add(article1, article2)
	assert.Equal(t, articles.Len(), 2)

	videos := articles.Videos()
	assert.NotNil(t, videos)
}

func TestArticles_Quotes(t *testing.T) {

	articles := NewArticles()

	assert.Equal(t, articles.Len(), 0)

	article1 := article.NewArticle()
	article2 := article.NewArticle()

	articles.Add(article1, article2)
	assert.Equal(t, articles.Len(), 2)

	quotes := articles.Quotes()
	assert.NotNil(t, quotes)
}

func TestArticles_Socials(t *testing.T) {

	articles := NewArticles()

	assert.Equal(t, articles.Len(), 0)

	article1 := article.NewArticle()
	article2 := article.NewArticle()

	articles.Add(article1, article2)
	assert.Equal(t, articles.Len(), 2)

	socials := articles.Socials()
	assert.NotNil(t, socials)
}

func TestArticles_Tags(t *testing.T) {

	articles := NewArticles()

	assert.Equal(t, articles.Len(), 0)

	article1 := article.NewArticle()
	article2 := article.NewArticle()

	articles.Add(article1, article2)
	assert.Equal(t, articles.Len(), 2)

	tags := articles.Tags()
	assert.NotNil(t, tags)
}

func TestArticles_Len(t *testing.T) {

	articles := NewArticles()

	assert.Equal(t, articles.Len(), 0)

	article1 := article.NewArticle()
	article2 := article.NewArticle()

	articles.Add(article1, article2)
	assert.Equal(t, articles.Len(), 2)
}

func TestArticles_NewArticles(t *testing.T) {

	articles := NewArticles()

	assert.NotNil(t, articles)
}

func TestArticles_FilterFn(t *testing.T) {

	articles := NewArticles()

	assert.NotNil(t, articles)
}
