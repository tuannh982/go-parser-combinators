package combinators

import (
	"fmt"
	go_parser_combinators "github.com/tuannh982/go-parser-combinators"
	"regexp"
	"strings"
)

var whiteSpace = regexp.MustCompile(`^\s*`)

func skipWhitespace(input go_parser_combinators.Input) go_parser_combinators.Input {
	loc := whiteSpace.FindStringIndex(input.Source)
	if loc != nil && loc[0] == 0 {
		return go_parser_combinators.Advance(input, loc[1])
	}
	return input
}

func Lit(str string) go_parser_combinators.Parser[string] {
	fn := go_parser_combinators.ParserFunc[string](func(input go_parser_combinators.Input) go_parser_combinators.ParseResult[string] {
		start := skipWhitespace(input)
		if strings.HasPrefix(start.Source, str) {
			return go_parser_combinators.ParseResult[string]{
				Result: str,
				Rest:   go_parser_combinators.Advance(start, len(str)),
			}
		}
		found := "EOF"
		if len(start.Source) > 0 {
			found = fmt.Sprintf("'%c'", start.Source[0])
		}
		return go_parser_combinators.ParseResult[string]{Err: fmt.Errorf("'%s' expected but %s found", str, found)}
	})
	return go_parser_combinators.ParserImpl[string]{
		Name: fmt.Sprintf("Lit(%s)", str),
		Fn:   fn,
	}
}

func Re(re *regexp.Regexp) go_parser_combinators.Parser[string] {
	fn := go_parser_combinators.ParserFunc[string](func(input go_parser_combinators.Input) go_parser_combinators.ParseResult[string] {
		start := skipWhitespace(input)
		loc := re.FindStringIndex(start.Source)
		if loc != nil && loc[0] == 0 {
			match := start.Source[loc[0]:loc[1]]
			return go_parser_combinators.ParseResult[string]{
				Result: match,
				Rest:   go_parser_combinators.Advance(start, loc[1]),
			}
		}
		return go_parser_combinators.ParseResult[string]{
			Err: fmt.Errorf("regex match failed: %s", re.String()),
		}
	})
	return go_parser_combinators.ParserImpl[string]{
		Name: fmt.Sprintf("Re(%s)", re.String()),
		Fn:   fn,
	}
}
