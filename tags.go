package article

import (
	"encoding/json"
	"strings"
)

type Tags struct {
	tags []string
}

func NewTags(tags ...string) *Tags {

	// skip empty tags
	validTags := []string{}

	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" && len(tag) <= 255 {
			validTags = append(validTags, tag)
		}
	}

	return &Tags{tags: validTags}
}

// String returns a comma-separated list of tags
func (list *Tags) String() string {
	return strings.Join(list.tags, ",")
}

// Add adds tags to the collection
func (list *Tags) Add(tags ...string) *Tags {
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" && len(tag) <= 255 {
			list.tags = append(list.tags, tag)
		}
	}
	return list
}

// Slice returns a slice of all tags
func (list *Tags) Slice() []string {
	return list.tags
}

// Contains returns true if the tag exists in the collection
func (list *Tags) Contains(tag string) bool {
	for _, match := range list.tags {
		if match == tag {
			return true
		}
	}
	return false
}

// Remove removes a tag from the collection
func (list *Tags) Remove(tag string) *Tags {
	filteredTags := []string{}
	for _, match := range list.tags {
		if match != tag {
			filteredTags = append(filteredTags, match)
		}
	}
	list.tags = filteredTags
	return list
}

// Len returns the number of tags
func (list *Tags) Len() int {
	return len(list.tags)
}

// UnmarshalJSON to array of images using encoding/json
func (list *Tags) UnmarshalJSON(data []byte) error {

	// Unmarshal to a slice of Image
	var tags []string
	if err := json.Unmarshal(data, &tags); err != nil {
		return err
	}

	// Create a new Images collection
	*list = *NewTags(tags...)

	return nil
}

// MarshalJSON from array of images using encoding/json
func (list Tags) MarshalJSON() ([]byte, error) {
	return json.Marshal(list.tags)
}
