package deepcopy

import (
	"reflect"
	"testing"

	"golang.org/x/tools/go/packages"
)

func TestGetReceiverNames(t *testing.T) {
	config := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps | packages.NeedImports | packages.NeedSyntax,
	}
	path := "../testdata/issue22"
	pkgs, err := packages.Load(config, path)
	if err != nil {
		t.Fatal(err)
	}
	if len(pkgs) != 1 {
		t.Fatalf("unexpected number of packages: %d", len(pkgs))
	}
	if packages.PrintErrors(pkgs) > 0 {
		t.Fatal("packages contain errors")
	}

	got, err := getReceiverNames(pkgs[0])
	if err != nil {
		t.Fatal(err)
	}
	exp := map[string]string{
		"Foo": "xxx",
		"Bar": "yyy",
	}
	if !reflect.DeepEqual(got, exp) {
		t.Fatalf("expected: %v, got: %v", exp, got)
	}
}
