package app_test

import (
    . "github.com/CharLemAznable/go-lannister/app"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestFixedConfig(t *testing.T) {
    a := assert.New(t)

    config := &Config{}
    FixedConfig(config)
    a.Equal(4791, config.Port)
    a.Equal("lannister", config.ContextPath)
    a.Equal("info", config.LogLevel)
}
