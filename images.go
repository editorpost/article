package article

import (
	"encoding/json"
	"log"
	"log/slog"
)

// Images represents a collection of Image pointers
type Images struct {
	items []*Image
}

// NewImagesStrict creates a collection and validates every image
func NewImagesStrict(images ...*Image) (*Images, error) {

	var valid []*Image

	for _, image := range images {

		if err := validate.Struct(image); err != nil {
			return nil, err
		}

		valid = append(valid, image)
	}

	return &Images{items: valid}, nil
}

// NewImages creates a collection, skips invalid items, and logs errors
func NewImages(images ...*Image) *Images {

	var valid []*Image

	for _, image := range images {
		if err := validate.Struct(image); err == nil {
			valid = append(valid, image)
		} else {
			slog.Debug("Invalid image skipped: %v", err)
		}
	}

	return &Images{items: valid}
}

// Normalize validates and trims the fields of all items
func (list *Images) Normalize() {
	for _, img := range list.items {
		img.Normalize()
	}
}

// Get returns the image by ID
func (list *Images) Get(id string) (*Image, bool) {
	for _, img := range list.items {
		if img.ID == id {
			return img, true
		}
	}
	return nil, false
}

// Slice returns a slice of all items
func (list *Images) Slice() []*Image {
	return list.items
}

// Add adds items to the collection
func (list *Images) Add(images ...*Image) *Images {
	for _, img := range images {
		if img != nil && img.ID != "" {
			list.items = append(list.items, img)
		} else {
			log.Printf("Invalid image skipped: %+v", img)
		}
	}
	return list
}

// Remove removes items by ID
func (list *Images) Remove(ids ...string) *Images {

	if len(ids) == 0 {
		return list
	}

	idSet := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}

	var filteredImages []*Image
	for _, img := range list.items {
		if _, found := idSet[img.ID]; !found {
			filteredImages = append(filteredImages, img)
		}
	}

	list.items = filteredImages
	return list
}

// IDs returns a slice of all image IDs
func (list *Images) IDs() []string {
	ids := make([]string, len(list.items))
	for idx, img := range list.items {
		ids[idx] = img.ID
	}
	return ids
}

// Len returns the number of items
func (list *Images) Len() int {
	return len(list.items)
}

// Filter returns a new Images collection filtered by the provided functions
func (list *Images) Filter(fns ...func(*Image) bool) *Images {
	var filteredImages []*Image
	for _, img := range list.items {
		include := true
		for _, fn := range fns {
			if !fn(img) {
				include = false
				break
			}
		}
		if include {
			filteredImages = append(filteredImages, img)
		}
	}
	return &Images{items: filteredImages}
}

// ReplaceURLs replaces the URLs of the Images from the provided map.
// Returns a slice of image IDs that failed to be replaced.
func (list *Images) ReplaceURLs(m map[string]string) []string {

	failed := make([]string, 0)

	for _, img := range list.Slice() {

		url, exists := m[img.URL]

		if exists {
			img.URL = url
			continue
		}

		failed = append(failed, img.ID)
	}

	return failed
}

// ReplaceOrRemoveURLs replaces the URLs of the Images from the provided map.
// If an image URL is not found in the map, the image is removed from the Images collection.
func (list *Images) ReplaceOrRemoveURLs(m map[string]string) []string {
	failed := list.ReplaceURLs(m)
	list.Remove(failed...)
	return failed
}

// UnmarshalJSON to array of items using encoding/json
func (list *Images) UnmarshalJSON(data []byte) error {

	// Unmarshal to a slice of Image
	var images []*Image
	if err := json.Unmarshal(data, &images); err != nil {
		return err
	}

	// Create a new Images collection
	*list = *NewImages(images...)

	return nil
}

// MarshalJSON from array of items using encoding/json
func (list *Images) MarshalJSON() ([]byte, error) {
	return json.Marshal(list.items)
}
