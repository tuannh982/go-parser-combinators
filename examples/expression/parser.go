package expression

import (
	go_parser_combinators "github.com/tuannh982/go-parser-combinators"
	"github.com/tuannh982/go-parser-combinators/combinators"
	"regexp"
)

var (
	notToken      = combinators.Lit("NOT")
	andToken      = combinators.Lit("AND")
	orToken       = combinators.Lit("OR")
	lparen        = combinators.Lit("(")
	rparen        = combinators.Lit(")")
	fieldName     = combinators.Re(regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*`))
	fieldValue    = combinators.Re(regexp.MustCompile("\"" + `([^"\x00-\x1F\x7F\\]|\\[\\'"bfnrt]|\\u[a-fA-F0-9]{4})*` + "\""))
	binaryOpToken = combinators.Re(regexp.MustCompile("eq|neq|lte|gte|lt|gt"))
	unaryOpToken  = combinators.Lit("IS")
)

func buildParser() go_parser_combinators.Parser[Expression] {
	var not go_parser_combinators.Parser[Expression]
	var query go_parser_combinators.Parser[Expression]
	var term go_parser_combinators.Parser[Expression]
	var factor go_parser_combinators.Parser[Expression]
	var simpleExpr go_parser_combinators.Parser[Expression]
	var unaryExpr go_parser_combinators.Parser[Expression]
	var binaryExpr go_parser_combinators.Parser[Expression]
	binaryExpr = combinators.Map(
		combinators.Seq(combinators.Seq(fieldName, binaryOpToken), fieldValue),
		func(v combinators.Tuple[combinators.Tuple[string, string], string]) Expression {
			return &BinaryExpression{
				Field: v.First.First,
				Op:    v.First.Second,
				Value: v.Second[1 : len(v.Second)-1],
			}
		})
	unaryExpr = combinators.Map(
		combinators.Seq(combinators.SeqL(fieldName, unaryOpToken), fieldValue),
		func(v combinators.Tuple[string, string]) Expression {
			return &UnaryExpression{
				Field: v.First,
				Value: v.Second[1 : len(v.Second)-1],
			}
		})
	simpleExpr = combinators.Or(unaryExpr, binaryExpr)
	not = go_parser_combinators.Lazy(func() go_parser_combinators.Parser[Expression] {
		return combinators.Map(
			combinators.Seq(combinators.Seq(notToken, lparen), combinators.Seq(query, rparen)),
			func(v combinators.Tuple[combinators.Tuple[string, string], combinators.Tuple[Expression, string]]) Expression {
				return &Not{v.Second.First}
			},
		)
	})
	query = go_parser_combinators.Lazy(func() go_parser_combinators.Parser[Expression] {
		return combinators.Map(
			combinators.Seq(term, combinators.Rep(combinators.SeqR(orToken, term))),
			func(v combinators.Tuple[Expression, []Expression]) Expression {
				if len(v.Second) == 0 {
					return v.First
				}
				expressions := []Expression{v.First}
				for _, filter := range v.Second {
					expressions = append(expressions, filter)
				}
				return &Ors{Exps: expressions}
			},
		)
	})
	term = go_parser_combinators.Lazy(func() go_parser_combinators.Parser[Expression] {
		return combinators.Map(
			combinators.Seq(factor, combinators.Rep(combinators.SeqR(andToken, factor))),
			func(v combinators.Tuple[Expression, []Expression]) Expression {
				if len(v.Second) == 0 {
					return v.First
				}
				expressions := []Expression{v.First}
				for _, filter := range v.Second {
					expressions = append(expressions, filter)
				}
				return &Ands{Exps: expressions}
			},
		)
	})
	factor = go_parser_combinators.Lazy(func() go_parser_combinators.Parser[Expression] {
		return combinators.Or(
			simpleExpr,
			combinators.Or(
				combinators.SeqR(lparen, combinators.SeqL(query, rparen)),
				not,
			),
		)
	})
	return query
}
