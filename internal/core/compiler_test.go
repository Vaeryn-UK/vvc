package core

import (
	"reflect"
	"strings"
	"testing"
)

func TestCompile(t *testing.T) {
	t.Run("basic-1", func (t *testing.T) {
		assertCompilation(
			t,
			`start:	ADD	13	15`,
			[]Word{0x01, 0x0d, 0x0f},
		)
	})

	t.Run("comments-1", func (t *testing.T) {
		assertCompilation(
			t,
			`start:	ADD	13	15 # bar foo comment
# foo bar comment`,
			[]Word{0x01, 0x0d, 0x0f},
		)
	})
}

func assertCompilation(t *testing.T, input string, expected Words) {
	code := strings.NewReader(input)
	compiler := NewCompiler()

	err, actual := compiler.Compile(code)

	if err != nil {
		t.Fatalf("Unexpected compile error: %s", err.Error())
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Compilation did not result in expected result.\n# Expected: \n%s\n# Actual:\n%s", expected.ToString(), actual.ToString())
	}
}
