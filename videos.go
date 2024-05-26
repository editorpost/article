package article

import (
	"errors"
	"log"
)

// Videos represents a collection of Video pointers
type Videos struct {
	videos []*Video
}

// NewVideosStrict creates a collection and validates every video
func NewVideosStrict(videos ...*Video) (*Videos, error) {
	for _, video := range videos {
		if err := validate.Struct(video); err != nil {
			return nil, errors.New("invalid video: " + err.Error())
		}
	}
	return &Videos{videos: videos}, nil
}

// NewVideos creates a collection, skips invalid videos, and logs errors
func NewVideos(videos ...*Video) *Videos {
	validVideos := []*Video{}
	for _, video := range videos {
		if err := validate.Struct(video); err == nil {
			validVideos = append(validVideos, video)
		} else {
			log.Printf("Invalid video skipped: %+v, error: %v", video, err)
		}
	}
	return &Videos{videos: validVideos}
}

// Get returns the video by ID
func (v *Videos) Get(id string) (*Video, bool) {
	for _, video := range v.videos {
		if video.ID == id {
			return video, true
		}
	}
	return nil, false
}

// Add adds videos to the collection
func (v *Videos) Add(videos ...*Video) *Videos {
	for _, video := range videos {
		if err := validate.Struct(video); err == nil {
			v.videos = append(v.videos, video)
		} else {
			log.Printf("Invalid video skipped: %+v, error: %v", video, err)
		}
	}
	return v
}

// Remove removes videos by ID
func (v *Videos) Remove(ids ...string) *Videos {
	idSet := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}

	filteredVideos := []*Video{}
	for _, video := range v.videos {
		if _, found := idSet[video.ID]; !found {
			filteredVideos = append(filteredVideos, video)
		}
	}

	v.videos = filteredVideos
	return v
}

// IDs returns a slice of all video IDs
func (v *Videos) IDs() []string {
	ids := make([]string, len(v.videos))
	for idx, video := range v.videos {
		ids[idx] = video.ID
	}
	return ids
}

// Count returns the number of videos
func (v *Videos) Count() int {
	return len(v.videos)
}

// Filter returns a new Videos collection filtered by the provided functions
func (v *Videos) Filter(fns ...func(*Video) bool) *Videos {
	filteredVideos := []*Video{}
	for _, video := range v.videos {
		include := true
		for _, fn := range fns {
			if !fn(video) {
				include = false
				break
			}
		}
		if include {
			filteredVideos = append(filteredVideos, video)
		}
	}
	return &Videos{videos: filteredVideos}
}

// Normalize validates and trims the fields of all videos
func (v *Videos) Normalize() {
	for _, video := range v.videos {
		video.Normalize()
	}
}
