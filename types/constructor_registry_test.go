package types

import (
    "github.com/CharLemAznable/sqlx"
    "github.com/stretchr/testify/assert"
    "testing"
)

var daoConsRegistry = NewDaoConstructorRegistry("Dao")

func TestDaoConstructorRegistry(t *testing.T) {
    a := assert.New(t)

    a.Panics(func() {
        daoConsRegistry.GetDaoConstructor(nil)
    })

    db := sqlx.NewDb(nil, "fake")
    a.Panics(func() {
        daoConsRegistry.GetDaoConstructor(db)
    })

    flag := ""
    daoConsRegistry.Register("fake",
        func(db *sqlx.DB) interface{} {
            flag = "fake"
            return nil
        })
    a.Nil(daoConsRegistry.GetDaoConstructor(db).
    (func(db *sqlx.DB) interface{})(db))
    a.Equal("fake", flag)
}
