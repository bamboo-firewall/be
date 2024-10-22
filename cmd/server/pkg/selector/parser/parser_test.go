package parser

import (
	"testing"
)

func TestSomething(t *testing.T) {
	//input := "! has(my-label) || my-label starts with 'prod' && role in {'frontend','business'} && type == 'production'"
	input := "role in {'agent'} && project == 'atao'"
	sel, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sel.Evaluate(map[string]string{"role": "agent", "project": "atao"}))

	//tokens, err := tokenizer.Tokenize(input)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(tokens)
}
