package yaml

import (
	"reflect"
)

// YamlError is the top level detailed error.
type YamlError struct {
	// Cause will be a structured error with additional details
	// about the error itself.
	Cause error
	// Original is a 1 sentence error for brevity.
	Original error
}



// GoLangStructError are errors that happen when decoding the golang struct,
// before the yaml is touched.
// These errors are usually internal server errors as it's a mistake
// on the Go dev, not the yaml document.
type GoLangStructError struct {
	Err error
}

// YamlTextError happens when mapping the yaml text to a golang struct.
// These errors concern the user who wrote the yaml document.
type YamlTextError struct {
	// Node is the yaml node when the error took place.
	Node Node
	// Path is the yaml hierarchy path to the node encountered for the error.
	// The path include the field attempting to be decoded.
	Path string

	// Cause is the error that caused the TextError.
	Cause error
	// To is the GoLang value the yaml text was attempted to be decoded into.
	To reflect.Value
}

func NewYamlDecodeError(err error, n Node, path string, out reflect.Value) error {
	return YamlError{
		Cause:    YamlTextError{
			Node:  n,
			Path:  path,
			Cause: err,
			To:    out,
		},
		Original: err,
	}
}


func NewAlreadyDefinedError(err error, n Node, path string, out reflect.Value, key string, value bool) error {
	return YamlError{
		Cause:    YamlTextError{
			Node:  n,
			Path:  path,
			Cause: AlreadyDefinedError{
				Key: key,
				Value: value,
			},
			To:    out,
		},
		Original: err,
	}
}

func NewUnknownFieldError(err error, n Node, path string, out reflect.Value, field string) error {
	return YamlError{
		Cause:    YamlTextError{
			Node:  n,
			Path:  path,
			Cause: UnknownFieldError{Field: field},
			To:    out,
		},
		Original: err,
	}
}

func NewWrongTypeError(err error, n Node, path string, out reflect.Value) error {
	var exp Kind

	t := out.Type()
	switch t.Kind() {
	case reflect.Array, reflect.Slice:
		exp = SequenceNode
	case reflect.Map, reflect.Struct:
		exp = MappingNode
	default:
		exp = ScalarNode
	}

	return YamlError{
		Cause:    YamlTextError{
			Node:  n,
			Path:  path,
			Cause: WrongTypeError{
				Expected: exp,
				Of:       underlyingPrimitive(out),
			},
			To:    out,
		},
		Original: err,
	}
}




func (w YamlTextError) Error() string {
	return ""
}

func (w GoLangStructError) Error() string {
	return ""
}

func NewGoLangStructError(err error) error {
	return YamlError{
		Cause: GoLangStructError{
			Err: err,
		},
		Original: err,
	}
}


type UnknownFieldError struct {
	Field string
}



func (w UnknownFieldError) Error() string {
	return ""
}


// WrongTypeError is when the type of the yaml does not match the golang type.
type WrongTypeError struct {
	Expected Kind
	Of       string
}

type AlreadyDefinedError struct {
	Key string
	// Value
	// true: Value already set
	// false: Key already set
	Value bool
}


func (w AlreadyDefinedError) Error() string {
	return ""
}


func (w YamlError) Error() string {
	return ""
}


func (w WrongTypeError) Error() string {
	return ""
}

func underlyingPrimitive(i interface{}) string {
	t := reflect.TypeOf(i)
	switch t.Kind() {
	case reflect.Struct:
		return "key:value" // TODO: @emyrk idk about this one
	case reflect.Map:
		return "key:value"
	case reflect.Slice:
		return "[]value"
	case reflect.Ptr:
		// TODO: @emyrk Catch panic?
		return underlyingPrimitive(reflect.ValueOf(i).Elem().Interface())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "int"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return "uint"
	case reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return "float"
	case reflect.String:
		return "string"
	default:
		return "unknown"
	}
}
