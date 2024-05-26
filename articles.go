package article

import "encoding/json"

// Article represents a news article with various types of content.
// This structure provides a flexible and universal foundation for storing and working with various types of content,
// allowing for easy creation and modification of articles, as well as integration of media and social elements.
//type Article struct {
//	ID          string    `json:"id" validate:"required,uuid4,max=36"`
//	Title       string    `json:"title" validate:"required,max=255"`
//	Summary      string    `json:"summary" validate:"max=255"`
//	HTML        string    `json:"html" validate:"required,max=65000"`
//	Text string    `json:"text" validate:"required,max=65000"`
//	Excerpt     string    `json:"excerpt" validate:"max=500"`
//	Published   time.Time `json:"published" validate:"required"`
//	Modified    time.Time `json:"modified"`
//	Images      *Images   `json:"images"`
//	Videos      *Videos   `json:"videos"`
//	Quotes      *Quotes   `json:"quotes"`
//	Tags        *Tags     `json:"tags"`
//	Socials     *Socials  `json:"socials"`
//	Source      string    `json:"source" validate:"omitempty,url,max=4096"`
//	Language    string    `json:"language" validate:"max=255"`
//	Category    string    `json:"category" validate:"max=255"`
//	SiteName    string    `json:"site" validate:"max=255"`
//}

type Articles struct {
	items []*Article
}

func NewArticles(articles ...*Article) *Articles {

	var valid []*Article

	for _, article := range articles {
		if article != nil {
			valid = append(valid, article)
		}
	}

	return &Articles{items: valid}
}

// Get returns the article by ID
func (list *Articles) Get(id string) (*Article, bool) {
	for _, article := range list.items {
		if article.ID == id {
			return article, true
		}
	}
	return nil, false
}

// Slice returns a slice of all articles
func (list *Articles) Slice() []*Article {
	return list.items
}

// Add adds articles to the collection
func (list *Articles) Add(articles ...*Article) *Articles {
	for _, article := range articles {
		if article != nil {
			list.items = append(list.items, article)
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

	var filtered []*Article
	for _, article := range list.items {
		if _, found := idSet[article.ID]; !found {
			filtered = append(filtered, article)
		}
	}

	list.items = filtered
	return list
}

// IDs returns a slice of all article IDs
func (list *Articles) IDs() []string {
	ids := make([]string, len(list.items))
	for i, article := range list.items {
		ids[i] = article.ID
	}
	return ids
}

// Normalize nested structures
func (list *Articles) Normalize() {
	for _, article := range list.items {
		article.Normalize()
	}
}

// Maps converts the Articles struct to a []map[string]any, including nested structures.
func (list *Articles) Maps() []map[string]any {

	result := make([]map[string]any, 0, len(list.items))

	for _, article := range list.items {
		result = append(result, article.Map())
	}

	return result
}

// FilterFn is a function to filter articles
type FilterFn func(article *Article) bool

// Filter filters articles using the provided function
func (list *Articles) Filter(fn FilterFn) *Articles {
	var filtered []*Article
	for _, article := range list.items {
		if fn(article) {
			filtered = append(filtered, article)
		}
	}
	return &Articles{items: filtered}
}

func (list *Articles) Images() *Images {
	images := NewImages()
	for _, article := range list.items {
		images.Add(article.Images.Slice()...)
	}
	return images
}

func (list *Articles) Videos() *Videos {
	videos := NewVideos()
	for _, article := range list.items {
		videos.Add(article.Videos.Slice()...)
	}
	return videos
}

func (list *Articles) Quotes() *Quotes {
	quotes := NewQuotes()
	for _, article := range list.items {
		quotes.Add(article.Quotes.Slice()...)
	}
	return quotes
}

func (list *Articles) Socials() *Socials {
	socials := NewSocials()
	for _, article := range list.items {
		socials.Add(article.Socials.Slice()...)
	}
	return socials
}

func (list *Articles) Tags() *Tags {
	tags := NewTags()
	for _, article := range list.items {
		tags.Add(article.Tags.Slice()...)
	}
	return tags
}

func (list *Articles) Len() int {
	return len(list.items)
}

// UnmarshalJSON to array of items using encoding/json
func (list *Articles) UnmarshalJSON(data []byte) error {

	// Unmarshal to a slice of Image
	var items []*Article
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}

	// Create a new Images collection
	*list = *NewArticles(items...)

	return nil
}

// MarshalJSON from array of items using encoding/json
func (list *Articles) MarshalJSON() ([]byte, error) {
	return json.Marshal(list.items)
}
