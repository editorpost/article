package article

import (
	"encoding/json"
	"errors"
	"log"
)

// Socials represents a collection of Social pointers
type Socials struct {
	items []*Social
}

// NewSocials creates a collection, skips invalid social items, and logs errors
func NewSocials(profiles ...*Social) *Socials {
	validProfiles := []*Social{}
	for _, profile := range profiles {
		if err := validate.Struct(profile); err == nil {
			validProfiles = append(validProfiles, profile)
		} else {
			log.Printf("Invalid social profile skipped: %+v, error: %v", profile, err)
		}
	}
	return &Socials{items: validProfiles}
}

// NewSocialsStrict creates a collection and validates every social profile
func NewSocialsStrict(profiles ...*Social) (*Socials, error) {
	for _, profile := range profiles {
		if err := validate.Struct(profile); err != nil {
			return nil, errors.New("invalid social profile: " + err.Error())
		}
	}
	return &Socials{items: profiles}, nil
}

// Get returns the social profile by ID
func (list *Socials) Get(id string) (*Social, bool) {
	for _, profile := range list.items {
		if profile.ID == id {
			return profile, true
		}
	}
	return nil, false
}

// Slice returns a slice of all social items
func (list *Socials) Slice() []*Social {
	return list.items
}

// Add adds social items to the collection
func (list *Socials) Add(profiles ...*Social) *Socials {
	for _, profile := range profiles {
		if err := validate.Struct(profile); err == nil {
			list.items = append(list.items, profile)
		} else {
			log.Printf("Invalid social profile skipped: %+v, error: %v", profile, err)
		}
	}
	return list
}

// Remove removes social items by ID
func (list *Socials) Remove(ids ...string) *Socials {
	idSet := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}

	filteredProfiles := []*Social{}
	for _, profile := range list.items {
		if _, found := idSet[profile.ID]; !found {
			filteredProfiles = append(filteredProfiles, profile)
		}
	}

	list.items = filteredProfiles
	return list
}

// IDs returns a slice of all social profile IDs
func (list *Socials) IDs() []string {
	ids := make([]string, len(list.items))
	for idx, profile := range list.items {
		ids[idx] = profile.ID
	}
	return ids
}

// Len returns the number of social items
func (list *Socials) Len() int {
	return len(list.items)
}

// Filter returns a new Socials collection filtered by the provided functions
func (list *Socials) Filter(fns ...func(*Social) bool) *Socials {
	filteredProfiles := []*Social{}
	for _, profile := range list.items {
		include := true
		for _, fn := range fns {
			if !fn(profile) {
				include = false
				break
			}
		}
		if include {
			filteredProfiles = append(filteredProfiles, profile)
		}
	}
	return &Socials{items: filteredProfiles}
}

// Normalize validates and trims the fields of all social items
func (list *Socials) Normalize() {
	for _, profile := range list.items {
		profile.Normalize()
	}
}

// UnmarshalJSON to array of images using encoding/json
func (list *Socials) UnmarshalJSON(data []byte) error {

	// Unmarshal to a slice of Image
	var socials []*Social
	if err := json.Unmarshal(data, &socials); err != nil {
		return err
	}

	// Create a new Images collection
	*list = *NewSocials(socials...)

	return nil
}

// MarshalJSON from array of images using encoding/json
func (list *Socials) MarshalJSON() ([]byte, error) {
	return json.Marshal(list.items)
}
