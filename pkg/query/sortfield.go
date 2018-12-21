package query

import (
	"github.com/kucjac/jsonapi/pkg/internal"
	"github.com/kucjac/jsonapi/pkg/mapping"
	"strings"
)

// Order is an enumerator that describes the order of sorting
type Order int

const (
	AscendingOrder Order = iota
	DescendingOrder
)

// SortField is a field that describes the sorting rules for given
type SortField struct {
	*mapping.StructField

	// Order defines if the sorting order (ascending or descending)
	Order Order

	// SubFields is the relationship sub field sorts
	SubFields []*SortField
}

func (s *SortField) copy() *SortField {
	sort := &SortField{StructField: s.StructField, Order: s.Order}
	if len(s.SubFields) != 0 {
		sort.SubFields = make([]*SortField, len(s.SubFields))
		for i, v := range s.SubFields {
			sort.SubFields[i] = v.copy()
		}
	}
	return sort

}

var maxNestedRelLevel int = 1

func newSortField(sort string, order Order, scope *Scope) (invalidField bool) {
	var (
		sField    *mapping.StructField
		ok        bool
		sortField *SortField
	)

	splitted := strings.Split(sort, internal.AnnotationNestedSeperator)
	l := len(splitted)
	switch {
	// for length == 1 the sort must be an attribute or a primary field
	case l == 1:
		if sort == internal.AnnotationID {
			sField = scope.Struct.Primary()
		} else {
			sField, ok = scope.Struct.Attr(sort)
			if !ok {
				invalidField = true
				return
			}
		}
		sortField = &SortField{StructField: sField, Order: order}
		scope.Sorts = append(scope.Sorts, sortField)
	case l <= (maxNestedRelLevel + 1):
		sField, ok = scope.Struct.RelationField(splitted[0])
		if !ok {

			invalidField = true
			return
		}
		// if true then the nested should be an attribute for given
		var found bool
		for i := range scope.Sorts {
			if scope.Sorts[i].FieldIndex() == sField.FieldIndex() {
				sortField = scope.Sorts[i]
				found = true
				break
			}
		}
		if !found {
			sortField = &SortField{StructField: sField}
		}
		invalidField = sortField.setSubfield(splitted[1:], order)
		if !found && !invalidField {
			scope.Sorts = append(scope.Sorts, sortField)
		}
	default:
		invalidField = true
	}
	return
}

// setSubfield sets sortfield for subfield of given relationship field.
func (s *SortField) setSubfield(sortSplitted []string, order Order) (invalidField bool) {
	var (
		subField *SortField
		sField   *mapping.StructField
	)

	// Subfields are available only for the relationships
	if !s.IsRelationship() {
		invalidField = true
		return
	}

	// sort splitted is splitted sort query entry
	// i.e. a sort query for
	switch len(sortSplitted) {
	case 0:
		invalidField = true
		return
	case 1:
		// if len is equal to one then it should be primary or attribute field
		sort := sortSplitted[0]
		if sort == internal.AnnotationID {
			sField = s.Relationship().ModelStruct().Primary()
		} else {
			var ok bool
			sField, ok = s.Relationship().ModelStruct().Attr(sortSplitted[0])
			if !ok {
				invalidField = true
				return
			}
		}

		s.SubFields = append(s.SubFields, &SortField{StructField: sField, Order: order})
	default:
		// if length is more than one -> there is a relationship
		var ok bool
		sField, ok := s.Relationship().ModelStruct().RelationField(sortSplitted[0])
		if !ok {
			invalidField = true
			return
		}

		// search for the subfields if already created
		for i := range s.SubFields {
			if s.SubFields[i].FieldIndex() == sField.FieldIndex() {
				subField = s.SubFields[i]
				break
			}
		}

		// if none found create new
		if subField == nil {
			subField = &SortField{StructField: sField}
		}

		//
		invalidField = subField.setSubfield(sortSplitted[1:], order)
		if !invalidField {
			// if found keep the subfield in subfields
			s.SubFields = append(s.SubFields, subField)
		}
	}
	return
}