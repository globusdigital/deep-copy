package deepcopy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

func TestGetReceiverNames(t *testing.T) {
	config := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps | packages.NeedImports | packages.NeedSyntax,
	}
	path := "../testdata/issue22"
	pkgs, err := packages.Load(config, path)
	require.NoError(t, err)

	require.Len(t, pkgs, 1)
	if packages.PrintErrors(pkgs) > 0 {
		t.Fatal("packages contain errors")
	}

	got, err := getReceiverNames(pkgs[0])
	require.NoError(t, err)

	exp := map[string]string{
		"Foo": "xxx",
		"Bar": "yyy",
	}
	assert.Equal(t, exp, got)
}
