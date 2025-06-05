go-parser-combinators
===
A lightweight, functional-style parser combinator library written in Go. Designed to enable rapid creation of complex
parsers by composing small, reusable parsing functions.

## Features

- ✨ Composable and reusable parser functions
- 📦 No dependencies outside the Go standard library
- 🧪 Easy to test and debug
- 🛠️ Supports basic and advanced grammars

## Installation

```bash
go get github.com/tuannh982/go-parser-combinators
```

## Quick Start

```go
package main

import (
    "fmt"
    "regexp"
    go_parser_combinators "github.com/tuannh982/go-parser-combinators"
    "github.com/tuannh982/go-parser-combinators/combinators"
)

func main() {
    // Parse "hello" followed by a number
    helloParser := combinators.Lit("hello")
    numberParser := combinators.Re(regexp.MustCompile(`\d+`))
    parser := combinators.SeqR(helloParser, numberParser)
    
	query := "hello 123"
    input := go_parser_combinators.NewInput(query)
    result := parser.Apply(input)
    
    if result.Err != nil || len(query) != result.Rest.Position.Offset {
        fmt.Printf("Parse error: %v\n", result.Err)
    } else {
        fmt.Printf("Parsed: %s\n", result.Result)  // Output: Parsed: 123
    }
}
```

## Provided Combinators

### Sequencing Combinators
- **`Seq[A, B](p1, p2)`** - Sequences two parsers and returns both results as a `Tuple[A, B]`
- **`SeqL[A, B](p1, p2)`** - Sequences two parsers but only returns the left result (type `A`)
- **`SeqR[A, B](p1, p2)`** - Sequences two parsers but only returns the right result (type `B`)

### Choice Combinators
- **`Or[A](p1, p2)`** - Tries the first parser, and if it fails, tries the second parser

### Repetition Combinators
- **`Rep[A](p)`** - Applies a parser zero or more times, returning a slice of results
- **`RepNM[A](p, min, max)`** - Applies a parser between `min` and `max` times (inclusive)

### Transformation Combinators
- **`Map[A, B](p, f)`** - Transforms the result of a successful parse using function `f`

### String/Regex Combinators
- **`Lit(str)`** - Matches a literal string exactly (automatically skips leading whitespace)
- **`Re(re)`** - Matches text using a regular expression (automatically skips leading whitespace)

## Recursive Parsers

For parsing recursive grammars (like balanced parentheses), use the `Lazy` function to avoid infinite recursion:

This example creates a parser that recognizes nested parentheses around the letter "x". It can parse simple content like "x", or the same content wrapped in any number of balanced parentheses like "(x)", "((x))", "(((x)))", etc.

```go
func buildParenParser() go_parser_combinators.Parser[string] {
    var parenExpr go_parser_combinators.Parser[string]
    
	// Use Lazy to defer parser construction until needed
    parenExpr = go_parser_combinators.Lazy(func() go_parser_combinators.Parser[string] {
        return combinators.Or(
            combinators.Lit("x"), // Base case
            combinators.Map(
                combinators.SeqR(combinators.Lit("("), combinators.SeqL(parenExpr, combinators.Lit(")"))),
                func(inner string) string {
                    return inner
                },
            ),
        )
    })
	
    return parenExpr
}
// Parses: "x", "(x)", "((x))", "(((x)))", etc.
```

The `Lazy` function wraps a parser generator function, deferring the actual parser construction until the parser is applied to input. This prevents nil reference errors that would otherwise occur when the parser tries to reference itself during initialization.

### Complex Recursive Grammars

While the parentheses example is simple, `Lazy` enables much more complex recursive grammars. For instance, the repository includes a complete boolean expression parser (`examples/expression/`) that handles:

- Binary operators: `eq`, `neq`, `lte`, `gte`, `lt`, `gt`
- Logical operators: `AND`, `OR`, `NOT`
- Parenthetical grouping with proper precedence
- Field names and quoted string values

Example expressions it can parse:
```
fieldA eq "value"
NOT(fieldB gte "123")
((fieldA eq "1") AND (fieldB lt "2")) OR (fieldC gte "3")
```

This demonstrates how multiple recursive parsers can work together to handle sophisticated grammars with proper operator precedence and associativity.

## Use Cases

- DSL interpreters
- Configuration file parsers
- Interpreting simple programming languages
- Log or data format processors

## Contributing

Contributions are welcome! Please open an issue or pull request with your ideas or improvements.

## License

MIT License. See [LICENSE](LICENSE) for details.