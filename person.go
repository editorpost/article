package article

import (
	"github.com/google/uuid"
	"log/slog"
)

type Person struct {
	ID      string   `json:"id" validate:"required,max=36"`
	Name    string   `json:"name" validate:"required,max=255"`
	Role    string   `json:"role" validate:"max=255"`
	Images  *Images  `json:"images"`
	Socials *Socials `json:"socials"`
}

// NewPerson creates a new Person with a random ID.
func NewPerson(name string) *Person {
	return &Person{
		ID:   uuid.New().String(),
		Name: name,
	}
}

// Normalize validates and trims the fields of the Person.
func (p *Person) Normalize() {

	if p.ID == "" {
		p.ID = uuid.New().String()
	}

	p.ID = TrimToMaxLen(p.ID, 36)
	p.Name = TrimToMaxLen(p.Name, 255)
	p.Role = TrimToMaxLen(p.Role, 255)

	if p.Socials != nil {
		for _, social := range p.Socials.Slice() {
			social.Normalize()
		}
	}
}

// Map converts the Person struct to a map[string]any.
func (p *Person) Map() map[string]any {

	images := make([]map[string]any, p.Images.Len())
	for i, image := range p.Images.Slice() {
		images[i] = image.Map()
	}

	socials := make([]map[string]any, p.Socials.Len())
	for i, social := range p.Socials.Slice() {
		socials[i] = social.Map()
	}

	return map[string]any{
		"id":      p.ID,
		"name":    p.Name,
		"role":    p.Role,
		"images":  images,
		"socials": socials,
	}
}

type Persons struct {
	items []*Person
}

// NewPersons creates a new Persons collection.
func NewPersons(persons ...*Person) *Persons {

	valid := make([]*Person, 0, len(persons))
	for _, person := range persons {
		if err := validate.Struct(person); err == nil {
			valid = append(valid, person)
		} else {
			slog.Debug("Invalid person skipped", slog.String("error", err.Error()))
		}
	}

	return &Persons{items: valid}
}

// Get returns the Person by ID.
func (list *Persons) Get(id string) (*Person, bool) {
	for _, person := range list.items {
		if person.ID == id {
			return person, true
		}
	}
	return nil, false
}

// Slice returns a slice of all Persons.
func (list *Persons) Slice() []*Person {
	return list.items
}

// Add adds Persons to the collection.
func (list *Persons) Add(persons ...*Person) *Persons {
	for _, person := range persons {
		if err := validate.Struct(person); err == nil {
			list.items = append(list.items, person)
		} else {
			slog.Debug("Invalid person skipped", slog.String("error", err.Error()))
		}
	}
	return list
}
