package common

import (
    "github.com/CharLemAznable/sqlx"
    "github.com/stretchr/testify/assert"
    "testing"
)

var sqlRegistry = NewSqlRegistry("SqlBundle")

func TestSqlRegistry(t *testing.T) {
    a := assert.New(t)

    a.Panics(func() {
        sqlRegistry.GetSql(nil)
    })

    db := sqlx.NewDb(nil, "foo")
    a.Nil(sqlRegistry.GetSql(db))

    sqlRegistry.Register("foo", "error")
    a.Equal("error", sqlRegistry.GetSql(db))
}
