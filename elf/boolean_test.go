package elf

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestToBool(t *testing.T) {
    a := assert.New(t)

    a.False(ToBool(""))
    a.True(ToBool("true"))
    a.True(ToBool("TRUE"))
    a.True(ToBool("tRUe"))
    a.True(ToBool("on"))
    a.True(ToBool("tRUe"))
    a.True(ToBool("T"))
    a.False(ToBool("false"))
    a.False(ToBool("f"))
    a.False(ToBool("No"))
    a.False(ToBool("n"))
    a.True(ToBool("on"))
    a.True(ToBool("ON"))
    a.False(ToBool("off"))
    a.False(ToBool("oFf"))
    a.True(ToBool("yes"))
    a.True(ToBool("Y"))
    a.True(ToBool("1"))
    a.False(ToBool("0"))
    a.False(ToBool("blue"))
    a.False(ToBool("true "))
    a.False(ToBool("ono"))
    a.False(ToBool("x gti"))
}
