package compiler

import (
	"fmt"
	"github.com/vaeryn-uk/vvc/internal/core"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type CompileError struct {
	message string
	token   *token
}

func (e *CompileError) Error() string {
	if e.token == nil {
		return fmt.Sprintf("Compilation error: %s", e.message)
	}

	return fmt.Sprintf("Compilation error (%d:%d). %s", e.token.line, e.token.column, e.message)
}

type token struct {
	token       string
	instruction *core.Instruction
	isRegister  bool
	isLabelRef  bool
	label       string
	value       uint
	line        uint
	column      uint
}

type Compiler struct {
	line   uint
	column uint
	data   core.Words
}

func NewCompiler() *Compiler {
	return new(Compiler)
}

func (c *Compiler) Compile(code io.ByteReader) (error, core.Words) {
	c.line = 1

	data := make(core.Words, 0)
	wordAddress := core.Address(0)

	write := func(w core.Word) {
		data = append(data, w)
		wordAddress++
	}

	labelAddresses := make(map[string]core.Address)
	wordsToAddress := make(map[core.Address]string)

	for {
		err, done, line := c.consumeLine(code)

		if err != nil {
			return err, nil
		}

		var args []*token
		var inst *token = nil
		var label *token = nil

		for _, token := range line {
			if inst == nil {
				// Expecting label or inst.
				if token.token[len(token.token)-1] == ':' {
					// Set the lable then continue.
					token.label = token.token[:len(token.token)-1]
					label = token
					continue
				} else {
					// Must be an inst.
					if err, i := core.LookupInstruction(token.token); err != nil {
						return c.compileError("Invalid instruction ref `%s`", token, token.token)
					} else {
						token.instruction = &i
						inst = token
					}
				}
			} else {
				val := token.token

				isLetters := true

				for _, r := range token.token {
					if !unicode.IsLetter(r) {
						isLetters = false
					}
				}

				if isLetters {
					token.isLabelRef = true
					args = append(args, token)
				} else {
					if strings.Index(val, "r") == 0 {
						token.isRegister = true
						val = val[1:]
					}

					if val, err := strconv.Atoi(val); err != nil {
						return c.compileError("Arguments must be numeric. Got `%s`", token, token.token)
					} else if !core.IsValidWord(val) {
						return c.compileError("Invalid representation of word `%s`", token, token.token)
					} else {
						token.value = uint(val)
						args = append(args, token)
					}
				}
			}
		}

		if inst != nil {
			if uint8(len(args)) != inst.instruction.ArgCount() {
				return c.compileError("Instruction `%s` expects %d args. Got %d", inst, inst.token, inst.instruction.ArgCount(), len(args))
			}

			if inst.instruction.A == core.ArgTypeRegister && !args[0].isRegister {
				return c.compileError("Instruction `%s` arg A must be a register. Got `%s`", inst, inst.token, args[0].token)
			}

			if inst.instruction.B == core.ArgTypeRegister && !args[1].isRegister {
				return c.compileError("Instruction `%s` arg B must be a register. Got `%s`", inst, inst.token, args[1].token)
			}

			if label != nil {
				labelAddresses[label.label] = wordAddress
			}

			write(core.Word(inst.instruction.Opcode()))

			for _, arg := range args {
				if arg.isLabelRef {
					write(0)
					wordsToAddress[core.Address(len(data) - 1)] = arg.token
				} else {
					write(core.Word(arg.value))
				}
			}
		}

		if done {
			// End of code.
			break
		}
	}

	// Replace all label reference with their physical addresses.
	for addr, labelRef := range wordsToAddress {
		label, ok := labelAddresses[labelRef]

		if !ok {
			return c.compileError("Label `%s` not found", nil, label)
		}

		data[addr] = core.Word(label)
	}

	if len(data) == 0 {
		// Cannot have an empty program.
		return c.compileError("An empty program is invalid.", nil)
	}

	return nil, data
}

func (c *Compiler) compileError(msg string, token *token, args ...interface{}) (*CompileError, []core.Word) {
	return &CompileError{fmt.Sprintf(msg, args...), token}, nil
}

func (c *Compiler) consumeLine(code io.ByteReader) (error, bool, []*token) {
	current := ""
	result := make([]*token, 0)

	c.line++
	c.column = 1
	ignore := false

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
}
