package query

import (
	"fmt"
	"github.com/kucjac/jsonapi/pkg/flags"
	"github.com/kucjac/jsonapi/pkg/internal"
	"github.com/kucjac/jsonapi/pkg/mapping"
	"reflect"
)

// IncludeScope is the includes information scope
// it contains the field to include from the root scope
// related subscope, and subfields to include.
type IncludeField struct {
	*mapping.StructField

	// Scope is the scope that contains the values and filters for given
	// include field
	Scope *Scope

	// RelatedScope is a pointer to the scope where the IncludedField is stored.
	RelatedScope *Scope

	NotInFieldset bool
}

func (i *IncludeField) copy(relatedScope *Scope, root *Scope) *IncludeField {
	included := &IncludeField{StructField: i.StructField, NotInFieldset: i.NotInFieldset}
	included.Scope = i.Scope.copy(false, root)
	included.RelatedScope = relatedScope
	return included
}

// GetMissingPrimaries gets the id values from the RelatedScope, checks which id values were
// already stored within the colleciton root scope and return new ones.
func (i *IncludeField) GetMissingPrimaries() ([]interface{}, error) {
	return i.getMissingPrimaries()
}

func newIncludeField(field *mapping.StructField, scope *Scope) *IncludeField {
	includeField := new(IncludeField)
	includeField.StructField = field

	// Set NewScope for given field

	includeField.Scope = scope.createModelsScope(field.Relationship().ModelStruct())

	// Set the root collection scope for given scope
	includeField.Scope.collectionScope = scope.getOrCreateModelsRootScope(field.Relationship().ModelStruct())
	if _, ok := includeField.Scope.collectionScope.Fieldset[includeField.ApiName()]; !ok {
		includeField.NotInFieldset = true
		scope.hasFieldNotInFieldset = true
	}

	// Set relatedScope for given incldudedField
	includeField.RelatedScope = scope

	includeField.Scope.rootScope.totalIncludeCount++

	return includeField
}

func (i *IncludeField) getMissingPrimaries() ([]interface{}, error) {
	// uniqueMissing makes it possible to get unique ids that are not already used
	uniqueMissing := map[interface{}]struct{}{}

	// Lock the SafeHashMap for given collection
	i.Scope.collectionScope.IncludedValues.Lock()
	defer i.Scope.collectionScope.IncludedValues.Unlock()

	// Get the value from the RelatedScope
	v := reflect.ValueOf(i.RelatedScope.Value)

	switch v.Kind() {
	case reflect.Slice:
		for j := 0; j < v.Len(); j++ {
			elem := v.Index(j)
			if elem.IsNil() {
				continue
			}

			if err := i.getMissingFromSingle(elem, uniqueMissing); err != nil {
				return nil, err
			}

		}
	case reflect.Ptr:
		if err := i.getMissingFromSingle(v, uniqueMissing); err != nil {
			return nil, err
		}
	default:
		err := internal.IErrUnexpectedType
		fmt.Println(v)
		return nil, err
	}

	// Copy the notUsed map into array
	missingIDs := make([]interface{}, len(uniqueMissing))

	j := 0
	for uniqueID := range uniqueMissing {
		missingIDs[j] = uniqueID
		j++
	}

	return missingIDs, nil
}

func (i *IncludeField) getMissingFromSingle(
	value reflect.Value,
	uniqueMissing map[interface{}]struct{},
) error {

	var (
		fieldValue = value.Elem().Field(i.FieldIndex())
		primIndex  = i.Relationship().ModelStruct().Primary().FieldIndex()

		setCollectionValues = func(ptr reflect.Value) {
			primValue := ptr.Elem().Field(primIndex)

			if primValue.IsValid() {
				primary := primValue.Interface()

				if _, ok := i.Scope.collectionScope.IncludedValues.Values()[primary]; !ok {
					// add to collection IDs
					i.Scope.collectionScope.IncludedValues.Add(primary, nil)
					if _, ok = uniqueMissing[primary]; !ok {
						uniqueMissing[primary] = struct{}{}
					} else {
						fmt.Println("Already exists")
					}
				}
			}
		}
	)

	// Get the type of the value
	switch fieldValue.Kind() {
	case reflect.Slice:
		for j := 0; j < fieldValue.Len(); j++ {
			// set primary field within scope for given model struct
			elem := fieldValue.Index(j)
			if !elem.IsNil() {
				setCollectionValues(elem)
			}
		}
	case reflect.Ptr:
		if !fieldValue.IsNil() {

			primValue := fieldValue.Elem().Field(primIndex)
			if primValue.IsValid() {
				setCollectionValues(fieldValue)
			}
		}
	default:
		err := internal.IErrUnexpectedType
		return err
	}

	return nil
}

func (i *IncludeField) setRelationshipValue(relatedValue reflect.Value) {
	var includedScopeValue reflect.Value
	fieldValue := relatedValue.Field(i.FieldIndex())

	switch i.FieldType().Kind() {
	case reflect.Slice:
		if fieldValue.Len() == 0 {
			i.Scope.Value = reflect.New(i.FieldType()).Elem().Interface()
			return
		}
		includedScopeValue = reflect.New(i.FieldType()).Elem()
		includedScopeValue.Set(fieldValue)
	case reflect.Ptr:
		if fieldValue.IsNil() {
			i.Scope.Value = nil
			return
		}
		includedScopeValue = reflect.New(i.FieldType().Elem())
		includedScopeValue.Elem().Set(fieldValue.Elem())
	}

	i.Scope.Value = includedScopeValue.Interface()
	return
}

