package elf

import (
    "github.com/CharLemAznable/sqlx"
    "github.com/stretchr/testify/assert"
    "testing"
)

var sqlBundleRegistry = NewSqlBundleRegistry("SqlBundle")

func TestSqlBundleRegistry(t *testing.T) {
    a := assert.New(t)

    a.Panics(func() {
        sqlBundleRegistry.GetSqlBundle(nil)
    })

    db := sqlx.NewDb(nil, "foo")
    a.Nil(sqlBundleRegistry.GetSqlBundle(db))

    sqlBundleRegistry.Register("foo", "error")
    a.Equal("error", sqlBundleRegistry.GetSqlBundle(db))
}
