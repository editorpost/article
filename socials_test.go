package article_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/editorpost/article"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

var validate = validator.New()

func fakeSocialProfile() *article.Social {
	return &article.Social{
		ID:       gofakeit.UUID(),
		Platform: gofakeit.BuzzWord(),
		URL:      gofakeit.URL(),
	}
}

func TestNewSocialProfilesStrict(t *testing.T) {
	t.Run("valid social items", func(t *testing.T) {
		profile1 := fakeSocialProfile()
		profile2 := fakeSocialProfile()
		profiles, err := article.NewSocialsStrict(profile1, profile2)
		assert.NoError(t, err)
		assert.Equal(t, 2, profiles.Len())
	})

	t.Run("invalid social items", func(t *testing.T) {
		profile1 := fakeSocialProfile()
		profile2 := &article.Social{ID: "", URL: ""}
		profiles, err := article.NewSocialsStrict(profile1, profile2)
		assert.Error(t, err)
		assert.Nil(t, profiles)
	})
}

func TestNewSocialProfiles(t *testing.T) {
	t.Run("mixed valid and invalid social items", func(t *testing.T) {
		profile1 := fakeSocialProfile()
		profile2 := &article.Social{ID: "", URL: ""}
		profile3 := fakeSocialProfile()
		profiles := article.NewSocials(profile1, profile2, profile3)
		assert.Equal(t, 2, profiles.Len())
	})
}

func TestSocialProfiles_Get(t *testing.T) {
	profile1 := fakeSocialProfile()
	profile2 := fakeSocialProfile()
	profiles := article.NewSocials(profile1, profile2)

	t.Run("existing social profile", func(t *testing.T) {
		profile, found := profiles.Get(profile1.ID)
		assert.True(t, found)
		assert.Equal(t, profile1, profile)
	})

	t.Run("non-existing social profile", func(t *testing.T) {
		_, found := profiles.Get(gofakeit.UUID())
		assert.False(t, found)
	})
}

func TestSocialProfiles_Add(t *testing.T) {
	profiles := article.NewSocials()
	profile1 := fakeSocialProfile()
	profile2 := &article.Social{ID: "", URL: ""}
	profile3 := fakeSocialProfile()

	profiles.Add(profile1, profile2, profile3)
	assert.Equal(t, 2, profiles.Len())
}

func TestSocialProfiles_Remove(t *testing.T) {
	profile1 := fakeSocialProfile()
	profile2 := fakeSocialProfile()
	profiles := article.NewSocials(profile1, profile2)

	t.Run("remove existing social profile", func(t *testing.T) {
		profiles.Remove(profile1.ID)
		assert.Equal(t, 1, profiles.Len())
		_, found := profiles.Get(profile1.ID)
		assert.False(t, found)
	})

	t.Run("remove non-existing social profile", func(t *testing.T) {
		initialCount := profiles.Len()
		profiles.Remove(gofakeit.UUID())
		assert.Equal(t, initialCount, profiles.Len())
	})
}

func TestSocialProfiles_IDs(t *testing.T) {
	profile1 := fakeSocialProfile()
	profile2 := fakeSocialProfile()
	profiles := article.NewSocials(profile1, profile2)

	ids := profiles.IDs()
	assert.Equal(t, 2, len(ids))
	assert.Contains(t, ids, profile1.ID)
	assert.Contains(t, ids, profile2.ID)
}

func TestSocialProfiles_Count(t *testing.T) {
	profiles := article.NewSocials()
	assert.Equal(t, 0, profiles.Len())

	profile1 := fakeSocialProfile()
	profile2 := fakeSocialProfile()
	profiles.Add(profile1, profile2)
	assert.Equal(t, 2, profiles.Len())
}

func TestSocialProfiles_Filter(t *testing.T) {
	profile1 := fakeSocialProfile()
	profile2 := fakeSocialProfile()
	profiles := article.NewSocials(profile1, profile2)

	t.Run("filter with always true function", func(t *testing.T) {
		filtered := profiles.Filter(func(profile *article.Social) bool {
			return true
		})
		assert.Equal(t, 2, filtered.Len())
	})

	t.Run("filter with always false function", func(t *testing.T) {
		filtered := profiles.Filter(func(profile *article.Social) bool {
			return false
		})
		assert.Equal(t, 0, filtered.Len())
	})

	t.Run("filter with specific condition", func(t *testing.T) {
		filtered := profiles.Filter(func(profile *article.Social) bool {
			return profile.ID == profile1.ID
		})
		assert.Equal(t, 1, filtered.Len())
		assert.Equal(t, profile1.ID, filtered.IDs()[0])
	})
}

func TestSocialProfiles_Normalize(t *testing.T) {
	valid := fakeSocialProfile()
	invalid := &article.Social{ID: gofakeit.UUID(), URL: " sdf "}

	profiles := article.NewSocials(valid, invalid)
	profiles.Normalize()

	// Normalize should remove invalid social profile
	_, exist := profiles.Get(invalid.ID)
	assert.False(t, exist)

	// Normalize should keep valid social profile
	_, exist = profiles.Get(valid.ID)
	assert.True(t, exist)
}
