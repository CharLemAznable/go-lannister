package types

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

var _ = func() bool {
    testing.Init()
    return true
}()

type (
    TestComponent1 struct{}
    TestComponent2 struct{}
)

var registry = NewRegistry("Component")

func TestRegistry(t *testing.T) {
    a := assert.New(t)

    registry.Register("nil", nil)
    a.Nil(registry.Get("nil"))

    component1 := new(TestComponent1)
    component2 := new(TestComponent2)
    registry.Register("same", component1)
    registry.Register("same", component2)
    a.Equal(component1, registry.Get("same"))

    registry.Iterate(nil)
    registry.IterateSorted(nil)
}
