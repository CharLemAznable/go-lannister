package app

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestLoadSqlxDB(t *testing.T) {
    a := assert.New(t)

    config := &Config{
        DriverName:     "error",
        DataSourceName: "error",
    }
    a.Nil(loadSqlxDB(config))
}
