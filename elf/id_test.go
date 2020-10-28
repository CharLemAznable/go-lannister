package elf

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestNextId(t *testing.T) {
    a := assert.New(t)

    a.True("" != NextId())
}
