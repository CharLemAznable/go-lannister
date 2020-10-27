package app_test

import (
    . "github.com/CharLemAznable/go-lannister/app"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestLoadSqlxDB(t *testing.T) {
    a := assert.New(t)

    config := &Config{
        DriverName:     "error",
        DataSourceName: "error",
    }
    a.Nil(LoadSqlxDB(config))
}
