package yaml

import (
	"reflect"
)

type YamlError struct {
	Node Node
	Path string

	Cause error
}

// WrongTypeError is when the type of the yaml does not match the golang type.
type WrongTypeError struct {
	Expected Kind
	Of       string
}

type AlreadyDefinedError struct {
}

func NewYamlError(n *Node, path string, cause error) error {
	return YamlError{
		Node:  *n,
		Path:  path,
		Cause: cause,
	}
}

func (w YamlError) Error() string {
	return ""
}

func NewWrongTypeError(expected interface{}) error {
	var exp Kind

	t := reflect.TypeOf(expected)
	switch t.Kind() {
	case reflect.Array, reflect.Slice:
		exp = SequenceNode
	case reflect.Map, reflect.Struct:
		exp = MappingNode
	default:
		exp = ScalarNode
	}

	return WrongTypeError{
		Expected: exp,
		Of:       underlyingPrimitive(expected),
	}
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
