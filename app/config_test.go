package app

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestFixedConfig(t *testing.T) {
    a := assert.New(t)

    config := &Config{}
    fixedConfig(config)
    a.Equal(4791, config.Port)
    a.Equal("lannister", config.ContextPath)
    a.Equal("info", config.LogLevel)
}
