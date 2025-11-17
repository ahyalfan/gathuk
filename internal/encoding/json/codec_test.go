// Package json
package json

import (
	"fmt"
	"testing"

	customtests "github.com/ahyalfan/gathuk/internal/utils/custom-test"
)

type MyStruct struct {
	Name   string         `config:"name"`
	Age    float64        `config:"age"`
	Active bool           `config:"active"`
	Items  []Item         `config:"items"`
	Class  map[string]any `config:"class"`
	Tour   Tour
}

type Item struct {
	ID   string `config:"id"`
	Name string `config:"name"`
}

// type Class struct {
// 	ID   string `config:"id"`
// 	Name string `config:"name"`
// }

type Tour struct {
	ID   string `config:"id"`
	Name string `config:"name"`
	Size []any
}

func TestCodec(t *testing.T) {
	ast := ObjectNode{
		Value: map[string]ASTNode{
			"name":   StringNode{"Alice"},
			"age":    NumberNode{30},
			"active": BooleanNode{true},
			"items": ArrayNode{
				Value: []ASTNode{
					ObjectNode{Value: map[string]ASTNode{
						"id":   StringNode{"1"},
						"name": StringNode{"Item 1"},
					}},
					ObjectNode{Value: map[string]ASTNode{
						"id":   StringNode{"2"},
						"name": StringNode{"Item 2"},
					}},
				},
			},
			"class": ObjectNode{Value: map[string]ASTNode{
				"id":   StringNode{"C1"},
				"name": StringNode{"Math"},
			}},
			"tour": ObjectNode{Value: map[string]ASTNode{
				"id":   StringNode{"T1"},
				"name": StringNode{"Paris Tour"},
				"size": ArrayNode{Value: []ASTNode{
					NumberNode{10},
					StringNode{"large"},
					NumberNode{10.10},
				}},
			}},
		},
	}

	cdc := Codec[MyStruct]{}
	cdcAny := Codec[any]{}
	// cdcMap := Codec[map[string]any]{}

	var result MyStruct
	var resultAny any

	err := cdc.ASTToStruct(ast, &result)
	customtests.OK(t, err)
	err = cdcAny.ASTToStruct(ast, &resultAny)
	customtests.OK(t, err)

	fmt.Printf("%+v\n", result)
	fmt.Printf("%+v\n", resultAny)
	m := resultAny.(map[string]interface{})
	name := m["name"].(string)
	fmt.Printf("%s\n", name)

	t.Run("encode", func(t *testing.T) {
		b, err := cdc.Encode(result)
		customtests.OK(t, err)
		fmt.Println(string(b))
	})
}
