package deepcopy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGenerator(t *testing.T) {
	t.Run("no option", func(t *testing.T) {
		g := NewGenerator()
		assert.Equal(t, Generator{
			methodName: "DeepCopy",
			imports:    map[string]string{},
			fns:        [][]byte{},
		}, g)
	})

	t.Run("isPtrRecv", func(t *testing.T) {
		g := NewGenerator(IsPtrRecv(true))
		assert.Equal(t, Generator{
			methodName: "DeepCopy",
			isPtrRecv:  true,
			imports:    map[string]string{},
			fns:        [][]byte{},
		}, g)
	})

	t.Run("WithMethodName", func(t *testing.T) {
		g := NewGenerator(WithMethodName("FuncDeepCopy"))
		assert.Equal(t, Generator{
			methodName: "FuncDeepCopy",
			imports:    map[string]string{},
			fns:        [][]byte{},
		}, g)
	})

	t.Run("WithMaxDepth", func(t *testing.T) {
		g := NewGenerator(WithMaxDepth(15))
		assert.Equal(t, Generator{
			methodName: "DeepCopy",
			maxDepth:   15,
			imports:    map[string]string{},
			fns:        [][]byte{},
		}, g)
	})

	t.Run("WithSkipLists", func(t *testing.T) {
		sl := SkipLists([]map[string]struct{}{
			{"foo": struct{}{}},
			{"bar": struct{}{}},
		})
		g := NewGenerator(WithSkipLists(sl))
		assert.Equal(t, Generator{
			methodName: "DeepCopy",
			skipLists:  sl,
			imports:    map[string]string{},
			fns:        [][]byte{},
		}, g)
	})

	t.Run("WithBuildTags", func(t *testing.T) {
		bts := []string{"foo", "bar"}
		g := NewGenerator(WithBuildTags(bts))
		assert.Equal(t, Generator{
			methodName: "DeepCopy",
			buildTags:  bts,
			imports:    map[string]string{},
			fns:        [][]byte{},
		}, g)
	})

	t.Run("multiple options", func(t *testing.T) {
		g := NewGenerator(
			IsPtrRecv(true),
			WithMethodName("FuncDeepCopy"),
		)
		assert.Equal(t, Generator{
			isPtrRecv:  true,
			methodName: "FuncDeepCopy",
			imports:    map[string]string{},
			fns:        [][]byte{},
		}, g)
	})
}
