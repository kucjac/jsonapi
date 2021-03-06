package mapping

type fieldFlag int

func (f fieldFlag) containsFlag(other fieldFlag) bool {
	return f&other != 0
}

// Field flags
const (
	// fDefault is a default flag value
	fDefault fieldFlag = 1 << iota
	// fOmitEmpty is a field flag for omitting empty value.
	fOmitempty
	// fI18n is the i18n field flag.
	fI18n
	// fNoFilter is the 'no filter' field flag.
	fNoFilter
	// fLanguage is the language field flag.
	fLanguage
	// fSortable is a flag used for sortable fields.
	fSortable
	// fClientID is flag used to mark field as allowable to set ClientID.
	fClientID
	// fTime  is a flag used to mark field type as a Time.
	fTime
	// fMap is a flag used to mark field as a map.
	fMap
	// fPtr is a flag used to mark field as a pointer.
	fPtr
	// fArray is a flag used to mark field as an array.
	fArray
	// fSlice is a flag used to mark field as a slice.
	fSlice
	// fBasePtr is flag used to mark field as a based pointer.
	fBasePtr
	// fNestedStruct is a flag used to mark field as a nested structure.
	fNestedStruct
	// fNested is a flag used to mark field as nested.
	fNestedField
	// fCreatedAt defines the created at field
	fCreatedAt
	// fUpdatedAt defines the updated at field
	fUpdatedAt
	// fDeletedAt defines the deleted at field
	fDeletedAt
	// fDatabaseSkip is the flag that defines if a field should be skipped in the database.
	fDatabaseSkip
	// fDatabaseNotNull is the flag that defines that the field is database not nullable.
	fDatabaseNotNull
	// fDatabaseUnique is the flag that defines that the field is database unique.
	fDatabaseUnique
	// fCodecSkip is the flag that defines if a field should be skipped in the codecs.
	fCodecSkip
	// fCodecISO8601 is a time field flag marking it usable with IS08601 formatting.
	fCodecISO8601
)
