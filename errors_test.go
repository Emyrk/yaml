package yaml

import (
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

	t.Run("DuplicatedKey", func(t *testing.T) {
		type s struct {
			A    int
			Nest struct {
				A int
			} `yaml:",inline"`
		}

		y := `
A: 7
Nest:
 A: 2`
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

}
