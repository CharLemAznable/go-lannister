package elf_test

import (
    . "github.com/CharLemAznable/go-lannister/elf"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestNextId(t *testing.T) {
    a := assert.New(t)

    a.True("" != NextId())
}
