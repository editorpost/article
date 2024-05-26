package article

import (
	"encoding/json"
	"log"
	"log/slog"
)

// Medias represents a collection of Media pointers
type Medias struct {
	items []*Media
}

// NewMediasStrict creates a collection and validates every media
func NewMediasStrict(medias ...*Media) (*Medias, error) {

	valid := make([]*Media, 0, len(medias))

	for _, media := range medias {

		if err := validate.Struct(media); err != nil {
			return nil, err
		}

		valid = append(valid, media)
	}

	return &Medias{items: valid}, nil
}

// NewMedias creates a collection, skips invalid items, and logs errors
func NewMedias(medias ...*Media) *Medias {

	var valid []*Media

	for _, media := range medias {
		if err := validate.Struct(media); err == nil {
			valid = append(valid, media)
		} else {
			slog.Debug("Invalid media skipped: %v", err)
		}
	}

	return &Medias{items: valid}
}

// Normalize validates and trims the fields of all items
func (list *Medias) Normalize() {
	for _, img := range list.items {
		img.Normalize()
	}
}

// Get returns the media by ID
func (list *Medias) Get(id string) (*Media, bool) {
	for _, img := range list.items {
		if img.ID == id {
			return img, true
		}
	}
	return nil, false
}

// Slice returns a slice of all items
func (list *Medias) Slice() []*Media {
	return list.items
}

// Add adds items to the collection
func (list *Medias) Add(medias ...*Media) *Medias {
	for _, img := range medias {
		if img != nil && img.ID != "" {
			list.items = append(list.items, img)
		} else {
			log.Printf("Invalid media skipped: %+v", img)
		}
	}
	return list
}

// Remove removes items by ID
func (list *Medias) Remove(ids ...string) *Medias {
	idSet := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}

	var filteredMedias []*Media
	for _, img := range list.items {
		if _, found := idSet[img.ID]; !found {
			filteredMedias = append(filteredMedias, img)
		}
	}

	list.items = filteredMedias
	return list
}

// IDs returns a slice of all media IDs
func (list *Medias) IDs() []string {
	ids := make([]string, len(list.items))
	for idx, img := range list.items {
		ids[idx] = img.ID
	}
	return ids
}

// Len returns the number of items
func (list *Medias) Len() int {
	return len(list.items)
}

// Filter returns a new Medias collection filtered by the provided functions
func (list *Medias) Filter(fns ...func(*Media) bool) *Medias {
	var filteredMedias []*Media
	for _, img := range list.items {
		include := true
		for _, fn := range fns {
			if !fn(img) {
				include = false
				break
			}
		}
		if include {
			filteredMedias = append(filteredMedias, img)
		}
	}
	return &Medias{items: filteredMedias}
}

// UnmarshalJSON to array of items using encoding/json
func (list *Medias) UnmarshalJSON(data []byte) error {

	// Unmarshal to a slice of Media
	var medias []*Media
	if err := json.Unmarshal(data, &medias); err != nil {
		return err
	}

	// Create a new Medias collection
	*list = *NewMedias(medias...)

	return nil
}

// MarshalJSON from array of items using encoding/json
func (list *Medias) MarshalJSON() ([]byte, error) {
	return json.Marshal(list.items)
}