func (i *IncludeField) copyScopeBoundaries() {
	// copy primaries
	i.Scope.PrimaryFilters = make([]*FilterField, len(i.Scope.collectionScope.PrimaryFilters))
	copy(i.Scope.PrimaryFilters, i.Scope.collectionScope.PrimaryFilters)

	// copy attribute filters
	i.Scope.AttributeFilters = make([]*FilterField, len(i.Scope.collectionScope.AttributeFilters))
	copy(i.Scope.AttributeFilters, i.Scope.collectionScope.AttributeFilters)

	// copy filterKeyFilters
	i.Scope.FilterKeyFilters = make([]*FilterField, len(i.Scope.collectionScope.FilterKeyFilters))
	copy(i.Scope.FilterKeyFilters, i.Scope.collectionScope.FilterKeyFilters)

	// relationships
	i.Scope.RelationshipFilters = make([]*FilterField, len(i.Scope.collectionScope.RelationshipFilters))
	copy(i.Scope.RelationshipFilters, i.Scope.collectionScope.RelationshipFilters)

	//copy foreignKeyFilters
	i.Scope.ForeignKeyFilters = make([]*FilterField, len(i.Scope.collectionScope.ForeignKeyFilters))
	copy(i.Scope.ForeignKeyFilters, i.Scope.collectionScope.ForeignKeyFilters)

	// copy language filters

	if i.Scope.collectionScope.LanguageFilters != nil {
		i.Scope.LanguageFilters = i.Scope.collectionScope.LanguageFilters
	}

	// fieldset is taken by reference - copied if there is nested
	i.Scope.Fieldset = i.Scope.collectionScope.Fieldset
	i.Scope.ctx = i.Scope.collectionScope.ctx
	i.Scope.logger = i.Scope.collectionScope.logger

	if f, ok := i.Scope.collectionScope.Flags().Get(flags.UseLinks); ok {
		i.Scope.Flags().Set(flags.UseLinks, f)
	}

	if f, ok := i.Scope.collectionScope.Flags().Get(flags.ReturnPatchContent); ok {
		i.Scope.Flags().Set(flags.ReturnPatchContent, f)
	}

	for _, nested := range i.Scope.IncludedFields {
		// if the nested include is not found within the collection fieldset
		// the 'i'.Scope should have a new (not reference) Fieldset
		// with the nested field added to it
		if nested.NotInFieldset {
			// make a new fieldset if it is the same reference
			if len(i.Scope.Fieldset) == len(i.Scope.collectionScope.Fieldset) {
				// if there is more than one nested this would not happen
				i.Scope.Fieldset = make(map[string]*mapping.StructField)
				// copy fieldset
				for key, field := range i.Scope.collectionScope.Fieldset {
					i.Scope.Fieldset[key] = field
				}
			}

			//add nested
			i.Scope.Fieldset[nested.ApiName()] = nested.StructField
		}

		nested.copyScopeBoundaries()
	}

}

func (i *IncludeField) copyPresetFullParameters() {
	// copy primaries
	i.Scope.PrimaryFilters = make([]*FilterField, len(i.Scope.collectionScope.PrimaryFilters))
	copy(i.Scope.PrimaryFilters, i.Scope.collectionScope.PrimaryFilters)

	// copy attribute filters
	i.Scope.AttributeFilters = make([]*FilterField, len(i.Scope.collectionScope.AttributeFilters))
	copy(i.Scope.AttributeFilters, i.Scope.collectionScope.AttributeFilters)

	// copy filterKeyFilters
	i.Scope.FilterKeyFilters = make([]*FilterField, len(i.Scope.collectionScope.FilterKeyFilters))
	copy(i.Scope.FilterKeyFilters, i.Scope.collectionScope.FilterKeyFilters)

	// relationships
	i.Scope.RelationshipFilters = make([]*FilterField, len(i.Scope.collectionScope.RelationshipFilters))
	copy(i.Scope.RelationshipFilters, i.Scope.collectionScope.RelationshipFilters)

	//copy foreignKeyFilters
	i.Scope.ForeignKeyFilters = make([]*FilterField, len(i.Scope.collectionScope.ForeignKeyFilters))
	copy(i.Scope.ForeignKeyFilters, i.Scope.collectionScope.ForeignKeyFilters)

	// copy language filters
	if i.Scope.collectionScope.LanguageFilters != nil {
		i.Scope.LanguageFilters = i.Scope.collectionScope.LanguageFilters
	}

	// fieldset is taken by reference - copied if there is nested
	// i.Scope.Fieldset = i.Scope.collectionScope.Fieldset

	i.Scope.Sorts = make([]*SortField, len(i.Scope.collectionScope.Sorts))
	copy(i.Scope.Sorts, i.Scope.collectionScope.Sorts)

	i.Scope.Pagination = i.Scope.collectionScope.Pagination
	if f, ok := i.Scope.collectionScope.Flags().Get(flags.UseLinks); ok {
		i.Scope.Flags().Set(flags.UseLinks, f)
	}

	if f, ok := i.Scope.collectionScope.Flags().Get(flags.ReturnPatchContent); ok {
		i.Scope.Flags().Set(flags.ReturnPatchContent, f)
	}

	for _, nested := range i.Scope.IncludedFields {
		// if the nested include is not found within the collection fieldset
		// the 'i'.Scope should have a new (not reference) Fieldset
		// with the nested field added to it
		if nested.NotInFieldset {
			// make a new fieldset if it is the same reference
			if len(i.Scope.Fieldset) == len(i.Scope.collectionScope.Fieldset) {
				// if there is more than one nested this would not happen
				i.Scope.Fieldset = make(map[string]*mapping.StructField)
				// copy fieldset
				for key, field := range i.Scope.collectionScope.Fieldset {
					i.Scope.Fieldset[key] = field
				}
			}

			//add nested
			i.Scope.Fieldset[nested.ApiName()] = nested.StructField
		}

		nested.copyPresetFullParameters()
	}
}