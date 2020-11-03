package base

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestPrepareDB(t *testing.T) {
    a := assert.New(t)

    config := &Config{
        DriverName:     "error",
        DataSourceName: "error",
    }
    db := PrepareDB(config)
    a.Nil(db.DB)
    a.Equal("", db.DriverName())
}
