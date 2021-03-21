package core

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"unicode"
)

type CompileError struct {
	message string
	token *token
}

func (e *CompileError) Error() string {
	if e.token == nil {
		return fmt.Sprintf("Compilation error: %s", e.message)
	}

	return fmt.Sprintf("Compilation error (%d:%d). %s", e.token.line, e.token.column, e.message)
}

type token struct {
	token string
	instruction *Instruction
	value uint
	line uint
	column uint
}

func (t *token) IsLabel() bool {
	return t.token[len(t.token) - 1] == ':'
}

type Compiler struct {
	line uint
	column uint
	data Words
}

func NewCompiler() *Compiler {
	return new(Compiler)
}

func (c *Compiler) Compile(code io.ByteReader) (error, Words) {
	c.line = 1

	data := make(Words, 0)

	for {
		err, done, line := c.consumeLine(code)

		if err != nil {
			return err, nil
		}

		var args []*token
		var inst *token = nil

		for _, token := range line {
			if inst == nil {
				// Expecting label or inst.
				if token.token[len(token.token) - 1] == ':' {
					// Ignore labels for now.
					continue
				} else {
					// Must be an inst.
					if err, i := LookupInstruction(token.token); err != nil {
						return c.compileError("Invalid instruction ref `%s`", token, token.token)
					} else {
						token.instruction = &i
						inst = token
					}
				}
			} else {
				// Must be args.
				if val, err := strconv.Atoi(token.token); err != nil {
					return c.compileError("Arguments must be numeric. Got `%s`", token, token.token)
				} else if !IsValidWord(val) {
					return c.compileError("Invalid representation of word `%s`", token, token.token)
				} else {
					token.value = uint(val)
					args = append(args, token)
				}
			}
		}

		if inst != nil {
			if uint8(len(args)) != inst.instruction.args {
				return c.compileError("Instruction `%s` expects %d args. Got %d", inst, inst.token, inst.instruction.args, len(args))
			}

			data = append(data, Word(inst.instruction.opcode))

			for _, arg := range args {
				data = append(data, Word(arg.value))
			}
		}

		if done {
			// End of code.
			break
		}
	}

	if len(data) == 0 {
		// Cannot have an empty program.
		return c.compileError("An empty program is invalid.", nil)
	}

	return nil, data
}

func (c *Compiler) compileError(msg string, token *token, args ...interface{}) (*CompileError, []Word) {
	return &CompileError{fmt.Sprintf(msg, args...), token}, nil
}

func (c *Compiler) consumeLine(code io.ByteReader) (error, bool, []*token) {
	current := ""
	result := make([]*token, 0)

	c.line++
	c.column = 1
	ignore := false;

	tokenize := func(str string) *token {
		return &token{token: str, line: c.line, column: c.column}
	}

	pushCurrent := func() {
		if len(current) > 0 {
			result = append(result, tokenize(current))
		}

		current = ""
	}

	for {
		if b, err := code.ReadByte(); err == nil {
			c.column++
			if string(b) == "\n" {
				pushCurrent()
				return nil, false, result
			} else if ignore {
				// Ignoring if we've seen a comment marker.
				continue
			} else if string(b) == "#" {
				ignore = true
			} else if unicode.IsSpace(rune(b)) {
				pushCurrent()
			} else {
				current += string(b)
			}
		} else if err == io.EOF {
			pushCurrent()
			return nil, true, result
		} else {
			return err, false, result
		}
	}

	return errors.New("line read didn't term"), true, nil
}
