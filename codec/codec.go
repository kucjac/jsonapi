package codec

import (
	"io"

	"github.com/neuronlabs/neuron/mapping"
)

// StructTag is a constant used as a tag that defines models codecs.
const StructTag = "codec"

// Codec is an interface used to marshal and unmarshal data, payload and errors in some encoding type.
type Codec interface {
	// MarshalErrors marshals given errors.
	MarshalErrors(w io.Writer, errors ...*Error) error
	// UnmarshalErrors unmarshal provided errors.
	UnmarshalErrors(r io.Reader) (MultiError, error)
	// MimeType returns the mime type that this codec is defined for.
	MimeType() string
}

// ModelMarshaler is an interface that allows to marshal provided models.
type ModelMarshaler interface {
	// MarshalModels marshal provided models into given codec encoding type. The function should
	// simply encode only provided models without any additional payload like metadata.
	MarshalModels(models []mapping.Model, options ...MarshalOption) ([]byte, error)
	// MarshalModel marshal single models into given codec encoding type.
	MarshalModel(model mapping.Model, options ...MarshalOption) ([]byte, error)
}

// ModelUnmarshaler is an interface that allows to unmarshal provided models of given model struct.
type ModelUnmarshaler interface {
	// UnmarshalModels unmarshals provided data into mapping.Model slice. The data should simply be only encoded models.
	// Requires model or model struct option.
	UnmarshalModels(data []byte, options ...UnmarshalOption) ([]mapping.Model, error)
	// UnmarshalModel unmarshal single model from provided input data. Requires model or model struct option.
	UnmarshalModel(data []byte, options ...UnmarshalOption) (mapping.Model, error)
}
