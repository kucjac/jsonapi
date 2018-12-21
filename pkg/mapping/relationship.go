package mapping

import (
	"github.com/kucjac/jsonapi/pkg/internal/models"
)

type RelationshipKind models.RelationshipKind

const (
	RelUnknown RelationshipKind = iota
	RelBelongsTo
	RelHasOne
	RelHasMany
	RelMany2Many
)

func (r RelationshipKind) String() string {
	switch r {
	case RelUnknown:
	case RelBelongsTo:
		return "BelongsTo"
	case RelHasOne:
		return "HasOne"
	case RelHasMany:
		return "HasMany"
	case RelMany2Many:
		return "Many2Many"
	}
	return "Unknown"
}

type RelationshipOption models.RelationshipOption

const (
	Restrict RelationshipOption = iota
	NoAction
	Cascade
	SetNull
)

// By default the relationship many2many is local
// it can be synced with back-referenced relationship using the 'sync' flag
// The HasOne and HasMany is the relationship that by default is 'synced'
// The BelongsTo relationship is local by default.

type Relationship struct {
	*models.Relationship
}

// ModelStruct returns relationships model struct - related model structure
func (r *Relationship) ModelStruct() *ModelStruct {
	mStruct := models.RelationshipMStruct(r.Relationship)
	if mStruct == nil {
		return nil
	}
	return &ModelStruct{mStruct}
}

// Kind returns relationship Kind
func (r *Relationship) Kind() RelationshipKind {
	k := models.RelationshipGetKind(r.Relationship)
	return RelationshipKind(k)
}

// ForeignKey returns foreign key for given relationship
func (r *Relationship) ForeignKey() *StructField {
	fk := models.RelationshipForeignKey(r.Relationship)
	if fk == nil {
		return nil
	}
	return &StructField{fk}
}
