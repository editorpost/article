package article_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/editorpost/article"
	"github.com/stretchr/testify/assert"
)

func NewMediaValid() *article.Media {
	return &article.Media{
		ID:  gofakeit.UUID(),
		URL: gofakeit.URL(),
	}
}

func NewMediaInvalid() *article.Media {
	return &article.Media{
		ID:  gofakeit.UUID(),
		URL: "invalid-url",
	}
}

func TestNewMediasStrict(t *testing.T) {

	t.Run("valid items", func(t *testing.T) {
		img1 := NewMediaValid()
		img2 := NewMediaValid()
		medias, err := article.NewMediasStrict(img1, img2)
		assert.NoError(t, err)
		assert.Equal(t, 2, medias.Len())
	})

	t.Run("invalid items", func(t *testing.T) {
		img1 := NewMediaValid()
		img2 := &article.Media{ID: ""}
		medias, err := article.NewMediasStrict(img1, img2)
		assert.Error(t, err)
		assert.Nil(t, medias)
	})
}

func TestNewMedias(t *testing.T) {

	t.Run("mixed valid and invalid items", func(t *testing.T) {
		img1 := NewMediaValid()
		img2 := &article.Media{ID: ""}
		img3 := NewMediaValid()
		medias := article.NewMedias(img1, img2, img3)
		assert.Equal(t, 2, medias.Len())
	})
}

func TestMedias_Get(t *testing.T) {
	img1 := NewMediaValid()
	img2 := NewMediaValid()
	medias := article.NewMedias(img1, img2)

	t.Run("existing media", func(t *testing.T) {
		img, found := medias.Get(img1.ID)
		assert.True(t, found)
		assert.Equal(t, img1, img)
	})

	t.Run("non-existing media", func(t *testing.T) {
		_, found := medias.Get(gofakeit.UUID())
		assert.False(t, found)
	})
}

func TestMedias_Add(t *testing.T) {
	medias := article.NewMedias()
	img1 := NewMediaValid()
	img2 := &article.Media{ID: ""}
	img3 := NewMediaValid()

	medias.Add(img1, img2, img3)
	assert.Equal(t, 2, medias.Len())
}

func TestMedias_Remove(t *testing.T) {
	img1 := NewMediaValid()
	img2 := NewMediaValid()
	medias := article.NewMedias(img1, img2)

	t.Run("remove existing media", func(t *testing.T) {
		medias.Remove(img1.ID)
		assert.Equal(t, 1, medias.Len())
		_, found := medias.Get(img1.ID)
		assert.False(t, found)
	})

	t.Run("remove non-existing media", func(t *testing.T) {
		initialCount := medias.Len()
		medias.Remove(gofakeit.UUID())
		assert.Equal(t, initialCount, medias.Len())
	})
}

func TestMedias_IDs(t *testing.T) {
	img1 := NewMediaValid()
	img2 := NewMediaValid()
	medias := article.NewMedias(img1, img2)

	ids := medias.IDs()
	assert.Equal(t, 2, len(ids))
	assert.Contains(t, ids, img1.ID)
	assert.Contains(t, ids, img2.ID)
}

func TestMedias_Count(t *testing.T) {
	medias := article.NewMedias()
	assert.Equal(t, 0, medias.Len())

	img1 := NewMediaValid()
	img2 := NewMediaValid()
	medias.Add(img1, img2)
	assert.Equal(t, 2, medias.Len())
}

func TestMedias_Filter(t *testing.T) {
	img1 := NewMediaValid()
	img2 := NewMediaValid()
	medias := article.NewMedias(img1, img2)

	t.Run("filter with always true function", func(t *testing.T) {
		filtered := medias.Filter(func(img *article.Media) bool {
			return true
		})
		assert.Equal(t, 2, filtered.Len())
	})

	t.Run("filter with always false function", func(t *testing.T) {
		filtered := medias.Filter(func(img *article.Media) bool {
			return false
		})
		assert.Equal(t, 0, filtered.Len())
	})

	t.Run("filter with specific condition", func(t *testing.T) {
		filtered := medias.Filter(func(img *article.Media) bool {
			return img.ID == img1.ID
		})
		assert.Equal(t, 1, filtered.Len())
		assert.Equal(t, img1.ID, filtered.IDs()[0])
	})
}
