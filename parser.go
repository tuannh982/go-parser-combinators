package go_parser_combinators

type ParseResult[T any] struct {
	Result T
	Rest   Input
	Err    error
}

type Parser[T any] interface {
	Apply(input Input) ParseResult[T]
	String() string
}

type ParserFunc[T any] func(Input) ParseResult[T]

type ParserImpl[T any] struct {
	Name string
	Fn   ParserFunc[T]
}

func (p ParserImpl[T]) Apply(input Input) ParseResult[T] {
	return p.Fn(input)
}

func (p ParserImpl[T]) String() string {
	return p.Name
}

type ParserGenerator[T any] func() Parser[T]

func Lazy[T any](f ParserGenerator[T]) Parser[T] {
	return ParserImpl[T]{
		Name: "lazy",
		Fn: func(input Input) ParseResult[T] {
			return f().Apply(input)
		},
	}
}
