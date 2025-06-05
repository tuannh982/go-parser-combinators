package expression

import (
	"github.com/stretchr/testify/assert"
	go_parser_combinators "github.com/tuannh982/go-parser-combinators"
	"testing"
)

var tests = []struct {
	Query    string
	Expected Expression
}{
	{
		Query:    `NOT(fieldD gte "123")`,
		Expected: &Not{Exp: &BinaryExpression{Field: "fieldD", Op: "gte", Value: "123"}},
	},
	{
		Query: `((fieldA eq "1") AND (((fieldB lt "2")))) OR (fieldC gte "3")`,
		Expected: &Ors{
			Exps: []Expression{
				&Ands{
					Exps: []Expression{
						&BinaryExpression{Field: "fieldA", Op: "eq", Value: "1"},
						&BinaryExpression{Field: "fieldB", Op: "lt", Value: "2"},
					},
				},
				&BinaryExpression{Field: "fieldC", Op: "gte", Value: "3"},
			},
		},
	},
}

func TestExpressions(t *testing.T) {
	parser := buildParser()
	for testNo, test := range tests {
		r := parser.Apply(go_parser_combinators.NewInput(test.Query))
		assert.NoError(t, r.Err)
		assert.Equal(t, len(test.Query), r.Rest.Position.Offset)
		assert.True(t, test.Expected.Equals(r.Result), "test#%d: expected: %v, actual: %v", testNo, test.Expected, r.Result)
	}
}
