package article

import (
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
func (i *Images) Normalize() {
	for _, img := range i.images {
		img.Normalize()
	}
}

// Get returns the image by ID
func (i *Images) Get(id string) (*Image, bool) {
	for _, img := range i.images {
		if img.ID == id {
			return img, true
		}
	}
	return nil, false
}

// Slice returns a slice of all images
func (i *Images) Slice() []*Image {
	return i.images
}

// Add adds images to the collection
func (i *Images) Add(images ...*Image) *Images {
	for _, img := range images {
		if img != nil && img.ID != "" {
			i.images = append(i.images, img)
		} else {
			log.Printf("Invalid image skipped: %+v", img)
		}
	}
	return i
}

// Remove removes images by ID
func (i *Images) Remove(ids ...string) *Images {
	idSet := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}

	filteredImages := []*Image{}
	for _, img := range i.images {
		if _, found := idSet[img.ID]; !found {
			filteredImages = append(filteredImages, img)
		}
	}

	i.images = filteredImages
	return i
}

// IDs returns a slice of all image IDs
func (i *Images) IDs() []string {
	ids := make([]string, len(i.images))
	for idx, img := range i.images {
		ids[idx] = img.ID
	}
	return ids
}

// Len returns the number of images
func (i *Images) Len() int {
	return len(i.images)
}

// Filter returns a new Images collection filtered by the provided functions
func (i *Images) Filter(fns ...func(*Image) bool) *Images {
	filteredImages := []*Image{}
	for _, img := range i.images {
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
