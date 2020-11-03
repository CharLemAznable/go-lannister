package base

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestPrepareConfig(t *testing.T) {
    a := assert.New(t)

    config := PrepareConfig()
    a.Equal(4791, config.Port)
    a.Equal("lannister", config.ContextPath)
    a.Equal("info", config.LogLevel)
    a.NotSame(globalConfig, config)
}
