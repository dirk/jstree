package jstree

import _ "fmt"
import "log"
import "os"
import "path"
import "testing"
import "github.com/stretchr/testify/require"

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func pathRelativeToPackageRoot(p string) string {
	dir, err := os.Getwd()
	if err != nil { log.Fatal(err) }

	return path.Join(dir, "..", "..", p)
}

var __parsedProgram *Program = nil

func getParsedProgram() *Program {
	if __parsedProgram != nil { return __parsedProgram }

	program, err := ParseFile(pathRelativeToPackageRoot("test/jstree_test.js"))
	if err != nil { log.Fatal(err) }

	__parsedProgram = program
	return __parsedProgram
}


// Tests ---------------------------------------------------------------------

func TestParseFile(t *testing.T) {
	file := pathRelativeToPackageRoot("test/jstree_test.js")

	_, err := ParseFile(file)
	if err != nil { t.Error(err) }
}

func TestParseProgram(t *testing.T) {
	require := require.New(t)
	program := getParsedProgram()

	require.Equal(3, len(program.body))
}

func TestParseVariableDeclaration(t *testing.T) {
	require := require.New(t)
	program := getParsedProgram()

	// Expect `var ...` and `const ...` to be variable delcarations.
	for i := 0; i <= 1; i++ {
		variableDeclaration    := program.body[i]
		nilVariableDeclaration := (*VariableDeclaration)(nil)
		require.IsType(nilVariableDeclaration, variableDeclaration)
	}

	varDeclaration := (program.body[0]).(*VariableDeclaration)
	require.Equal("var", varDeclaration.kind)

	constDeclaration := (program.body[1]).(*VariableDeclaration)
	require.Equal("const", constDeclaration.kind)
}

func TestParseFunctionDeclaration(t *testing.T) {
	require := require.New(t)
	program := getParsedProgram()

	functionDeclaration    := program.body[2]
	nilFunctionDeclaration := (*FunctionDeclaration)(nil)
	require.IsType(nilFunctionDeclaration, functionDeclaration)
}
