package combinators

import (
	"fmt"
	"github.com/tuannh982/go-parser-combinators"
)

type Tuple[A, B any] struct {
	First  A
	Second B
}

func Seq[A, B any](p1 go_parser_combinators.Parser[A], p2 go_parser_combinators.Parser[B]) go_parser_combinators.Parser[Tuple[A, B]] {
	return go_parser_combinators.ParserFunc[Tuple[A, B]](func(input go_parser_combinators.Input) go_parser_combinators.ParseResult[Tuple[A, B]] {
		res1 := p1.Apply(input)
		if res1.Err != nil {
			return go_parser_combinators.ParseResult[Tuple[A, B]]{Err: res1.Err}
		}
		res2 := p2.Apply(res1.Rest)
		if res2.Err != nil {
			return go_parser_combinators.ParseResult[Tuple[A, B]]{Err: res2.Err}
		}
		return go_parser_combinators.ParseResult[Tuple[A, B]]{
			Result: Tuple[A, B]{res1.Result, res2.Result},
			Rest:   res2.Rest,
		}
	})
}

func SeqL[A, B any](p1 go_parser_combinators.Parser[A], p2 go_parser_combinators.Parser[B]) go_parser_combinators.Parser[A] {
	return go_parser_combinators.ParserFunc[A](func(input go_parser_combinators.Input) go_parser_combinators.ParseResult[A] {
		res1 := p1.Apply(input)
		if res1.Err != nil {
			return go_parser_combinators.ParseResult[A]{Err: res1.Err}
		}
		res2 := p2.Apply(res1.Rest)
		if res2.Err != nil {
			return go_parser_combinators.ParseResult[A]{Err: res2.Err}
		}
		return go_parser_combinators.ParseResult[A]{
			Result: res1.Result,
			Rest:   res2.Rest,
		}
	})
}

func SeqR[A, B any](p1 go_parser_combinators.Parser[A], p2 go_parser_combinators.Parser[B]) go_parser_combinators.Parser[B] {
	return go_parser_combinators.ParserFunc[B](func(input go_parser_combinators.Input) go_parser_combinators.ParseResult[B] {
		res1 := p1.Apply(input)
		if res1.Err != nil {
			return go_parser_combinators.ParseResult[B]{Err: res1.Err}
		}
		res2 := p2.Apply(res1.Rest)
		if res2.Err != nil {
			return go_parser_combinators.ParseResult[B]{Err: res2.Err}
		}
		return go_parser_combinators.ParseResult[B]{
			Result: res2.Result,
			Rest:   res2.Rest,
		}
	})
}

func Or[A any](p1, p2 go_parser_combinators.Parser[A]) go_parser_combinators.Parser[A] {
	return go_parser_combinators.ParserFunc[A](func(input go_parser_combinators.Input) go_parser_combinators.ParseResult[A] {
		res1 := p1.Apply(input)
		if res1.Err == nil {
			return res1
		}
		return p2.Apply(input)
	})
}

func Rep[A any](p1 go_parser_combinators.Parser[A]) go_parser_combinators.Parser[[]A] {
	return go_parser_combinators.ParserFunc[[]A](func(input go_parser_combinators.Input) go_parser_combinators.ParseResult[[]A] {
		var results []A
		rest := input
		for {
			res := p1.Apply(rest)
			if res.Err != nil {
				break
			}
			results = append(results, res.Result)
			rest = res.Rest
		}
		return go_parser_combinators.ParseResult[[]A]{
			Result: results,
			Rest:   rest,
		}
	})
}

func RepNM[A any](p go_parser_combinators.Parser[A], min, max int) go_parser_combinators.Parser[[]A] {
	return go_parser_combinators.ParserFunc[[]A](func(input go_parser_combinators.Input) go_parser_combinators.ParseResult[[]A] {
		var results []A
		rest := input
		for i := 0; i < max; i++ {
			res := p.Apply(rest)
			if res.Err != nil {
				break
			}
			results = append(results, res.Result)
			rest = res.Rest
		}
		if len(results) < min {
			return go_parser_combinators.ParseResult[[]A]{
				Err: fmt.Errorf("expected at least %d repetitions, got %d", min, len(results)),
			}
		}
		return go_parser_combinators.ParseResult[[]A]{
			Result: results,
			Rest:   rest,
		}
	})
}

func Map[A, B any](p go_parser_combinators.Parser[A], f func(A) B) go_parser_combinators.Parser[B] {
	return go_parser_combinators.ParserFunc[B](func(input go_parser_combinators.Input) go_parser_combinators.ParseResult[B] {
		res := p.Apply(input)
		if res.Err != nil {
			return go_parser_combinators.ParseResult[B]{Err: res.Err}
		}
		return go_parser_combinators.ParseResult[B]{
			Result: f(res.Result),
			Rest:   res.Rest,
		}
	})
}
