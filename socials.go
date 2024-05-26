package article

import (
	"errors"
	"log"
)

// SocialProfiles represents a collection of SocialProfile pointers
type SocialProfiles struct {
	profiles []*SocialProfile
}

// NewSocialProfiles creates a collection, skips invalid social profiles, and logs errors
func NewSocialProfiles(profiles ...*SocialProfile) *SocialProfiles {
	validProfiles := []*SocialProfile{}
	for _, profile := range profiles {
		if err := validate.Struct(profile); err == nil {
			validProfiles = append(validProfiles, profile)
		} else {
			log.Printf("Invalid social profile skipped: %+v, error: %v", profile, err)
		}
	}
	return &SocialProfiles{profiles: validProfiles}
}

// NewSocialProfilesStrict creates a collection and validates every social profile
func NewSocialProfilesStrict(profiles ...*SocialProfile) (*SocialProfiles, error) {
	for _, profile := range profiles {
		if err := validate.Struct(profile); err != nil {
			return nil, errors.New("invalid social profile: " + err.Error())
		}
	}
	return &SocialProfiles{profiles: profiles}, nil
}

// Get returns the social profile by ID
func (sp *SocialProfiles) Get(id string) (*SocialProfile, bool) {
	for _, profile := range sp.profiles {
		if profile.ID == id {
			return profile, true
		}
	}
	return nil, false
}

// Slice returns a slice of all social profiles
func (sp *SocialProfiles) Slice() []*SocialProfile {
	return sp.profiles
}

// Add adds social profiles to the collection
func (sp *SocialProfiles) Add(profiles ...*SocialProfile) *SocialProfiles {
	for _, profile := range profiles {
		if err := validate.Struct(profile); err == nil {
			sp.profiles = append(sp.profiles, profile)
		} else {
			log.Printf("Invalid social profile skipped: %+v, error: %v", profile, err)
		}
	}
	return sp
}

// Remove removes social profiles by ID
func (sp *SocialProfiles) Remove(ids ...string) *SocialProfiles {
	idSet := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}

	filteredProfiles := []*SocialProfile{}
	for _, profile := range sp.profiles {
		if _, found := idSet[profile.ID]; !found {
			filteredProfiles = append(filteredProfiles, profile)
		}
	}

	sp.profiles = filteredProfiles
	return sp
}

// IDs returns a slice of all social profile IDs
func (sp *SocialProfiles) IDs() []string {
	ids := make([]string, len(sp.profiles))
	for idx, profile := range sp.profiles {
		ids[idx] = profile.ID
	}
	return ids
}

// Len returns the number of social profiles
func (sp *SocialProfiles) Len() int {
	return len(sp.profiles)
}

// Filter returns a new SocialProfiles collection filtered by the provided functions
func (sp *SocialProfiles) Filter(fns ...func(*SocialProfile) bool) *SocialProfiles {
	filteredProfiles := []*SocialProfile{}
	for _, profile := range sp.profiles {
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
	return &SocialProfiles{profiles: filteredProfiles}
}

// Normalize validates and trims the fields of all social profiles
func (sp *SocialProfiles) Normalize() {
	for _, profile := range sp.profiles {
		profile.Normalize()
	}
}
