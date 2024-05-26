package article_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/editorpost/article"
	"github.com/stretchr/testify/assert"
)

func fakeQuote() *article.Quote {
	return &article.Quote{
		ID:       gofakeit.UUID(),
		Text:     gofakeit.LoremIpsumSentence(10),
		Author:   gofakeit.Name(),
		Source:   gofakeit.URL(),
		Platform: gofakeit.BuzzWord(),
	}
}

func TestNewQuotesStrict(t *testing.T) {
	t.Run("valid items", func(t *testing.T) {
		quote1 := fakeQuote()
		quote2 := fakeQuote()
		quotes, err := article.NewQuotesStrict(quote1, quote2)
		assert.NoError(t, err)
		assert.Equal(t, 2, quotes.Len())
	})

	t.Run("invalid items", func(t *testing.T) {
		quote1 := fakeQuote()
		quote2 := &article.Quote{ID: "", Text: ""}
		quotes, err := article.NewQuotesStrict(quote1, quote2)
		assert.Error(t, err)
		assert.Nil(t, quotes)
	})
}

func TestNewQuotes(t *testing.T) {
	t.Run("mixed valid and invalid items", func(t *testing.T) {
		quote1 := fakeQuote()
		quote2 := &article.Quote{ID: "", Text: ""}
		quote3 := fakeQuote()
		quotes := article.NewQuotes(quote1, quote2, quote3)
		assert.Equal(t, 2, quotes.Len())
	})
}

func TestQuotes_Get(t *testing.T) {
	quote1 := fakeQuote()
	quote2 := fakeQuote()
	quotes := article.NewQuotes(quote1, quote2)

	t.Run("existing quote", func(t *testing.T) {
		quote, found := quotes.Get(quote1.ID)
		assert.True(t, found)
		assert.Equal(t, quote1, quote)
	})

	t.Run("non-existing quote", func(t *testing.T) {
		_, found := quotes.Get(gofakeit.UUID())
		assert.False(t, found)
	})
}

func TestQuotes_Add(t *testing.T) {
	quotes := article.NewQuotes()
	quote1 := fakeQuote()
	quote2 := &article.Quote{ID: "", Text: ""}
	quote3 := fakeQuote()

	quotes.Add(quote1, quote2, quote3)
	assert.Equal(t, 2, quotes.Len())
}

func TestQuotes_Remove(t *testing.T) {
	quote1 := fakeQuote()
	quote2 := fakeQuote()
	quotes := article.NewQuotes(quote1, quote2)

	t.Run("remove existing quote", func(t *testing.T) {
		quotes.Remove(quote1.ID)
		assert.Equal(t, 1, quotes.Len())
		_, found := quotes.Get(quote1.ID)
		assert.False(t, found)
	})

	t.Run("remove non-existing quote", func(t *testing.T) {
		initialCount := quotes.Len()
		quotes.Remove(gofakeit.UUID())
		assert.Equal(t, initialCount, quotes.Len())
	})
}

func TestQuotes_IDs(t *testing.T) {
	quote1 := fakeQuote()
	quote2 := fakeQuote()
	quotes := article.NewQuotes(quote1, quote2)

	ids := quotes.IDs()
	assert.Equal(t, 2, len(ids))
	assert.Contains(t, ids, quote1.ID)
	assert.Contains(t, ids, quote2.ID)
}

func TestQuotes_Count(t *testing.T) {
	quotes := article.NewQuotes()
	assert.Equal(t, 0, quotes.Len())

	quote1 := fakeQuote()
	quote2 := fakeQuote()
	quotes.Add(quote1, quote2)
	assert.Equal(t, 2, quotes.Len())
}

func TestQuotes_Filter(t *testing.T) {
	quote1 := fakeQuote()
	quote2 := fakeQuote()
	quotes := article.NewQuotes(quote1, quote2)

	t.Run("filter with always true function", func(t *testing.T) {
		filtered := quotes.Filter(func(quote *article.Quote) bool {
			return true
		})
		assert.Equal(t, 2, filtered.Len())
	})

	t.Run("filter with always false function", func(t *testing.T) {
		filtered := quotes.Filter(func(quote *article.Quote) bool {
			return false
		})
		assert.Equal(t, 0, filtered.Len())
	})

	t.Run("filter with specific condition", func(t *testing.T) {
		filtered := quotes.Filter(func(quote *article.Quote) bool {
			return quote.ID == quote1.ID
		})
		assert.Equal(t, 1, filtered.Len())
		assert.Equal(t, quote1.ID, filtered.IDs()[0])
	})
}

func TestQuotes_Normalize(t *testing.T) {
	valid := fakeQuote()
	invalid := &article.Quote{ID: gofakeit.UUID(), Text: " sdf "}

	quotes := article.NewQuotes(valid, invalid)

	// Normalize should remove invalid quote
	quotes.Normalize()

	_, exist := quotes.Get(invalid.ID)
	assert.False(t, exist)

	// Normalize should keep valid quote
	_, exist = quotes.Get(valid.ID)
	assert.True(t, exist)
}
