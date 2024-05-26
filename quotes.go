package article

import (
	"errors"
	"log"
)

// Quotes represents a collection of Quote pointers
type Quotes struct {
	quotes []*Quote
}

// NewQuotesStrict creates a collection and validates every quote
func NewQuotesStrict(quotes ...*Quote) (*Quotes, error) {
	for _, quote := range quotes {
		if err := validate.Struct(quote); err != nil {
			return nil, errors.New("invalid quote: " + err.Error())
		}
	}
	return &Quotes{quotes: quotes}, nil
}

// NewQuotes creates a collection, skips invalid quotes, and logs errors
func NewQuotes(quotes ...*Quote) *Quotes {
	validQuotes := []*Quote{}
	for _, quote := range quotes {
		if err := validate.Struct(quote); err == nil {
			validQuotes = append(validQuotes, quote)
		} else {
			log.Printf("Invalid quote skipped: %+v, error: %v", quote, err)
		}
	}
	return &Quotes{quotes: validQuotes}
}

// Get returns the quote by ID
func (q *Quotes) Get(id string) (*Quote, bool) {
	for _, quote := range q.quotes {
		if quote.ID == id {
			return quote, true
		}
	}
	return nil, false
}

// Slice returns a slice of all quotes
func (q *Quotes) Slice() []*Quote {
	return q.quotes
}

// Add adds quotes to the collection
func (q *Quotes) Add(quotes ...*Quote) *Quotes {
	for _, quote := range quotes {
		if err := validate.Struct(quote); err == nil {
			q.quotes = append(q.quotes, quote)
		} else {
			log.Printf("Invalid quote skipped: %+v, error: %v", quote, err)
		}
	}
	return q
}

// Remove removes quotes by ID
func (q *Quotes) Remove(ids ...string) *Quotes {
	idSet := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}

	filteredQuotes := []*Quote{}
	for _, quote := range q.quotes {
		if _, found := idSet[quote.ID]; !found {
			filteredQuotes = append(filteredQuotes, quote)
		}
	}

	q.quotes = filteredQuotes
	return q
}

// IDs returns a slice of all quote IDs
func (q *Quotes) IDs() []string {
	ids := make([]string, len(q.quotes))
	for idx, quote := range q.quotes {
		ids[idx] = quote.ID
	}
	return ids
}

// Len returns the number of quotes
func (q *Quotes) Len() int {
	return len(q.quotes)
}

// Filter returns a new Quotes collection filtered by the provided functions
func (q *Quotes) Filter(fns ...func(*Quote) bool) *Quotes {
	filteredQuotes := []*Quote{}
	for _, quote := range q.quotes {
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
	return &Quotes{quotes: filteredQuotes}
}

// Normalize validates and trims the fields of all quotes
func (q *Quotes) Normalize() {
	for _, quote := range q.quotes {
		quote.Normalize()
	}
}
