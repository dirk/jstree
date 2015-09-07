package jstree

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

func TestParseFile(t *testing.T) {
	file := pathRelativeToPackageRoot("test/jstree_test.js")

	_, err := ParseFile(file)
	if err != nil { t.Error(err) }
}

func TestParseProgram(t *testing.T) {
	require := require.New(t)

	program, err := ParseFile(pathRelativeToPackageRoot("test/jstree_test.js"))
	require.NoError(err)

	require.Equal(2, len(program.body))
}
