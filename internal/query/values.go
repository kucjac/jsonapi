package query

import (
	"fmt"
	"github.com/neuronlabs/neuron/internal"
	"reflect"
	"strconv"
	"time"
)

func setPrimaryField(value string, fieldValue reflect.Value) (err error) {
	// if the id field is of string type set it to the strValue
	t := fieldValue.Type()

	switch t.Kind() {
	case reflect.String:
		fieldValue.SetString(value)
	case reflect.Int:
		err = setIntField(value, fieldValue, 64)
	case reflect.Int16:
		err = setIntField(value, fieldValue, 16)
	case reflect.Int32:
		err = setIntField(value, fieldValue, 32)
	case reflect.Int64:
		err = setIntField(value, fieldValue, 64)
	case reflect.Uint:
		err = setUintField(value, fieldValue, 64)
	case reflect.Uint16:
		err = setUintField(value, fieldValue, 16)
	case reflect.Uint32:
		err = setUintField(value, fieldValue, 32)
	case reflect.Uint64:
		err = setUintField(value, fieldValue, 64)
	default:
		// should never happen - model checked at precomputation.
		/**

		TO DO:

		Panic - recover
		for internals

		*/
		err = internal.IErrInvalidType
		// err = fmt.Errorf("Internal error. Invalid model primary field format: %v", t)
	}
	return
}

func setAttributeField(value string, fieldValue reflect.Value) (err error) {
	// the attribute can be:
	t := fieldValue.Type()
	switch t.Kind() {
	case reflect.Int:
		err = setIntField(value, fieldValue, 64)
	case reflect.Int8:
		err = setIntField(value, fieldValue, 8)
	case reflect.Int16:
		err = setIntField(value, fieldValue, 16)
	case reflect.Int32:
		err = setIntField(value, fieldValue, 32)
	case reflect.Int64:
		err = setIntField(value, fieldValue, 64)
	case reflect.Uint:
		err = setUintField(value, fieldValue, 64)
	case reflect.Uint16:
		err = setUintField(value, fieldValue, 16)
	case reflect.Uint32:
		err = setUintField(value, fieldValue, 32)
	case reflect.Uint64:
		err = setUintField(value, fieldValue, 64)
	case reflect.String:
		fieldValue.SetString(value)
	case reflect.Bool:
		err = setBoolField(value, fieldValue)
	case reflect.Float32:
		err = setFloatField(value, fieldValue, 32)
	case reflect.Float64:
		err = setFloatField(value, fieldValue, 64)
	case reflect.Struct:
		// check if it is time

		if _, ok := fieldValue.Elem().Interface().(time.Time); ok {
			// it is time
		} else {
			// structs are not allowed as attribute
			err = fmt.Errorf("The struct is not allowed as an attribute. FieldName: '%s'",
				t.Name())
		}
	default:
		// unknown field
		err = fmt.Errorf("Unsupported field type as an attribute: '%s'.", t.Name())
	}
	return
}

func setTimeField(value string, fieldValue reflect.Value) (err error) {
	return
}

func setUintField(value string, fieldValue reflect.Value, bitSize int) (err error) {
	var uintValue uint64

	// Parse unsigned int
	uintValue, err = strconv.ParseUint(value, 10, bitSize)

	if err != nil {
		return err
	}

	// Set uint
	fieldValue.SetUint(uintValue)
	return nil
}

func setIntField(value string, fieldValue reflect.Value, bitSize int) (err error) {
	var intValue int64
	intValue, err = strconv.ParseInt(value, 10, bitSize)
	if err != nil {
		return err
	}

	// Set value if no error
	fieldValue.SetInt(intValue)
	return nil
}

func setFloatField(value string, fieldValue reflect.Value, bitSize int) (err error) {
	var floatValue float64

	// Parse float
	floatValue, err = strconv.ParseFloat(value, bitSize)
	if err != nil {
		return err
	}
	fieldValue.SetFloat(floatValue)
	return nil
}

func setBoolField(value string, fieldValue reflect.Value) (err error) {
	var boolValue bool
	// set default if empty
	if value == "" {
		value = "false"
	}
	boolValue, err = strconv.ParseBool(value)
	if err != nil {
		return err
	}
	fieldValue.SetBool(boolValue)
	return nil
}
