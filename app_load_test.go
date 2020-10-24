package lannister

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

var _ = func() bool {
    testing.Init()
    return true
}()

func TestAppLoad(t *testing.T) {
    a := assert.New(t)
    var appConfig AppConfig

    loadAppConfig("appConfigError.toml", &appConfig)
    a.Equal(0, appConfig.Port)
    a.Equal("", appConfig.ContextPath)
    a.Equal("", appConfig.LogLevel)

    fixedAppConfig(&appConfig)
    a.Equal(4791, appConfig.Port)
    a.Equal("lannister", appConfig.ContextPath)
    a.Equal("info", appConfig.LogLevel)
}
