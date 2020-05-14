package query

import (
	"github.com/neuronlabs/neuron/errors"
	"github.com/neuronlabs/neuron/mapping"
)

// SplitBracketParameter splits the parameters within the '[' and ']' brackets.
func SplitBracketParameter(bracketed string) (values []string, err error) {
	doubleOpen := func() error {
		err := errors.NewDet(ClassFilterFormat, "double open square brackets")
		err.SetDetailsf("open square bracket '[' found, without closing ']' in: '%s'", bracketed)
		return err
	}

	// set initial indexes
	startIndex := -1
	endIndex := -1

	for i := 0; i < len(bracketed); i++ {
		c := bracketed[i]
		switch c {
		case mapping.AnnotationOpenedBracket:
			if startIndex > endIndex {
				err = doubleOpen()
				return nil, err
			}
			startIndex = i
		case mapping.AnnotationClosedBracket:
			// if opening bracket not set or in case of more than one brackets
			// if start was not set before this endIndex
			if startIndex == -1 || startIndex < endIndex {
				err := errors.NewDet(ClassFilterFormat, "no opening bracket found")
				err.SetDetailsf("close square bracket ']' found, without opening '[' in '%s'", bracketed)
				return nil, err
			}
			endIndex = i
			values = append(values, bracketed[startIndex+1:endIndex])
		}
	}
	if (startIndex != -1 && endIndex == -1) || startIndex > endIndex {
		err = doubleOpen()
		return nil, err
	}
	return values, nil
}
