package article

import "strings"

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

// NewTagsFromString creates a collection of tags from a comma-separated string
func NewTagsFromString(str string) *Tags {
	return NewTags(strings.Split(str, ",")...)
}

// String returns a comma-separated list of tags
func (t *Tags) String() string {
	return strings.Join(t.tags, ",")
}

// Add adds tags to the collection
func (t *Tags) Add(tags ...string) *Tags {
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" && len(tag) <= 255 {
			t.tags = append(t.tags, tag)
		}
	}
	return t
}

// Slice returns a slice of all tags
func (t *Tags) Slice() []string {
	return t.tags
}

// Contains returns true if the tag exists in the collection
func (t *Tags) Contains(tag string) bool {
	for _, match := range t.tags {
		if match == tag {
			return true
		}
	}
	return false
}

// Remove removes a tag from the collection
func (t *Tags) Remove(tag string) *Tags {
	filteredTags := []string{}
	for _, match := range t.tags {
		if match != tag {
			filteredTags = append(filteredTags, match)
		}
	}
	t.tags = filteredTags
	return t
}

// Len returns the number of tags
func (t *Tags) Len() int {
	return len(t.tags)
}
