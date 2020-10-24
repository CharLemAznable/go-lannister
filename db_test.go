package lannister

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestLoadSqlxDB(t *testing.T) {
    a := assert.New(t)

    appConfig := AppConfig{
        DriverName:     "error",
        DataSourceName: "error",
    }
    a.Nil(loadSqlxDB(&appConfig))
}
