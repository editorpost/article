package article

import (
	"encoding/json"
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
	var validVideos []*Video
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
func (list *Videos) Get(id string) (*Video, bool) {
	for _, video := range list.videos {
		if video.ID == id {
			return video, true
		}
	}
	return nil, false
}

// Slice returns a slice of all videos
func (list *Videos) Slice() []*Video {
	return list.videos
}

// Add adds videos to the collection
func (list *Videos) Add(videos ...*Video) *Videos {
	for _, video := range videos {
		if err := validate.Struct(video); err == nil {
			list.videos = append(list.videos, video)
		} else {
			log.Printf("Invalid video skipped: %+v, error: %v", video, err)
		}
	}
	return list
}

// Remove removes videos by ID
func (list *Videos) Remove(ids ...string) *Videos {
	idSet := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}

	var filteredVideos []*Video
	for _, video := range list.videos {
		if _, found := idSet[video.ID]; !found {
			filteredVideos = append(filteredVideos, video)
		}
	}

	list.videos = filteredVideos
	return list
}

// IDs returns a slice of all video IDs
func (list *Videos) IDs() []string {
	ids := make([]string, len(list.videos))
	for idx, video := range list.videos {
		ids[idx] = video.ID
	}
	return ids
}

// Len returns the number of videos
func (list *Videos) Len() int {
	return len(list.videos)
}

// Filter returns a new Videos collection filtered by the provided functions
func (list *Videos) Filter(fns ...func(*Video) bool) *Videos {
	var filteredVideos []*Video
	for _, video := range list.videos {
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
func (list *Videos) Normalize() {
	for _, video := range list.videos {
		video.Normalize()
	}
}

// UnmarshalJSON to array of items using encoding/json
func (list *Videos) UnmarshalJSON(data []byte) error {

	// Unmarshal to a slice of Image
	var videos []*Video
	if err := json.Unmarshal(data, &videos); err != nil {
		return err
	}

	// Create a new Images collection
	*list = *NewVideos(videos...)

	return nil
}

// MarshalJSON from array of items using encoding/json
func (list *Videos) MarshalJSON() ([]byte, error) {
	return json.Marshal(list.videos)
}
