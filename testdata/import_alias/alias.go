package import_alias

import (
	anotherItem "github.com/globusdigital/deep-copy/testdata/import_alias/another/item"
	"github.com/globusdigital/deep-copy/testdata/import_alias/item"
)

type Data struct {
	Items        []item.Item
	AnotherItems []anotherItem.Item
}
