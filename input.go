package go_parser_combinators

type Position struct {
	Offset int
	Line   int
	Column int
}

type Input struct {
	Source   string
	Position Position
}

func NewInput(text string) Input {
	return Input{
		Source: text,
		Position: Position{
			Offset: 0,
			Line:   0,
			Column: 0,
		},
	}
}

func Advance(input Input, consumed int) Input {
	offset := input.Position.Offset + consumed
	line := input.Position.Line
	column := input.Position.Column
	consumedCharSeq := input.Source[:consumed]
	for _, char := range consumedCharSeq {
		if char == '\n' {
			line++
			column = 0
		} else {
			column++
		}
	}
	return Input{
		Source: input.Source[consumed:],
		Position: Position{
			Offset: offset,
			Line:   line,
			Column: column,
		},
	}
}
