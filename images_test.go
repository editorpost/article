package article_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/editorpost/article"
	"github.com/stretchr/testify/assert"
)

func fakeImage() *article.Image {
	return &article.Image{
		ID: gofakeit.UUID(),
	}
}

func TestNewImagesStrict(t *testing.T) {

	t.Run("valid images", func(t *testing.T) {
		img1 := fakeImage()
		img2 := fakeImage()
		images, err := article.NewImagesStrict(img1, img2)
		assert.NoError(t, err)
		assert.Equal(t, 2, images.Count())
	})

	t.Run("invalid images", func(t *testing.T) {
		img1 := fakeImage()
		img2 := &article.Image{ID: ""}
		images, err := article.NewImagesStrict(img1, img2)
		assert.Error(t, err)
		assert.Nil(t, images)
	})
}

func TestNewImages(t *testing.T) {

	t.Run("mixed valid and invalid images", func(t *testing.T) {
		img1 := fakeImage()
		img2 := &article.Image{ID: ""}
		img3 := fakeImage()
		images := article.NewImages(img1, img2, img3)
		assert.Equal(t, 2, images.Count())
	})
}

func TestImages_Get(t *testing.T) {
	img1 := fakeImage()
	img2 := fakeImage()
	images := article.NewImages(img1, img2)

	t.Run("existing image", func(t *testing.T) {
		img, found := images.Get(img1.ID)
		assert.True(t, found)
		assert.Equal(t, img1, img)
	})

	t.Run("non-existing image", func(t *testing.T) {
		_, found := images.Get(gofakeit.UUID())
		assert.False(t, found)
	})
}

func TestImages_Add(t *testing.T) {
	images := article.NewImages()
	img1 := fakeImage()
	img2 := &article.Image{ID: ""}
	img3 := fakeImage()

	images.Add(img1, img2, img3)
	assert.Equal(t, 2, images.Count())
}

func TestImages_Remove(t *testing.T) {
	img1 := fakeImage()
	img2 := fakeImage()
	images := article.NewImages(img1, img2)

	t.Run("remove existing image", func(t *testing.T) {
		images.Remove(img1.ID)
		assert.Equal(t, 1, images.Count())
		_, found := images.Get(img1.ID)
		assert.False(t, found)
	})

	t.Run("remove non-existing image", func(t *testing.T) {
		initialCount := images.Count()
		images.Remove(gofakeit.UUID())
		assert.Equal(t, initialCount, images.Count())
	})
}

func TestImages_IDs(t *testing.T) {
	img1 := fakeImage()
	img2 := fakeImage()
	images := article.NewImages(img1, img2)

	ids := images.IDs()
	assert.Equal(t, 2, len(ids))
	assert.Contains(t, ids, img1.ID)
	assert.Contains(t, ids, img2.ID)
}

func TestImages_Count(t *testing.T) {
	images := article.NewImages()
	assert.Equal(t, 0, images.Count())

	img1 := fakeImage()
	img2 := fakeImage()
	images.Add(img1, img2)
	assert.Equal(t, 2, images.Count())
}

func TestImages_Filter(t *testing.T) {
	img1 := fakeImage()
	img2 := fakeImage()
	images := article.NewImages(img1, img2)

	t.Run("filter with always true function", func(t *testing.T) {
		filtered := images.Filter(func(img *article.Image) bool {
			return true
		})
		assert.Equal(t, 2, filtered.Count())
	})

	t.Run("filter with always false function", func(t *testing.T) {
		filtered := images.Filter(func(img *article.Image) bool {
			return false
		})
		assert.Equal(t, 0, filtered.Count())
	})

	t.Run("filter with specific condition", func(t *testing.T) {
		filtered := images.Filter(func(img *article.Image) bool {
			return img.ID == img1.ID
		})
		assert.Equal(t, 1, filtered.Count())
		assert.Equal(t, img1.ID, filtered.IDs()[0])
	})
}
