package article

import (
	"encoding/json"
	"errors"
	"log"
)

// Quotes represents a collection of Quote pointers
type Quotes struct {
	items []*Quote
}

// NewQuotesStrict creates a collection and validates every quote
func NewQuotesStrict(quotes ...*Quote) (*Quotes, error) {
	for _, quote := range quotes {
		if err := validate.Struct(quote); err != nil {
			return nil, errors.New("invalid quote: " + err.Error())
		}
	}
	return &Quotes{items: quotes}, nil
}

// NewQuotes creates a collection, skips invalid items, and logs errors
func NewQuotes(quotes ...*Quote) *Quotes {
	validQuotes := []*Quote{}
	for _, quote := range quotes {
		if err := validate.Struct(quote); err == nil {
			validQuotes = append(validQuotes, quote)
		} else {
			log.Printf("Invalid quote skipped: %+v, error: %v", quote, err)
		}
	}
	return &Quotes{items: validQuotes}
}

// Get returns the quote by ID
func (list *Quotes) Get(id string) (*Quote, bool) {
	for _, quote := range list.items {
		if quote.ID == id {
			return quote, true
		}
	}
	return nil, false
}

// Slice returns a slice of all items
func (list *Quotes) Slice() []*Quote {
	return list.items
}

// Add adds items to the collection
func (list *Quotes) Add(quotes ...*Quote) *Quotes {
	for _, quote := range quotes {
		if err := validate.Struct(quote); err == nil {
			list.items = append(list.items, quote)
		} else {
			log.Printf("Invalid quote skipped: %+v, error: %v", quote, err)
		}
	}
	return list
}

// Remove removes items by ID
func (list *Quotes) Remove(ids ...string) *Quotes {
	idSet := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}

	filteredQuotes := []*Quote{}
	for _, quote := range list.items {
		if _, found := idSet[quote.ID]; !found {
			filteredQuotes = append(filteredQuotes, quote)
		}
	}

	list.items = filteredQuotes
	return list
}

// IDs returns a slice of all quote IDs
func (list *Quotes) IDs() []string {
	ids := make([]string, len(list.items))
	for idx, quote := range list.items {
		ids[idx] = quote.ID
	}
	return ids
}

// Len returns the number of items
func (list *Quotes) Len() int {
	return len(list.items)
}

// Filter returns a new Quotes collection filtered by the provided functions
func (list *Quotes) Filter(fns ...func(*Quote) bool) *Quotes {
	filteredQuotes := []*Quote{}
	for _, quote := range list.items {
		include := true
		for _, fn := range fns {
			if !fn(quote) {
				include = false
				break
			}
		}
		if include {
			filteredQuotes = append(filteredQuotes, quote)
		}
	}
	return &Quotes{items: filteredQuotes}
}

// Normalize validates and trims the fields of all items
func (list *Quotes) Normalize() {
	for _, quote := range list.items {
		quote.Normalize()
	}
}

// UnmarshalJSON to array of images using encoding/json
func (list *Quotes) UnmarshalJSON(data []byte) error {

	// Unmarshal to a slice of Image
	var quotes []*Quote
	if err := json.Unmarshal(data, &quotes); err != nil {
		return err
	}

	// Create a new Images collection
	*list = *NewQuotes(quotes...)

	return nil
}

// MarshalJSON from array of images using encoding/json
func (list *Quotes) MarshalJSON() ([]byte, error) {
	return json.Marshal(list.items)
}
