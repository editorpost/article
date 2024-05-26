package article_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/editorpost/article"
	"github.com/stretchr/testify/assert"
)

func fakeVideo() *article.Video {
	return &article.Video{
		ID:        gofakeit.UUID(),
		URL:       gofakeit.URL(),
		EmbedCode: gofakeit.LoremIpsumSentence(10),
		Caption:   gofakeit.LoremIpsumSentence(5),
	}
}

func TestNewVideosStrict(t *testing.T) {
	t.Run("valid videos", func(t *testing.T) {
		video1 := fakeVideo()
		video2 := fakeVideo()
		videos, err := article.NewVideosStrict(video1, video2)
		assert.NoError(t, err)
		assert.Equal(t, 2, videos.Count())
	})

	t.Run("invalid videos", func(t *testing.T) {
		video1 := fakeVideo()
		video2 := &article.Video{ID: "", URL: "invalid"}
		videos, err := article.NewVideosStrict(video1, video2)
		assert.Error(t, err)
		assert.Nil(t, videos)
	})
}

func TestNewVideos(t *testing.T) {
	t.Run("mixed valid and invalid videos", func(t *testing.T) {
		video1 := fakeVideo()
		video2 := &article.Video{ID: "", URL: "invalid"}
		video3 := fakeVideo()
		videos := article.NewVideos(video1, video2, video3)
		assert.Equal(t, 2, videos.Count())
	})
}

func TestVideos_Get(t *testing.T) {
	video1 := fakeVideo()
	video2 := fakeVideo()
	videos := article.NewVideos(video1, video2)

	t.Run("existing video", func(t *testing.T) {
		video, found := videos.Get(video1.ID)
		assert.True(t, found)
		assert.Equal(t, video1, video)
	})

	t.Run("non-existing video", func(t *testing.T) {
		_, found := videos.Get(gofakeit.UUID())
		assert.False(t, found)
	})
}

func TestVideos_Add(t *testing.T) {
	videos := article.NewVideos()
	video1 := fakeVideo()
	video2 := &article.Video{ID: "", URL: "invalid"}
	video3 := fakeVideo()

	videos.Add(video1, video2, video3)
	assert.Equal(t, 2, videos.Count())
}

func TestVideos_Remove(t *testing.T) {
	video1 := fakeVideo()
	video2 := fakeVideo()
	videos := article.NewVideos(video1, video2)

	t.Run("remove existing video", func(t *testing.T) {
		videos.Remove(video1.ID)
		assert.Equal(t, 1, videos.Count())
		_, found := videos.Get(video1.ID)
		assert.False(t, found)
	})

	t.Run("remove non-existing video", func(t *testing.T) {
		initialCount := videos.Count()
		videos.Remove(gofakeit.UUID())
		assert.Equal(t, initialCount, videos.Count())
	})
}

func TestVideos_IDs(t *testing.T) {
	video1 := fakeVideo()
	video2 := fakeVideo()
	videos := article.NewVideos(video1, video2)

	ids := videos.IDs()
	assert.Equal(t, 2, len(ids))
	assert.Contains(t, ids, video1.ID)
	assert.Contains(t, ids, video2.ID)
}

func TestVideos_Count(t *testing.T) {
	videos := article.NewVideos()
	assert.Equal(t, 0, videos.Count())

	video1 := fakeVideo()
	video2 := fakeVideo()
	videos.Add(video1, video2)
	assert.Equal(t, 2, videos.Count())
}

func TestVideos_Filter(t *testing.T) {
	video1 := fakeVideo()
	video2 := fakeVideo()
	videos := article.NewVideos(video1, video2)

	t.Run("filter with always true function", func(t *testing.T) {
		filtered := videos.Filter(func(video *article.Video) bool {
			return true
		})
		assert.Equal(t, 2, filtered.Count())
	})

	t.Run("filter with always false function", func(t *testing.T) {
		filtered := videos.Filter(func(video *article.Video) bool {
			return false
		})
		assert.Equal(t, 0, filtered.Count())
	})

	t.Run("filter with specific condition", func(t *testing.T) {
		filtered := videos.Filter(func(video *article.Video) bool {
			return video.ID == video1.ID
		})
		assert.Equal(t, 1, filtered.Count())
		assert.Equal(t, video1.ID, filtered.IDs()[0])
	})
}

func TestVideos_Normalize(t *testing.T) {

	valid := fakeVideo()
	invalid := &article.Video{ID: gofakeit.UUID(), URL: " sdf "}

	videos := article.NewVideos(valid, invalid)

	// Normalize should remove invalid video
	_, exist := videos.Get(invalid.ID)
	assert.False(t, exist)

	// Normalize should keep valid video
	_, exist = videos.Get(valid.ID)
	assert.True(t, exist)
}
