package go_parser_combinators

type ParseResult[T any] struct {
	Result T
	Rest   Input
	Err    error
}

type Parser[T any] interface {
	Apply(input Input) ParseResult[T]
}

type ParserFunc[T any] func(Input) ParseResult[T]

func (p ParserFunc[T]) Apply(input Input) ParseResult[T] {
	return p(input)
}

type ParserGenerator[T any] func() Parser[T]

func Lazy[T any](f ParserGenerator[T]) Parser[T] {
	return ParserFunc[T](func(input Input) ParseResult[T] {
		return f().Apply(input)
	})
}
