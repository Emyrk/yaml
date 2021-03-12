package yaml

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestErrorStrings(t *testing.T) {
	t.Run("MappingOnSequence", func(t *testing.T) {
		type s struct {
			Sequence  []string `yaml:"sequence"`
			Sequence2 []string `yaml:"sequence-2"`
		}

		y := `
sequence-2:
 a: b
sequence:
 a: b

`
		var tmp s
		err := Unmarshal([]byte(y), &tmp)
		// Default:
		// yaml: unmarshal errors:
		// 	line 3: cannot unmarshal !!map into []string
		// 	line 5: cannot unmarshal !!map into []string
		fmt.Println(err)
	})

	t.Run("DuplicatedKeyInGoLangStruct", func(t *testing.T) {
		type s struct {
			A    int
			Nest struct {
				A int
			} `yaml:",inline"`
		}

		y := `
A: 7
`
		var tmp s
		err := Unmarshal([]byte(y), &tmp)
		// Default:
		// panic: duplicated key 'a' in struct yaml.s [recovered]
		// 		panic: duplicated key 'a' in struct yaml.s [recovered]
		// 		panic: duplicated key 'a' in struct yaml.s
		fmt.Println(err)
	})

	t.Run("KeyAlreadyDefined", func(t *testing.T) {
		type s struct {
			A int
		}

		y := `
A: 3
A: 4`
		var tmp s
		err := Unmarshal([]byte(y), &tmp)
		// Default:
		// yaml: unmarshal errors:
		// 	line 3: mapping key "A" already defined at line 2
		fmt.Println(err)
	})

	t.Run("InvalidMapKey", func(t *testing.T) {
		type s map[string]string

		y := `
A: 3
B:
 C: "Hello"`
		var tmp s
		err := Unmarshal([]byte(y), &tmp)
		// Default:
		// yaml: unmarshal errors:
		// 	line 3: mapping key "A" already defined at line 2
		fmt.Println(err)
	})

	t.Run("MappingOnScalar", func(t *testing.T) {
		type s map[string]interface{}

		y := `
A: 3
 B: "Hello""`
		var tmp s
		err := Unmarshal([]byte(y), &tmp)
		// Default:
		// 	yaml: line 3: mapping values are not allowed in this context
		fmt.Println(err)
	})

	t.Run("InvalidMapKey", func(t *testing.T) {
		type s alphabet

		y := `
A:
  S: "Hello?"
B:
 P:
  S: "Goodbye"
  I: "Hello"
  B: 
    - 1
`
		var tmp s
		err := Unmarshal([]byte(y), &tmp)
		// Default:
		// 	<nil>
		data, _ := json.Marshal(tmp)
		var buf bytes.Buffer
		json.Indent(&buf, data, "", "\t")
		fmt.Println(buf.String() + "\n")
		fmt.Println(tmp, err)
	})

}


type alphabet struct {
	A prims `yaml:"A"`
	B nest `yaml:"B"`
}

type nest struct {
	P prims `yaml:"P"`
}

type prims struct {
	S string `yaml:"S" json:",omitempty"`
	I int `yaml:"I" json:",omitempty"`
	B bool `yaml:"B" json:",omitempty"`
	F float64 `yaml:"F" json:",omitempty"`
}