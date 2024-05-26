package article

import (
	"encoding/json"
	"errors"
	"log"
)

// Images represents a collection of Image pointers
type Images struct {
	images []*Image
}

// NewImagesStrict creates a collection and validates every image
func NewImagesStrict(images ...*Image) (*Images, error) {
	for _, img := range images {
		if img == nil || img.ID == "" {
			return nil, errors.New("invalid image: must not be nil and must have an ID")
		}
	}
	return &Images{images: images}, nil
}

// NewImages creates a collection, skips invalid images, and logs errors
func NewImages(images ...*Image) *Images {
	validImages := []*Image{}
	for _, img := range images {
		if img != nil && img.ID != "" {
			validImages = append(validImages, img)
		} else {
			log.Printf("Invalid image skipped: %+v", img)
		}
	}
	return &Images{images: validImages}
}

// Normalize validates and trims the fields of all images
func (list *Images) Normalize() {
	for _, img := range list.images {
		img.Normalize()
	}
}

// Get returns the image by ID
func (list *Images) Get(id string) (*Image, bool) {
	for _, img := range list.images {
		if img.ID == id {
			return img, true
		}
	}
	return nil, false
}

// Slice returns a slice of all images
func (list *Images) Slice() []*Image {
	return list.images
}

// Add adds images to the collection
func (list *Images) Add(images ...*Image) *Images {
	for _, img := range images {
		if img != nil && img.ID != "" {
			list.images = append(list.images, img)
		} else {
			log.Printf("Invalid image skipped: %+v", img)
		}
	}
	return list
}

// Remove removes images by ID
func (list *Images) Remove(ids ...string) *Images {
	idSet := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}

	filteredImages := []*Image{}
	for _, img := range list.images {
		if _, found := idSet[img.ID]; !found {
			filteredImages = append(filteredImages, img)
		}
	}

	list.images = filteredImages
	return list
}

// IDs returns a slice of all image IDs
func (list *Images) IDs() []string {
	ids := make([]string, len(list.images))
	for idx, img := range list.images {
		ids[idx] = img.ID
	}
	return ids
}

// Len returns the number of images
func (list *Images) Len() int {
	return len(list.images)
}

// Filter returns a new Images collection filtered by the provided functions
func (list *Images) Filter(fns ...func(*Image) bool) *Images {
	filteredImages := []*Image{}
	for _, img := range list.images {
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
	return &Images{images: filteredImages}
}

// UnmarshalJSON to array of images using encoding/json
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

// MarshalJSON from array of images using encoding/json
func (list Images) MarshalJSON() ([]byte, error) {
	return json.Marshal(list.images)
}
